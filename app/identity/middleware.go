package identity

import (
	"net/http"

	"github.com/WeCanHearYou/wechy/app"
)

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(tenantService TenantService) app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(c app.Context) error {
			hostname := stripPort(c.Request().Host)
			tenant, err := tenantService.GetByDomain(hostname)
			if err == nil {
				c.SetTenant(tenant)
				return next(c)
			}

			c.Logger().Infof("Tenant not found for '%s'.", hostname)
			return c.NoContent(http.StatusNotFound)
		}
	}
}
