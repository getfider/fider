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

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/imagic"
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
			userByID := &query.GetUserByID{UserID: id}
			err := bus.Dispatch(c, userByID)
			if err == nil && userByID.Result.Tenant.ID == c.Tenant().ID {
				if userByID.Result.Email != "" {
					url := fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=%d&d=404", crypto.MD5(strings.ToLower(userByID.Result.Email)), size)
					cacheKey := fmt.Sprintf("gravatar:%s", url)

					//If gravatar was found in cache
					if image, found := c.Engine().Cache().Get(cacheKey); found {
						log.Debugf(c, "Gravatar found in cache: @{GravatarURL}", dto.Props{
							"GravatarURL": cacheKey,
						})
						imageInBytes := image.([]byte)
						return c.Image(http.DetectContentType(imageInBytes), imageInBytes)
					}

					log.Debugf(c, "Requesting gravatar: @{GravatarURL}", dto.Props{
						"GravatarURL": url,
					})

					req := &cmd.HTTPRequest{
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
			q := &query.GetBlobByKey{Key: bkey}
			err := bus.Dispatch(c, q)
			if err != nil {
				return c.Failure(err)
			}
			bytes = q.Result.Content
			contentType = q.Result.ContentType
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

		opts := []imagic.ImageOperation{}
		if size > 0 {
			opts = append(opts, imagic.Padding(size*10/100))
			opts = append(opts, imagic.Resize(size))
		}

		if c.QueryParam("bg") != "" {
			opts = append(opts, imagic.ChangeBackground(color.White))
		}

		bytes, err = imagic.Apply(bytes, opts...)
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

		q := &query.GetBlobByKey{Key: bkey}
		err = bus.Dispatch(c, q)
		if err != nil {
			return c.Failure(err)
		}

		bytes := q.Result.Content
		if size > 0 {
			bytes, err = imagic.Apply(bytes, imagic.Resize(size))
			if err != nil {
				return c.Failure(err)
			}
		}

		return c.Image(q.Result.ContentType, bytes)
	}
}
