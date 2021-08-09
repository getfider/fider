package middlewares

import (
	"net/http"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// Maintenance returns a maintenance page when system is under maintenance
func Maintenance() web.MiddlewareFunc {
	if !env.Config.Maintenance.Enabled {
		return nil
	}

	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {

			c.Response.Header().Set("Retry-After", "3600")

			return c.Render(http.StatusServiceUnavailable, "maintenance.html", web.Props{
				Title:       "UNDER MAINTENANCE",
				Description: env.Config.Maintenance.Message,
				Data: web.Map{
					"message": env.Config.Maintenance.Message,
					"until":   env.Config.Maintenance.Until,
				},
			})
		}
	}
}
