package router

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(ctx *context.WchyContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hostname := stripPort(c.Request().Host)
			tenant, err := ctx.Tenant.GetByDomain(hostname)
			if err == nil {
				c.Set("Tenant", tenant)
				return next(c)
			}

			fmt.Printf("Tenant not found for %s.\n", hostname)
			return c.NoContent(http.StatusNotFound)
		}
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
