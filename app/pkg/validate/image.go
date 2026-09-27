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
//   - JPEG (colour): ~14 megapixels (e.g. 12MP phone photos)
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
// to create one full-size intermediate copy of it.
// Returns -1 if the dimensions are invalid or beyond MaxImageSide.
func EstimateDecodeBytes(width, height int, format string, model color.Model) int64 {
	if width <= 0 || height <= 0 || width > MaxImageSide || height > MaxImageSide {
		return -1
	}

	var decodeBytesPerPixel int64
	switch format {
	case "jpeg":
		// We can't tell from the header if the JPEG is progressive, so assume it is.
		// A progressive decode keeps all DCT coefficients (64 x int32 per 8x8 block,
		// i.e. 4 bytes per pixel per component) in addition to the decoded image.
		switch model {
		case color.GrayModel:
			decodeBytesPerPixel = 4*1 + 1
		case color.YCbCrModel:
			decodeBytesPerPixel = 4*3 + 3
		case color.RGBAModel: // Adobe RGB JPEG: decoded as YCbCr, then converted to RGBA
			decodeBytesPerPixel = 4*3 + 3 + 4
		default: // CMYK or unknown
			decodeBytesPerPixel = 4*4 + 4
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

	// Can't overflow: at most 16384 * 16384 * (20 + 4)
	return int64(width) * int64(height) * (decodeBytesPerPixel + fullSizeIntermediateBytesPerPixel)
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

	estimate := EstimateDecodeBytes(file.Width, file.Height, format, cfg.ColorModel)
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
