package middlewares

import (
	"net/http"
	"strings"

	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage"
)

// JwtGetter gets JWT token from cookie and add into context
func JwtGetter(users storage.User) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {

			if cookie, err := c.Cookie("auth"); err == nil {
				if claims, err := jwt.Decode(cookie.Value); err == nil {
					if user, err := users.GetByID(claims.UserID); err == nil {
						if user.Tenant.ID == c.Tenant().ID {
							c.SetUser(user)
						}
					}
				} else {
					c.Logger().Error(err)
				}
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {

			query := c.Request().URL.Query()

			jwt := query.Get("jwt")
			if jwt != "" {
				c.SetCookie(&http.Cookie{
					Name:     "auth",
					Value:    jwt,
					HttpOnly: true,
					Path:     "/",
				})

				query.Del("jwt")

				url := c.BaseURL() + c.Request().URL.Path
				querystring := query.Encode()
				if querystring != "" {
					url += "?" + querystring
				}

				return c.Redirect(http.StatusTemporaryRedirect, url)
			}

			return next(c)
		}
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
