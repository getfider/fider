package middlewares

import (
	"net/http"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage"
)

// Tenant adds either SingleTenant or MultiTenant to the pipeline
func Tenant(tenants storage.Tenant) web.MiddlewareFunc {
	if env.IsSingleHostMode() {
		return SingleTenant(tenants)
	}
	return MultiTenant(tenants)
}

// SingleTenant inject default tenant into current context
func SingleTenant(tenants storage.Tenant) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			tenant, err := tenants.First()
			if err == app.ErrNotFound {
				tenant = &models.Tenant{
					Name:      "Default",
					Subdomain: "default",
				}
				tenants.Add(tenant)
			}

			c.SetTenant(tenant)
			return next(c)
		}
	}
}

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant(tenants storage.Tenant) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			hostname := stripPort(c.Request().Host)
			u, err := url.Parse(c.AuthEndpoint())
			if err != nil {
				return c.Failure(err)
			}

			if hostname == stripPort(u.Host) {
				return next(c)
			}

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
			u, err := url.Parse("http://" + baseURL)
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
