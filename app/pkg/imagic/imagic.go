package imagic

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

var ErrNotSupported = errors.New("File not supported")

// File contains metadata of a given image
type File struct {
	Width  int
	Height int
	Size   int
}

// Parse returns a File if it's in a supported format
func Parse(file []byte) (*File, error) {
	reader := bytes.NewReader(file)
	cfg, format, err := image.DecodeConfig(reader)
	if err != nil || (format != "png" && format != "gif" && format != "jpeg" && format != "webp") {
		return nil, ErrNotSupported
	}

	return &File{
		Size:   len(file),
		Width:  cfg.Width,
		Height: cfg.Height,
	}, nil
}

// ImageOperation is an operation that can be performed on an image and return a modified version of it
type ImageOperation func(image.Image, string) image.Image

// ChangeBackground changes a transparent background to the given color (only if PNG)
func ChangeBackground(bgColor color.Color) ImageOperation {
	return func(src image.Image, format string) image.Image {
		if format != "png" && format != "webp" {
			return src
		}
		dst := image.NewRGBA(src.Bounds())
		draw.Draw(dst, dst.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return dst
	}
}

// Padding adds padding (in pixels) around the image
func Padding(padding int) ImageOperation {
	return func(src image.Image, format string) image.Image {
		if padding == 0 {
			return src
		}
		srcBounds := src.Bounds()
		srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

		dst := image.NewRGBA(image.Rect(0, 0, srcW+padding, srcH+padding))
		draw.Draw(dst, dst.Bounds(), src, image.Pt(-padding/2, -padding/2), draw.Src)
		return dst
	}
}

// Resize image to fit within the given size (whichever dimension is larger)
func Resize(size int) ImageOperation {
	return func(src image.Image, format string) image.Image {
		b := src.Bounds()
		srcW, srcH := b.Dx(), b.Dy()

		if size >= srcW && size >= srcH {
			return src
		}

		if srcW > srcH {
			return imaging.Resize(src, size, 0, imaging.NearestNeighbor)
		}
		return imaging.Resize(src, 0, size, imaging.NearestNeighbor)
	}
}

// Apply a list of operations on a given image
// Returns the final image bytes in WEBP format
func Apply(input []byte, operations ...ImageOperation) ([]byte, error) {
	img, format, err := decode(input)
	if err != nil {
		return nil, err
	}

	// Apply each operation in order
	for _, op := range operations {
		img = op(img, format)
	}

	// Encode final result as WebP
	return encodeWebP(img)
}

// decode attempts to read the image bytes into an image.Image, returning its format
func decode(file []byte) (image.Image, string, error) {
	src, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, "", err
	}
	return src, format, nil
}

// encodeWebP encodes the image.Image into WebP
func encodeWebP(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	options := &webp.Options{
		Lossless: false,
		Quality:  80,
	}
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
