package handlers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/letteravatar"
)

//Avatar returns a gravatar picture of fallsback to letter avatar based on name
func Avatar() web.HandlerFunc {
	return func(c web.Context) error {
		name := c.Param("name")
		size, _ := c.ParamAsInt("size")
		email := c.QueryParam("e")

		id, err := c.ParamAsInt("id")
		if err == nil && id > 0 && email == "" {
			user, err := c.Services().Users.GetByID(id)
			if err == nil && user.Tenant.ID == c.Tenant().ID {
				email = user.Email
				println(user.Email)
			}
		}

		if email != "" {
			hash := md5.Sum([]byte(email))
			url := fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d&d=404", hash, size)
			c.Logger().Debugf("Requesting gravatar: %s", url)
			resp, err := http.Get(url)
			if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
				bytes, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					return c.Blob(http.StatusOK, "image/png", bytes)
				}
			}
		}

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
