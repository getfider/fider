package middlewares

import (
	"fmt"

	"github.com/getfider/fider/app/pkg/web"
)

// Secure adds web security related Http Headers to response
func Secure() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			c.Response.Header().Add("Content-Security-Policy-Report-Only", fmt.Sprintf(web.CspPolicyTemplate, c.ContextID()))
			c.Response.Header().Add("X-XSS-Protection", "1; mode=block")
			c.Response.Header().Add("X-Content-Type-Options", "nosniff")
			c.Response.Header().Add("Referrer-Policy", "no-referrer-when-downgrade")
			return next(c)
		}
	}
}
