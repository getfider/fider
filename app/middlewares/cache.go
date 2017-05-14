package middlewares

import "github.com/getfider/fider/app/pkg/web"

// OneYearCache adds Cache-Control header for one year
func OneYearCache() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			c.Response().Header().Add("Cache-Control", "public, max-age=30672000")
			return next(c)
		}
	}
}
