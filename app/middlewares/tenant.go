package middlewares

import (
	"fmt"
	"net/http"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
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
			tenant, err := c.Services().Tenants.First()
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				c.SetTenant(tenant)
			}

			return next(c)
		}
	}
}

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			hostname := c.Request.URL.Hostname()

			// If no tenant is specified, redirect user to getfider.com
			// This is only valid for fider.io hosting
			if (env.IsProduction() && hostname == "fider.io") ||
				(env.IsDevelopment() && hostname == "dev.fider.io") {
				if c.Request.URL.Path == "" || c.Request.URL.Path == "/" {
					return c.Redirect("https://getfider.com")
				}
			}

			tenant, err := c.Services().Tenants.GetByDomain(hostname)
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				c.SetTenant(tenant)

				if tenant.CNAME != "" && !c.IsAjax() {
					baseURL := c.TenantBaseURL(tenant)
					if baseURL != c.BaseURL() {
						link := baseURL + c.Request.URL.RequestURI()
						c.Response.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"canonical\"", link))
					}
				}
			}

			return next(c)
		}
	}
}

// RequireTenant returns 404 if tenant is not available
func RequireTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant() == nil {
				if env.IsSingleHostMode() {
					return c.Redirect("/signup")
				}
				return c.NotFound()
			}
			return next(c)
		}
	}
}

// BlockPendingTenants blocks requests for pending tenants
func BlockPendingTenants() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant().Status == models.TenantPending {
				return c.Render(http.StatusOK, "pending-activation.html", web.Props{
					Title:       "Pending Activation",
					Description: "We sent you a confirmation email with a link to activate your site. Please check your inbox to activate it.",
				})
			}
			return next(c)
		}
	}
}

// CheckTenantPrivacy blocks requests of unauthenticated users for private tenants
func CheckTenantPrivacy() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant().IsPrivate && !c.IsAuthenticated() {
				return c.Redirect("/signin")
			}
			return next(c)
		}
	}
}
