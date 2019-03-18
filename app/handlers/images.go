package handlers

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/services/httpclient"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/img"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/letteravatar"
)

//LetterAvatar returns a letter gravatar picture based on given name
func LetterAvatar() web.HandlerFunc {
	return func(c *web.Context) error {
		id := c.Param("id")
		name := c.Param("name")
		if name == "" {
			name = "?"
		}

		size, err := c.QueryParamAsInt("size")
		if err != nil {
			return c.BadRequest(web.Map{})
		}
		size = between(size, 50, 200)

		img, err := letteravatar.Draw(size, strings.ToUpper(letteravatar.Extract(name)), &letteravatar.Options{
			PaletteKey: fmt.Sprintf("%s:%s", id, name),
		})
		if err != nil {
			return c.Failure(err)
		}

		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			return c.Failure(err)
		}

		return c.Image("image/png", buf.Bytes())
	}
}

//Gravatar returns a gravatar picture of fallsback to letter avatar based on name
func Gravatar() web.HandlerFunc {
	return func(c *web.Context) error {
		id, err := c.ParamAsInt("id")
		if err != nil {
			return c.NotFound()
		}

		size, err := c.QueryParamAsInt("size")
		if err != nil {
			return c.BadRequest(web.Map{})
		}

		size = between(size, 50, 200)

		if err == nil && id > 0 {
			user, err := c.Services().Users.GetByID(id)
			if err == nil && user.Tenant.ID == c.Tenant().ID {
				if user.Email != "" {
					url := fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=%d&d=404", crypto.MD5(strings.ToLower(user.Email)), size)
					cacheKey := fmt.Sprintf("gravatar:%s", url)

					//If gravatar was found in cache
					if image, found := c.Engine().Cache().Get(cacheKey); found {
						log.Debugf(c, "Gravatar found in cache: @{GravatarURL}", log.Props{
							"GravatarURL": cacheKey,
						})
						imageInBytes := image.([]byte)
						return c.Image(http.DetectContentType(imageInBytes), imageInBytes)
					}

					log.Debugf(c, "Requesting gravatar: @{GravatarURL}", log.Props{
						"GravatarURL": url,
					})

					req := &httpclient.Request{
						URL:    url,
						Method: "GET",
					}
					err := bus.Dispatch(c, req)
					if err == nil && req.ResponseStatusCode == http.StatusOK {
						bytes := req.ResponseBody
						c.Engine().Cache().Set(cacheKey, bytes, 24*time.Hour)
						return c.Image(http.DetectContentType(bytes), bytes)
					}
				}
			}
		}

		return LetterAvatar()(c)
	}
}

//Favicon returns the Fider favicon by given size
func Favicon() web.HandlerFunc {
	return func(c *web.Context) error {
		var (
			bytes       []byte
			err         error
			contentType string
		)

		bkey := c.Param("bkey")
		if bkey != "" {
			cmd := &blob.RetrieveBlob{Key: bkey}
			err := bus.Dispatch(c, cmd)
			if err != nil {
				return c.Failure(err)
			}
			bytes = cmd.Blob.Content
			contentType = cmd.Blob.ContentType
		} else {
			bytes, err = ioutil.ReadFile(env.Path("favicon.png"))
			contentType = "image/png"
			if err != nil {
				return c.Failure(err)
			}
		}

		size, err := c.QueryParamAsInt("size")
		if err != nil {
			return c.BadRequest(web.Map{})
		}

		size = between(size, 50, 200)

		opts := []img.ImageOperation{}
		if size > 0 {
			opts = append(opts, img.Padding(size*10/100))
			opts = append(opts, img.Resize(size))
		}

		if c.QueryParam("bg") != "" {
			opts = append(opts, img.ChangeBackground(color.White))
		}

		bytes, err = img.Apply(bytes, opts...)
		if err != nil {
			return c.Failure(err)
		}

		return c.Image(contentType, bytes)
	}
}

//ViewUploadedImage returns any uploaded image by given ID and size
func ViewUploadedImage() web.HandlerFunc {
	return func(c *web.Context) error {
		bkey := c.Param("bkey")

		size, err := c.QueryParamAsInt("size")
		if err != nil {
			return c.BadRequest(web.Map{})
		}

		size = between(size, 0, 2000)

		cmd := &blob.RetrieveBlob{Key: bkey}
		err = bus.Dispatch(c, cmd)
		if err != nil {
			return c.Failure(err)
		}

		bytes := cmd.Blob.Content
		if size > 0 {
			bytes, err = img.Apply(bytes, img.Resize(size))
			if err != nil {
				return c.Failure(err)
			}
		}

		return c.Image(cmd.Blob.ContentType, bytes)
	}
}
