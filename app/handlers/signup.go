package handlers

import (
	"time"

	"github.com/getfider/fider/app/tasks"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
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

			expiresAt := time.Now().Add(365 * 24 * time.Hour)
			token, err := jwt.Encode(jwt.FiderClaims{
				UserID:    user.ID,
				UserName:  user.Name,
				UserEmail: user.Email,
				Metadata: jwt.Metadata{
					ExpiresAt: expiresAt.Unix(),
				},
			})

			if err != nil {
				return c.Failure(err)
			}

			if env.IsSingleHostMode() {
				c.AddCookie(web.CookieAuthName, token, expiresAt)
			} else {
				c.AddDomainCookie(web.CookieSignUpAuthName, token)
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
		if c.Tenant().Status == models.TenantInactive {
			return VerifySignInKey(models.EmailVerificationKindSignUp)(c)
		}
		return c.NotFound()
	}
}
