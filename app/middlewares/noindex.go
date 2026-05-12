package middlewares

import (
	"github.com/getfider/fider/app/pkg/web"
)

// NoIndex adds noindex headers for tenants that have prevent_indexing enabled
func NoIndex() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			tenant := c.Tenant()
			if tenant != nil && tenant.PreventIndexing {
				c.Response.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive, nosnippet, notranslate, noimageindex")
			}
			return next(c)
		}
	}
}