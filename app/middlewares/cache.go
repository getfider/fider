package middlewares

import (
	"fmt"
	"time"

	"github.com/getfider/fider/app/pkg/web"
)

// ClientCache adds Cache-Control header for X seconds
func ClientCache(d time.Duration) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Response.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.f", d.Seconds()))
			return next(c)
		}
	}
}
