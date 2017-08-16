package handlers

import (
	"bytes"
	"image/png"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/letteravatar"
)

//LetterAvatar returns a letter avatar based on name
func LetterAvatar() web.HandlerFunc {
	return func(c web.Context) error {
		name := c.Param("name")
		size, _ := c.ParamAsInt("size")
		img, err := letteravatar.Draw(size, strings.ToUpper(letteravatar.Extract(name)), &letteravatar.Options{PaletteKey: name})
		if err != nil {
			return c.Failure(err)
		}

		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			return c.Failure(err)
		}

		return c.Blob(http.StatusOK, "image/png", buf.Bytes())
	}
}
