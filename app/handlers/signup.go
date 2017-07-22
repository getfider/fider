package handlers

import (
	"net/http"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/im"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
)

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c web.Context) error {
		subdomain := strings.ToLower(c.Param("subdomain"))

		if result := validate.Subdomain(c.Services().Tenants, subdomain); !result.Ok {
			return c.Ok(web.Map{
				"message": strings.Join(result.Messages, ","),
			})
		}

		return c.Ok(web.Map{})
	}
}

//CreateTenant creates a new tenant
func CreateTenant() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(im.CreateTenant)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		tenant, err := c.Services().Tenants.Add(input.Name, input.Subdomain)
		if err != nil {
			return c.Failure(err)
		}

		user := &models.User{
			Name:   input.UserClaims.OAuthName,
			Email:  input.UserClaims.OAuthEmail,
			Tenant: tenant,
			Role:   models.RoleAdministrator,
			Providers: []*models.UserProvider{
				{UID: input.UserClaims.OAuthID, Name: input.UserClaims.OAuthProvider},
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
