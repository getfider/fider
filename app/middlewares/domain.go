package middlewares

import (
	"net/http"
	"net/url"

	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage"
)

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(tenants storage.Tenant) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			hostname := stripPort(c.Request().Host)
			tenant, err := tenants.GetByDomain(hostname)
			if err == nil {
				c.SetTenant(tenant)
				return next(c)
			}

			c.Logger().Infof("Tenant not found for '%s'.", hostname)
			return c.NoContent(http.StatusNotFound)
		}
	}
}

// HostChecker checks for a specific host
func HostChecker(baseURL string) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			u, err := url.Parse(baseURL)
			if err != nil {
				return c.Failure(err)
			}

			if c.Request().Host != u.Host {
				c.Logger().Errorf("%s is not valid for this operation. Only %s is allowed.", c.Request().Host, u.Host)
				return c.NoContent(http.StatusBadRequest)
			}

			return next(c)
		}
	}
}
