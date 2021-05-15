package webutil

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

func encode(user *entity.User) string {
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
func AddAuthUserCookie(ctx *web.Context, user *entity.User) {
	AddAuthTokenCookie(ctx, encode(user))
}

//AddAuthTokenCookie adds given token to a cookie
func AddAuthTokenCookie(ctx *web.Context, token string) {
	expiresAt := time.Now().Add(365 * 24 * time.Hour)
	ctx.AddCookie(web.CookieAuthName, token, expiresAt)
}

//SetSignUpAuthCookie sets a temporary domain-wide Auth Token
func SetSignUpAuthCookie(ctx *web.Context, user *entity.User) {
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
