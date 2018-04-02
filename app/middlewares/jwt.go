package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

// JwtGetter gets JWT token from cookie and insert into context
func JwtGetter() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {

			cookie, err := c.Cookie(web.CookieAuthName)
			if err != nil {
				if errors.Cause(err) == http.ErrNoCookie {
					return next(c)
				}
				return err
			}

			claims, err := jwt.DecodeFiderClaims(cookie.Value)
			if err != nil {
				c.RemoveCookie(web.CookieAuthName)
				return next(c)
			}

			user, err := c.Services().Users.GetByID(claims.UserID)
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					c.RemoveCookie(web.CookieAuthName)
					return next(c)
				}
				return err
			}

			if c.Tenant() != nil && user.Tenant.ID == c.Tenant().ID {
				c.SetUser(user)
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {

			query := c.Request.URL.Query()

			token := query.Get("token")
			if token != "" {
				c.AddCookie(web.CookieAuthName, token, time.Now().Add(365*24*time.Hour))

				query.Del("token")

				url := c.BaseURL() + c.Request.URL.Path
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
