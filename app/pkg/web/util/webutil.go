package webutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/services/blob"
)

func encode(user *models.User) string {
	token, err := jwt.Encode(jwt.FiderClaims{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		Origin:    jwt.FiderClaimsOriginUI,
		Metadata: jwt.Metadata{
			ExpiresAt: time.Now().Add(365 * 24 * time.Hour).Unix(),
		},
	})

	if err != nil {
		panic(errors.Wrap(err, "failed to add auth cookie"))
	}

	return token
}

//AddAuthUserCookie generates Auth Token and adds a cookie
func AddAuthUserCookie(ctx *web.Context, user *models.User) {
	AddAuthTokenCookie(ctx, encode(user))
}

//AddAuthTokenCookie adds given token to a cookie
func AddAuthTokenCookie(ctx *web.Context, token string) {
	expiresAt := time.Now().Add(365 * 24 * time.Hour)
	ctx.AddCookie(web.CookieAuthName, token, expiresAt)
}

//SetSignUpAuthCookie sets a temporary domain-wide Auth Token
func SetSignUpAuthCookie(ctx *web.Context, user *models.User) {
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:     web.CookieSignUpAuthName,
		Domain:   env.MultiTenantDomain(),
		Value:    encode(user),
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		Secure:   ctx.Request.IsSecure,
	})
}

//GetSignUpAuthCookie returns the temporary temporary domain-wide Auth Token and removes it
func GetSignUpAuthCookie(ctx *web.Context) string {
	cookie, err := ctx.Request.Cookie(web.CookieSignUpAuthName)
	if err == nil {
		http.SetCookie(ctx.Response, &http.Cookie{
			Name:     web.CookieSignUpAuthName,
			Domain:   env.MultiTenantDomain(),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
			Expires:  time.Now().Add(-100 * time.Hour),
			Secure:   ctx.Request.IsSecure,
		})
		return cookie.Value
	}
	return ""
}

// ProcessMultiImageUpload uploads multiple image to blob (if it's a new one)
func ProcessMultiImageUpload(c *web.Context, imgs []*models.ImageUpload, preffix string) error {
	for _, img := range imgs {
		err := ProcessImageUpload(c, img, preffix)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessImageUpload uploads image to blob (if it's a new one)
func ProcessImageUpload(c *web.Context, img *models.ImageUpload, preffix string) error {
	if img.Upload != nil && len(img.Upload.Content) > 0 {
		bkey := fmt.Sprintf("%s/%s-%s", preffix, rand.String(64), blob.SanitizeFileName(img.Upload.FileName))
		err := bus.Dispatch(c, &cmd.StoreBlob{
			Key:         bkey,
			Content:     img.Upload.Content,
			ContentType: img.Upload.ContentType,
		})
		if err != nil {
			return errors.Wrap(err, "failed to upload new blob")
		}
		img.BlobKey = bkey
	}
	return nil
}
