package validate

import (
	"bytes"
	"image"

	// Register the decoders for every format we accept, so image.DecodeConfig
	// can read their headers regardless of which other packages are imported.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// MaxImagePixels is the maximum number of pixels (width * height) an image may have
// before we refuse to fully decode it.
//
// Decoding an image allocates roughly width*height*4 bytes (NRGBA), regardless of how
// small the compressed file is. A 12000x12000 PNG of a single colour compresses to a
// few hundred KB, but decodes to ~576MB. With a 40 megapixel budget the worst case
// decode is ~160MB, while still accepting photos from modern cameras and phones
// (e.g. 24MP phone photos, 36MP DSLR images) that fit within the 5MB upload limit.
const MaxImagePixels int64 = 40_000_000

// IsWithinPixelBudget returns true if an image of the given dimensions is safe to decode.
// Non-positive dimensions are rejected.
func IsWithinPixelBudget(width, height int) bool {
	if width <= 0 || height <= 0 {
		return false
	}
	// Check each side first so the multiplication below can't overflow int64
	if int64(width) > MaxImagePixels || int64(height) > MaxImagePixels {
		return false
	}
	return int64(width)*int64(height) <= MaxImagePixels
}

// ImageWithinPixelBudget reads only the image header (via image.DecodeConfig, which does
// not decode pixel data) and reports whether the image is safe to fully decode.
// An error is returned if the image format is unknown or the header can't be read.
func ImageWithinPixelBudget(content []byte) (bool, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return false, err
	}
	return IsWithinPixelBudget(cfg.Width, cfg.Height), nil
}
