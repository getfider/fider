package handlers

import (
	"net/http"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/validate"
	"github.com/labstack/echo"
)

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c web.Context) error {
		subdomain := c.Param("subdomain")
		ok, messages := validate.Subdomain(subdomain)
		if !ok {
			return c.Ok(echo.Map{
				"message": strings.Join(messages, ","),
			})
		}

		available, err := c.Services().Tenants.IsSubdomainAvailable(subdomain)
		if err != nil {
			return c.Failure(err)
		}

		if !available {
			return c.Ok(echo.Map{
				"message": "This subdomain is not available anymore",
			})
		}

		return c.Ok(echo.Map{})
	}
}

type createTenantInput struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
}

//CreateTenant creates a new tenant
func CreateTenant() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(createTenantInput)
		if err := c.Bind(input); err != nil {
			return c.Failure(err)
		}

		claims, err := jwt.DecodeOAuthClaims(input.Token)
		if err != nil {
			return c.Failure(err)
		}

		tenant := &models.Tenant{
			Name:      input.Name,
			Subdomain: input.Subdomain,
		}

		if env.IsSingleHostMode() {
			tenant.Subdomain = "Default"
		}

		if err := c.Services().Tenants.Add(tenant); err != nil {
			return c.Failure(err)
		}

		user := &models.User{
			Name:   claims.OAuthName,
			Email:  claims.OAuthEmail,
			Tenant: tenant,
			Role:   models.RoleAdministrator,
			Providers: []*models.UserProvider{
				{UID: claims.OAuthID, Name: claims.OAuthProvider},
			},
		}
		if err := c.Services().Users.Register(user); err != nil {
			return c.Failure(err)
		}

		return c.Ok(echo.Map{
			"id": tenant.ID,
		})
	}
}

//SignUp is the entry point for installation / signup
func SignUp() web.HandlerFunc {
	return func(c web.Context) error {
		if env.IsSingleHostMode() {
			tenant, err := c.Services().Tenants.First()
			if err != nil && err != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}
		}
		return c.Page(echo.Map{})
	}
}
