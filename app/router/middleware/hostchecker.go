package middleware

import (
	"net/http"

	"net/url"

	"github.com/labstack/echo"
)

// HostChecker checks for a specific host
func HostChecker(baseURL string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, _ := url.Parse(baseURL)

			if c.Request().Host != u.Host {
				c.Logger().Errorf("%s is not valid for this operation. Only %s is allowed.", c.Request().Host, u.Host)
				return c.NoContent(http.StatusBadRequest)
			}

			return next(c)
		}
	}
}
