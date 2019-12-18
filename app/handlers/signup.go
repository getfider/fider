package handlers

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"

	webutil "github.com/getfider/fider/app/pkg/web/util"

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
	return func(c *web.Context) error {
		subdomain := strings.ToLower(c.Param("subdomain"))

		messages, _ := validate.Subdomain(c, subdomain)
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
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		input := new(actions.CreateTenant)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		socialSignUp := input.Model.Token != ""

		status := enum.TenantPending
		if socialSignUp {
			status = enum.TenantActive
		}

		createTenant := &cmd.CreateTenant{
			Name:      input.Model.TenantName,
			Subdomain: input.Model.Subdomain,
			Status:    status,
		}
		err := bus.Dispatch(c, createTenant)
		if err != nil {
			return c.Failure(err)
		}

		c.SetTenant(createTenant.Result)

		user := &models.User{
			Tenant: createTenant.Result,
			Role:   enum.RoleAdministrator,
		}

		if socialSignUp {
			user.Name = input.Model.UserClaims.OAuthName
			user.Email = input.Model.UserClaims.OAuthEmail
			user.Providers = []*models.UserProvider{
				{UID: input.Model.UserClaims.OAuthID, Name: input.Model.UserClaims.OAuthProvider},
			}

			if err := bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
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

			err := bus.Dispatch(c, &cmd.SaveVerificationKey{
				Key:      input.Model.VerificationKey,
				Duration: 48 * time.Hour,
				Request:  input.Model,
			})
			if err != nil {
				return c.Failure(err)
			}

			c.Enqueue(tasks.SendSignUpEmail(input.Model, web.TenantBaseURL(c, createTenant.Result)))
		}

		return c.Ok(web.Map{})
	}
}

//SignUp is the entry point for installation / signup
func SignUp() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		if env.IsSingleHostMode() {
			firstTenant := &query.GetFirstTenant{}
			err := bus.Dispatch(c, firstTenant)
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if firstTenant.Result != nil {
				return c.Redirect("/")
			}
		} else {
			baseURL := web.OAuthBaseURL(c)
			if !strings.HasPrefix(c.Request.URL.String(), baseURL) {
				return c.Redirect(baseURL + "/signup")
			}
		}

		return c.Page(web.Props{
			Title:       "Sign up",
			Description: "Sign up for Fider and let your customers share, vote and discuss on suggestions they have to make your product even better.",
			ChunkName:   "SignUp.page",
		})
	}
}

// VerifySignUpKey checks if verify key is correct, activate the tenant and sign in user
func VerifySignUpKey() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Tenant().Status == enum.TenantPending {
			return VerifySignInKey(enum.EmailVerificationKindSignUp)(c)
		}
		return c.NotFound()
	}
}
