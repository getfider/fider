package validate

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"time"

	// Register the decoders for every format imagic accepts, so image.DecodeConfig
	// can read their headers regardless of which other packages are imported.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/goenning/imagic"
)

// MaxDecodeBytes is the estimated worst case amount of memory a single image decode
// (plus the intermediate images created by the operations applied to it) may use.
// Images over this are never decoded.
//
// A decode allocates memory proportional to width*height, regardless of how small the
// compressed file is, so a tiny file can be a "decompression bomb". See
// ImageHeader.EstimateBytes for the per-format estimates. With 320MB and a single
// operation (e.g. a resize) the largest images that can be decoded are roughly:
//   - JPEG (colour, baseline): ~30 megapixels
//   - JPEG (colour, progressive): ~14.6 megapixels (so 12MP phone photos pass)
//   - PNG (8-bit RGB/RGBA/gray): ~28 megapixels
//   - PNG (paletted): ~56 megapixels
//   - GIF: ~67 megapixels
//
// Together with maxConcurrentDecodes this bounds the worst case at ~640MB.
const MaxDecodeBytes int64 = 320 << 20

// MaxImageSide is the maximum width or height of an image that we'll decode.
// Some allocations (padding, resize weights) grow with the longest side, so extreme
// aspect ratios (e.g. 40000000x1) are rejected even if their pixel count is small.
const MaxImageSide = 16384

// maxConcurrentDecodes limits how many images are decoded at the same time across the
// whole process, so parallel requests can't multiply the peak memory usage.
const maxConcurrentDecodes = 2

// decodeWaitTimeout is how long a request waits for a decode slot before giving up with
// ErrDecodeBusy. Requests may hold a DB connection while waiting, so this must be bounded.
var decodeWaitTimeout = 5 * time.Second

// intermediateBytesPerPixel is the memory each operation (imagic.Padding, imagic.Resize,
// imagic.ChangeBackground) may allocate, per source pixel: at most one full-size RGBA
// image (Padding/ChangeBackground), or the horizontal pass of a resize, which is never
// larger than the source.
const intermediateBytesPerPixel = 4

var decodeSemaphore = make(chan struct{}, maxConcurrentDecodes)

// ErrImageTooLarge is returned when an image is estimated to need too much memory to decode
var ErrImageTooLarge = errors.New("image is too large to be decoded")

// ErrDecodeBusy is returned when no decode slot became available in time
var ErrDecodeBusy = errors.New("too many images are being processed, try again later")

// ImageHeader is the information about an image that can be read without decoding it
type ImageHeader struct {
	Width  int
	Height int
	// Format is "png", "gif" or "jpeg"
	Format     string
	ColorModel color.Model
	// Progressive is true for JPEGs that aren't (or might not be) baseline/extended sequential
	Progressive bool
}

// ContentType returns the MIME type of the image format
func (h *ImageHeader) ContentType() string {
	return "image/" + h.Format
}

// ReadImageHeader reads the image header without decoding the pixel data.
// Returns imagic.ErrNotSupported for anything imagic can't process.
func ReadImageHeader(content []byte) (*ImageHeader, error) {
	// Same logic as imagic.Parse (so uploads and serving accept the same formats), which
	// we can't use directly because it doesn't expose the color model we need to
	// estimate the memory needed to decode the image.
	cfg, format, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil || (format != "png" && format != "gif" && format != "jpeg") {
		return nil, imagic.ErrNotSupported
	}
	return &ImageHeader{
		Width:       cfg.Width,
		Height:      cfg.Height,
		Format:      format,
		ColorModel:  cfg.ColorModel,
		Progressive: format == "jpeg" && isProgressiveJPEG(content),
	}, nil
}

// EstimateBytes estimates the worst case memory needed to decode the image and then apply
// the given number of operations to it.
// Returns -1 if the dimensions are invalid or beyond MaxImageSide.
func (h *ImageHeader) EstimateBytes(operations int) int64 {
	width, height := h.Width, h.Height
	if width <= 0 || height <= 0 || width > MaxImageSide || height > MaxImageSide {
		return -1
	}

	var decodeBytesPerPixel int64
	switch h.Format {
	case "jpeg":
		// Go's jpeg decoder allocates whole MCUs (up to 32x32 pixels with 4x sampling factors)
		width, height = roundUp(width, 32), roundUp(height, 32)

		// Buffers allocated by image/jpeg, per pixel, assuming no chroma subsampling (worst case):
		//   - Gray (1 component): the image.Gray (1)
		//   - 3 components: the image.YCbCr (3), converted into an image.RGBA (4) when the
		//     image is Adobe RGB. That can be decided by an APP14 segment after the first
		//     SOS, which image.DecodeConfig never sees, so always assume it.
		//   - 4 components (CMYK/YCCK): the image.YCbCr (3) + black channel (1), converted
		//     into an image.CMYK (4)
		// A progressive decode additionally keeps all DCT coefficients until the end
		// (64 x int32 per 8x8 block, i.e. 4 bytes per pixel per component).
		var imageBytes, components int64
		switch h.ColorModel {
		case color.GrayModel:
			imageBytes, components = 1, 1
		case color.YCbCrModel, color.RGBAModel:
			imageBytes, components = 3+4, 3
		default:
			imageBytes, components = 3+1+4, 4
		}
		decodeBytesPerPixel = imageBytes
		if h.Progressive {
			decodeBytesPerPixel += 4 * components
		}
	case "png":
		// Interlaced PNGs decode each pass into its own image before merging, so
		// assume twice the size of the final image.
		decodeBytesPerPixel = 2 * pngDecodedBytesPerPixel(h.ColorModel)
	case "gif":
		// Only the first frame is decoded (as an image.Paletted), and Go rejects frames
		// larger than the logical screen
		decodeBytesPerPixel = 1
	default:
		return -1
	}

	if operations < 0 {
		operations = 0
	}

	// Can't overflow: at most 16416 * 16416 * (28 + 4 * operations)
	return int64(width) * int64(height) * (decodeBytesPerPixel + int64(operations)*intermediateBytesPerPixel)
}

