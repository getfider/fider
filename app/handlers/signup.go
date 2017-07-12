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
			return c.Ok(web.Map{
				"message": strings.Join(messages, ","),
			})
		}

		return c.Ok(web.Map{})
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
			return c.BadRequest(web.Map{
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

		if input.Name == "" {
			return c.BadRequest(web.Map{
				"message": "Name is required",
			})
		}

		ok, messages, err := validate.Subdomain(c.Services().Tenants, input.Subdomain)
		if err != nil {
			return c.Failure(err)
		}

		if !ok {
			return c.BadRequest(web.Map{
				"message": strings.Join(messages, ","),
			})
		}

		tenant, err := c.Services().Tenants.Add(input.Name, input.Subdomain)
		if err != nil {
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

		return c.Ok(web.Map{
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
		return c.Page(web.Map{})
	}
}
