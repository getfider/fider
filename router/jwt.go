package router

import (
	"net/http"
	"os"

	"github.com/WeCanHearYou/wchy/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var jwtSecret = os.Getenv("JWT_SECRET")

// JwtGetter gets JWT token from cookie and add into context
func JwtGetter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if cookie, err := c.Cookie("auth"); err == nil {
				claims := &auth.WchyClaims{}
				token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})

				if err == nil && token.Valid {
					c.Set("Claims", claims)
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

				url := scheme + "://" + c.Request().Host + c.Request().URL.Path + "?" + query.Encode()
				return c.Redirect(http.StatusTemporaryRedirect, url)
			}

			return next(c)
		}
	}
}
