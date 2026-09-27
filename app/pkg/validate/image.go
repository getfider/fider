package validate

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"

	// Register the decoders for every format imagic accepts, so image.DecodeConfig
	// can read their headers regardless of which other packages are imported.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/goenning/imagic"
)

// MaxDecodeBytes is the estimated worst case amount of memory a single image decode
// (plus one full-size intermediate image) may use. Images over this are never decoded.
//
// A decode allocates memory proportional to width*height, regardless of how small the
// compressed file is, so a tiny file can be a "decompression bomb". See
// EstimateDecodeBytes for the per-format estimates. With 256MB the largest images that
// can be decoded are roughly:
//   - JPEG (colour, baseline): ~38 megapixels
//   - JPEG (colour, progressive): ~14 megapixels
//   - PNG (8-bit RGB/RGBA): ~22 megapixels
//   - PNG (8-bit gray/paletted): ~44 megapixels
//   - GIF: ~53 megapixels
const MaxDecodeBytes int64 = 256 << 20

// MaxImageSide is the maximum width or height of an image that we'll decode.
// Some allocations (padding, resize weights) grow with the longest side, so extreme
// aspect ratios (e.g. 40000000x1) are rejected even if their pixel count is small.
const MaxImageSide = 16384

// maxConcurrentDecodes limits how many images are decoded at the same time across the
// whole process, so parallel requests can't multiply the peak memory usage.
const maxConcurrentDecodes = 2

// fullSizeIntermediateBytesPerPixel accounts for one full-size RGBA copy of the image made by
// an operation before the resize (e.g. imagic.Padding in the Favicon handler).
const fullSizeIntermediateBytesPerPixel = 4

var decodeSemaphore = make(chan struct{}, maxConcurrentDecodes)

// ErrImageTooLarge is returned when an image is estimated to need too much memory to decode
var ErrImageTooLarge = errors.New("image is too large to be decoded")

// EstimateDecodeBytes estimates the worst case memory needed to decode an image with
// the given header info (format and color model as returned by image.DecodeConfig) and
// to create one full-size intermediate copy of it. progressive is only used for JPEGs
// (see isProgressiveJPEG).
// Returns -1 if the dimensions are invalid or beyond MaxImageSide.
func EstimateDecodeBytes(width, height int, format string, model color.Model, progressive bool) int64 {
	if width <= 0 || height <= 0 || width > MaxImageSide || height > MaxImageSide {
		return -1
	}

	var decodeBytesPerPixel int64
	switch format {
	case "jpeg":
		// Go's jpeg decoder allocates whole MCUs (up to 32x32 pixels with 4x sampling factors)
		width, height = roundUp(width, 32), roundUp(height, 32)

		// Buffers allocated by image/jpeg, per pixel, assuming no chroma subsampling (worst case):
		//   - Gray: the image.Gray (1)
		//   - YCbCr: the image.YCbCr (3)
		//   - Adobe RGB (RGBAModel): the image.YCbCr (3), converted into an image.RGBA (4)
		//   - CMYK (or unknown): the image.YCbCr (3) + black channel (1), converted into an image.CMYK (4)
		// A progressive decode additionally keeps all DCT coefficients until the end
		// (64 x int32 per 8x8 block, i.e. 4 bytes per pixel per component).
		var imageBytes, components int64
		switch model {
		case color.GrayModel:
			imageBytes, components = 1, 1
		case color.YCbCrModel:
			imageBytes, components = 3, 3
		case color.RGBAModel:
			imageBytes, components = 3+4, 3
		default:
			imageBytes, components = 3+1+4, 4
		}
		decodeBytesPerPixel = imageBytes
		if progressive {
			decodeBytesPerPixel += 4 * components
		}
	case "png":
		// Interlaced PNGs decode each pass into its own image before merging, so
		// assume twice the size of the final image.
		decodeBytesPerPixel = 2 * bytesPerPixel(model)
	case "gif":
		// Only the first frame is decoded, and Go rejects frames larger than the logical screen
		decodeBytesPerPixel = bytesPerPixel(model)
	default:
		decodeBytesPerPixel = 2 * bytesPerPixel(model)
	}

	// Can't overflow: at most 16416 * 16416 * (24 + 4)
	return int64(width) * int64(height) * (decodeBytesPerPixel + fullSizeIntermediateBytesPerPixel)
}

func roundUp(n, multiple int) int {
	return (n + multiple - 1) / multiple * multiple
}

// maxJPEGHeaderScan limits how many bytes isProgressiveJPEG looks at to find the SOF marker.
// Only metadata segments (EXIF, ICC profile, etc) come before it, each at most 64KB.
const maxJPEGHeaderScan = 1 << 20

// isProgressiveJPEG walks the JPEG segments up to the first Start Of Frame marker and
// reports whether it's anything other than baseline/extended sequential (SOF0/SOF1).
// Anything unexpected, malformed or truncated is reported as progressive, so that the
// decode budget falls back to the worst case. When this returns false, image/jpeg follows
// the exact same segments and sees the same SOF0/SOF1 marker (it rejects multiple SOFs).
func isProgressiveJPEG(content []byte) bool {
	data := content
	if len(data) > maxJPEGHeaderScan {
		data = data[:maxJPEGHeaderScan]
	}
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

func bytesPerPixel(model color.Model) int64 {
	if _, ok := model.(color.Palette); ok {
		return 1
	}
	switch model {
	case color.GrayModel, color.AlphaModel:
		return 1
	case color.Gray16Model, color.Alpha16Model:
		return 2
	case color.RGBAModel, color.NRGBAModel, color.CMYKModel:
		return 4
	case color.RGBA64Model, color.NRGBA64Model:
		return 8
	default:
		return 8
	}
}

// CheckDecodeBudget reads only the image header and checks that decoding the image
// stays within MaxDecodeBytes and MaxImageSide.
// Returns imagic.ErrNotSupported for unsupported formats and ErrImageTooLarge when over budget.
func CheckDecodeBudget(content []byte) error {
	// imagic.Parse decides which formats are supported, so uploads and serving agree
	file, err := imagic.Parse(content)
	if err != nil {
		return err
	}

	// imagic.Parse doesn't expose the color model, which we need to estimate the
	// bytes per pixel, so read the header again (this does not decode pixel data).
	cfg, format, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return imagic.ErrNotSupported
	}

	progressive := format == "jpeg" && isProgressiveJPEG(content)
	estimate := EstimateDecodeBytes(file.Width, file.Height, format, cfg.ColorModel, progressive)
	if estimate < 0 || estimate > MaxDecodeBytes {
		return ErrImageTooLarge
	}
	return nil
}

// SafeApply is the only way images should be decoded: it checks the decode budget
// (returning ErrImageTooLarge without decoding if exceeded) and then applies the given
// operations, limiting how many decodes happen concurrently.
func SafeApply(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	if err := CheckDecodeBudget(content); err != nil {
		return nil, err
	}
	return applyWithSemaphore(ctx, content, operations...)
}

// ApplyTrusted applies the operations without checking the decode budget. Only use it for
// trusted images bundled with Fider, never for user content.
func ApplyTrusted(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	return applyWithSemaphore(ctx, content, operations...)
}

func applyWithSemaphore(ctx context.Context, content []byte, operations ...imagic.ImageOperation) ([]byte, error) {
	select {
	case decodeSemaphore <- struct{}{}:
		defer func() { <-decodeSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return imagic.Apply(content, operations...)
}