// pngDecodedBytesPerPixel returns the size of each pixel of the image created by image/png
// for the given color model (as returned by image.DecodeConfig).
func pngDecodedBytesPerPixel(model color.Model) int64 {
	if _, ok := model.(color.Palette); ok {
		return 1 // image.Paletted
	}
	switch model {
	case color.GrayModel:
		// 1/2/4/8-bit gray is decoded to an image.Gray, but to an image.NRGBA if there's a
		// tRNS chunk, which image.DecodeConfig doesn't read for non paletted images.
		return 4
	case color.Gray16Model:
		// image.Gray16, or image.NRGBA64 with tRNS
		return 8
	case color.RGBAModel, color.NRGBAModel:
		// 8-bit RGB (image.RGBA, image.NRGBA with tRNS), gray+alpha and RGBA (image.NRGBA)
		return 4
	default:
		// 16-bit RGB (image.RGBA64, image.NRGBA64 with tRNS), gray+alpha and RGBA (image.NRGBA64)
		return 8
	}
}

func roundUp(n, multiple int) int {
	return (n + multiple - 1) / multiple * multiple
}

// isProgressiveJPEG walks the JPEG segments up to the first Start Of Frame marker and
// reports whether it's anything other than baseline/extended sequential (SOF0/SOF1).
// Anything unexpected, malformed or truncated is reported as progressive, so that the
// decode budget falls back to the worst case. When this returns false, image/jpeg follows
// the exact same segments and sees the same SOF0/SOF1 marker (it rejects multiple SOFs).
// Segments are skipped by their length, so this is cheap even for large files.
func isProgressiveJPEG(data []byte) bool {
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 { // SOI
		return true
	}

	i := 2
	for {
		if i >= len(data) || data[i] != 0xFF {
			return true
		}
		// A marker may be preceded by any number of 0xFF fill bytes
		for i < len(data) && data[i] == 0xFF {
			i++
		}
		if i >= len(data) {
			return true
		}
		marker := data[i]
		i++

		switch {
		case marker == 0xC0 || marker == 0xC1: // SOF0 (baseline), SOF1 (extended sequential)
			return false
		case marker == 0xC4 || marker == 0xC8 || marker == 0xCC: // DHT, JPG, DAC: not frames
		case marker >= 0xC2 && marker <= 0xCF: // SOF2 (progressive), SOF6/SOF10/SOF14 and others
			return true
		case marker == 0x01 || (marker >= 0xD0 && marker <= 0xD7): // TEM, RSTn: no length
			continue
		case marker == 0x00 || marker == 0xD8 || marker == 0xD9 || marker == 0xDA:
			// Stuffed byte, SOI, EOI or SOS before any SOF: not a well formed header
			return true
		}

		if i+2 > len(data) {
			return true
		}
		length := int(data[i])<<8 | int(data[i+1])
		if length < 2 {
			return true
		}
		i += length
	}
}

// CheckDecodeBudget reads only the image header and checks that decoding the image and
// applying the given number of operations stays within MaxDecodeBytes and MaxImageSide.
// Returns imagic.ErrNotSupported for unsupported formats, and ErrImageTooLarge (with the
// header) when over budget.
func CheckDecodeBudget(content []byte, operations int) (*ImageHeader, error) {
	header, err := ReadImageHeader(content)
	if err != nil {
		return nil, err
	}

	estimate := header.EstimateBytes(operations)
	if estimate < 0 || estimate > MaxDecodeBytes {
		return header, ErrImageTooLarge
	}
	return header, nil
}

// SafeApply is the way user images should be decoded: it checks the decode budget
// (returning ErrImageTooLarge without decoding if exceeded) and then applies the given
// operations, limiting how many decodes happen concurrently (see decodeWithinBudget).
func SafeApply(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	if _, err := CheckDecodeBudget(content, len(operations)); err != nil {
		return nil, err
	}
	return decodeWithinBudget(ctx, content, operations...)
}

// SafeResize is like SafeApply with imagic.Resize(size), but returns the content unchanged
// (without decoding it) when the image is already within size.
func SafeResize(ctx context.Context, content []byte, size int) ([]byte, error) {
	header, err := CheckDecodeBudget(content, 1)
	if err != nil {
		return nil, err
	}
	if size >= header.Width && size >= header.Height {
		return content, nil
	}
	return decodeWithinBudget(ctx, content, imagic.Resize(size))
}

// ApplyTrusted applies the operations without checking the decode budget. Only use it for
// trusted images bundled with Fider, never for user content.
func ApplyTrusted(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	return decodeWithinBudget(ctx, content, operations...)
}

// decodeWithinBudget applies the operations once a decode slot is available. The caller
// must have checked the decode budget. Returns ErrDecodeBusy if no slot becomes available
// within decodeWaitTimeout, or the context error if the request is cancelled.
func decodeWithinBudget(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	timer := time.NewTimer(decodeWaitTimeout)
	defer timer.Stop()

	select {
	case decodeSemaphore <- struct{}{}:
		defer func() { <-decodeSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, ErrDecodeBusy
	}
	return imagic.Apply(content, operations...)
}
