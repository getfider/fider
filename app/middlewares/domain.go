package middlewares

import (
	"net/http"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// Tenant adds either SingleTenant or MultiTenant to the pipeline
func Tenant() web.MiddlewareFunc {
	if env.IsSingleHostMode() {
		return SingleTenant()
	}
	return MultiTenant()
}

// SingleTenant inject default tenant into current context
func SingleTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			tenants := c.Services().Tenants
			tenant, err := tenants.First()
			if err != nil {
				if err == app.ErrNotFound {
					return c.Redirect(http.StatusTemporaryRedirect, "/signup")
				}
				return err
			}

			c.SetTenant(tenant)
			return next(c)
		}
	}
}

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			tenants := c.Services().Tenants
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

			c.Logger().Debugf("Tenant not found for '%s'.", hostname)
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
