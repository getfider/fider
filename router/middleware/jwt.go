package middleware

import (
	"net/http"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/labstack/echo"
)

// JwtGetter gets JWT token from cookie and add into context
func JwtGetter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if cookie, err := c.Cookie("auth"); err == nil {
				if claims, err := auth.Decode(cookie.Value); err == nil {
					c.Set("Claims", claims)
				} else {
					c.Logger().Error(err)
				}
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			query := c.Request().URL.Query()

			jwt := query.Get("jwt")
			if jwt != "" {
				c.SetCookie(&http.Cookie{
					Name:     "auth",
					Value:    jwt,
					HttpOnly: true,
				})

				scheme := "http"
				if c.Request().TLS != nil {
					scheme = "https"
				}

				query.Del("jwt")

				url := scheme + "://" + c.Request().Host + c.Request().URL.Path
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
