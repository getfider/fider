package handlers

import (
	"net/http"
	"time"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
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
		input := new(actions.CreateTenant)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		socialSignUp := input.Model.Token != ""

		status := models.TenantInactive
		if socialSignUp {
			status = models.TenantActive
		}

		tenant, err := c.Services().Tenants.Add(input.Model.TenantName, input.Model.Subdomain, status)
		if err != nil {
			return c.Failure(err)
		}

		c.Services().Tenants.SetCurrentTenant(tenant)

		user := &models.User{
			Tenant: tenant,
			Role:   models.RoleAdministrator,
		}

		if socialSignUp {
			user.Name = input.Model.UserClaims.OAuthName
			user.Email = input.Model.UserClaims.OAuthEmail
			user.Providers = []*models.UserProvider{
				{UID: input.Model.UserClaims.OAuthID, Name: input.Model.UserClaims.OAuthProvider},
			}
		} else {
			user.Name = input.Model.Name
			user.Email = input.Model.Email

			err := c.Services().Tenants.SaveVerificationKey(input.Model.VerificationKey, 48*time.Hour, input.Model)
			if err != nil {
				return c.Failure(err)
			}

			err = c.Services().Emailer.Send("Fider", user.Email, "Confirm your new Fider instance", "signup_email", web.Map{
				"baseUrl":         c.TenantBaseURL(tenant),
				"verificationKey": input.Model.VerificationKey,
			})
			if err != nil {
				return c.Failure(err)
			}
		}

		if socialSignUp {
			if err := c.Services().Users.Register(user); err != nil {
				return c.Failure(err)
			}
			token, err := c.AddAuthCookie(user)
			if err != nil {
				return c.Failure(err)
			}
			return c.Ok(web.Map{
				"token": token,
			})
		}

		return c.Ok(web.Map{})
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

// VerifySignUpKey checks if verify key is correct, activate the tenant and sign in user
func VerifySignUpKey() web.HandlerFunc {
	return func(c web.Context) error {
		if c.Tenant().Status == models.TenantInactive {
			return VerifySignInKey(models.EmailVerificationKindSignUp)(c)
		}
		return c.NotFound()
	}
}
