package img

import (
	"bytes"
	"image"

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

//Resize image based on given size
func Resize(file []byte, size int) ([]byte, error) {
	src, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	srcBounds := src.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	if (srcW <= size && srcH <= size) || srcW != srcH {
		return file, nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, srcBounds, draw.Src, nil)

	writer := new(bytes.Buffer)
	if format == "png" {
		err = png.Encode(writer, dst)
	} else if format == "jpeg" {
		err = jpeg.Encode(writer, dst, nil)
	} else if format == "gif" {
		err = gif.Encode(writer, dst, nil)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to encode image to '%s'", format)
	}

	return writer.Bytes(), nil
}
