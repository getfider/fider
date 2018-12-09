package middlewares

import (
	"net/http"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// Maintenance returns a maintenance page when system is under maintenance
func Maintenance() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if env.GetEnvOrDefault("MAINTENANCE", "") != "true" {
				return next(c)
			}

			c.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Response.Header().Set("Retry-After", "3600")

			message := env.GetEnvOrDefault("MAINTENANCE_MESSAGE", "Sorry, we're down for scheduled maintenance right now.")
			until := env.GetEnvOrDefault("MAINTENANCE_UNTIL", "")

			return c.Render(http.StatusServiceUnavailable, "maintenance.html", web.Props{
				Title:       "UNDER MAINTENANCE",
				Description: message,
				Data: web.Map{
					"message": message,
					"until":   until,
				},
			})
		}
	}
}
