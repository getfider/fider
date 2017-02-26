package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// HostChecker checks for a specific host
func HostChecker(baseURL string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			scheme := "http"
			if c.Request().TLS != nil {
				scheme = "https"
			}
			url := scheme + "://" + c.Request().Host

			if strings.Index(url, baseURL) < 0 {
				c.Logger().Errorf("%s is not valid for this operation", url)
				return c.NoContent(http.StatusBadRequest)
			}
			return next(c)
		}
	}
}
