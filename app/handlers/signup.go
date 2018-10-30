package handlers

import (
	"time"

	"github.com/getfider/fider/app/pkg/web/util"

	"github.com/getfider/fider/app/tasks"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
)

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c web.Context) error {
		subdomain := strings.ToLower(c.Param("subdomain"))

		messages, _ := validate.Subdomain(c.Services().Tenants, subdomain)
		if len(messages) > 0 {
			return c.Ok(web.Map{
				"message": strings.Join(messages, ","),
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

		status := models.TenantPending
		if socialSignUp {
			status = models.TenantActive
		}

		tenant, err := c.Services().Tenants.Add(input.Model.TenantName, input.Model.Subdomain, status)
		if err != nil {
			return c.Failure(err)
		}

		c.SetTenant(tenant)

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

			if err := c.Services().Users.Register(user); err != nil {
				return c.Failure(err)
			}

			if env.IsSingleHostMode() {
				webutil.AddAuthUserCookie(c, user)
			} else {
				webutil.SetSignUpAuthCookie(c, user)
			}

		} else {
			user.Name = input.Model.Name
			user.Email = input.Model.Email

			err := c.Services().Tenants.SaveVerificationKey(input.Model.VerificationKey, 48*time.Hour, input.Model)
			if err != nil {
				return c.Failure(err)
			}

			c.Enqueue(tasks.SendSignUpEmail(input.Model, c.TenantBaseURL(tenant)))
		}

		return c.Ok(web.Map{})
	}
}

//SignUp is the entry point for installation / signup
func SignUp() web.HandlerFunc {
	return func(c web.Context) error {
		if env.IsSingleHostMode() {
			tenant, err := c.Services().Tenants.First()
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				return c.Redirect("/")
			}
		} else {
			baseURL := webutil.GetOAuthBaseURL(c)
			if !strings.HasPrefix(c.Request.URL.String(), baseURL) {
				return c.Redirect(baseURL + "/signup")
			}
		}

		return c.Page(web.Props{
			Title:       "Sign up",
			Description: "Sign up for Fider and let your customers share, vote and discuss on suggestions they have to make your product even better.",
		})
	}
}

// VerifySignUpKey checks if verify key is correct, activate the tenant and sign in user
func VerifySignUpKey() web.HandlerFunc {
	return func(c web.Context) error {
		if c.Tenant().Status == models.TenantPending {
			return VerifySignInKey(models.EmailVerificationKindSignUp)(c)
		}
		return c.NotFound()
	}
}
