package img

import (
	"bytes"
	"image"

	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"

	stdError "errors"

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

	image, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, ErrNotSupported
	}

	return &File{
		Size:   len(file),
		Width:  image.Width,
		Height: image.Height,
	}, nil
}

//ChangeBackground will change given image transparent background to given color
func ChangeBackground(file []byte, bgColor color.Color) ([]byte, error) {
	src, format, err := decode(file)
	if err != nil {
		return nil, err
	}

	if format != "png" {
		return file, nil
	}

	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
	return encode(dst, format)
}

//Resize image based on given size
func Resize(file []byte, size int, padding int) ([]byte, error) {
	src, format, err := decode(file)
	if err != nil {
		return nil, err
	}

	//TODO: Very slow to resize images with aspect ratio different than 1:1

	srcBounds := src.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	if (srcW <= size && srcH <= size) || srcW != srcH {
		size = srcW
	}

	padding = size * padding / 100
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	dstBounds := image.Rect(padding, padding, size-padding, size-padding)
	srcBounds = image.Rect(0, 0, srcBounds.Max.X, srcBounds.Max.Y)
	draw.CatmullRom.Scale(dst, dstBounds, src, srcBounds, draw.Src, nil)

	return encode(dst, format)
}

func decode(file []byte) (image.Image, string, error) {
	src, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to decode image")
	}
	return src, format, err
}

func encode(img image.Image, format string) ([]byte, error) {
	var err error
	writer := new(bytes.Buffer)
	if format == "png" {
		err = png.Encode(writer, img)
	} else if format == "jpeg" {
		err = jpeg.Encode(writer, img, nil)
	} else if format == "gif" {
		err = gif.Encode(writer, img, nil)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to encode image to '%s'", format)
	}

	return writer.Bytes(), nil
}
