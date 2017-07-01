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
		subdomain := strings.ToLower(c.Param("subdomain"))
		ok, messages, err := validate.Subdomain(c.Services().Tenants, subdomain)
		if err != nil {
			return c.Failure(err)
		}

		if !ok {
			return c.Ok(echo.Map{
				"message": strings.Join(messages, ","),
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

		if input.Token == "" {
			return c.BadRequest(echo.Map{
				"message": "Please identify yourself before proceeding",
			})
		}

		claims, err := jwt.DecodeOAuthClaims(input.Token)
		if err != nil {
			return c.Failure(err)
		}

		if env.IsSingleHostMode() {
			input.Subdomain = "default"
		}

		tenant := &models.Tenant{
			Name:      input.Name,
			Subdomain: input.Subdomain,
		}

		tenant.Subdomain = strings.ToLower(input.Subdomain)

		if tenant.Name == "" {
			return c.BadRequest(echo.Map{
				"message": "Name is required",
			})
		}

		ok, messages, err := validate.Subdomain(c.Services().Tenants, tenant.Subdomain)
		if err != nil {
			return c.Failure(err)
		}

		if !ok {
			return c.BadRequest(echo.Map{
				"message": strings.Join(messages, ","),
			})
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
