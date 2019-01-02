package img

import (
	"bytes"
	"image"

	"image/color"
	"image/png"

	stdError "errors"

	"github.com/disintegration/imaging"
	"github.com/getfider/fider/app/pkg/errors"
	"golang.org/x/image/draw"
)

//ErrNotSupported returned by Parse when given file is not in a supported format
var ErrNotSupported = stdError.New("File not supported")

//File is an image supported by Fider
type File struct {
	Width  int
	Height int
	Size   int
}

//Parse returns the a img.File if it's supported by fider
func Parse(file []byte) (*File, error) {
	reader := bytes.NewReader(file)

	image, format, err := image.DecodeConfig(reader)
	if err != nil || (format != "png" && format != "gif" && format != "jpeg") {
		return nil, ErrNotSupported
	}

	return &File{
		Size:   len(file),
		Width:  image.Width,
		Height: image.Height,
	}, nil
}

//ImageOperation is an operation that can be performed on an image and retun a modified version of it
type ImageOperation func(image.Image, string) image.Image

//ChangeBackground will change given image transparent background to given color
func ChangeBackground(bgColor color.Color) ImageOperation {
	return func(src image.Image, format string) image.Image {
		if format != "png" {
			return src
		}

		dst := image.NewRGBA(src.Bounds())
		draw.Draw(dst, dst.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return dst
	}
}

//Padding adds a padding based on given value
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

//Resize image based on given size
func Resize(size int) ImageOperation {
	return func(src image.Image, format string) image.Image {
		b := src.Bounds()
		srcW, srcH := b.Dx(), b.Dy()
		if size >= srcH && size >= srcW {
			return src
		}
		if srcW > srcH {
			return imaging.Resize(src, size, 0, imaging.Lanczos)
		}
		return imaging.Resize(src, 0, size, imaging.Lanczos)
	}
}

// Apply a list of operations on a given image
func Apply(input []byte, operations ...ImageOperation) ([]byte, error) {
	img, format, err := decode(input)
	if err != nil {
		return nil, err
	}

	for _, op := range operations {
		img = op(img, format)
	}

	return encode(img, format)
}

func decode(file []byte) (image.Image, string, error) {
	src, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to decode image")
	}
	return src, format, err
}

func encode(img image.Image, format string) ([]byte, error) {
	writer := new(bytes.Buffer)
	f, _ := imaging.FormatFromExtension(format)

	err := imaging.Encode(
		writer, img, f,
		imaging.PNGCompressionLevel(png.BestCompression),
		imaging.JPEGQuality(95),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode image to '%s'", format)
	}
	return writer.Bytes(), nil
}
