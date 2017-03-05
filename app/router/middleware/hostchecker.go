package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

// HostChecker checks for a specific host
func HostChecker(host string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Request().Host != host {
				c.Logger().Errorf("%s is not valid for this operation. Only %s is allowed.", c.Request().Host, host)
				return c.NoContent(http.StatusBadRequest)
			}

			return next(c)
		}
	}
}
