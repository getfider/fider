package img

import (
	"bytes"
	"image"

	// These are the supported image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	stdError "errors"
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
