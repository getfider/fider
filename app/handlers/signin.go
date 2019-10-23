package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
	"github.com/getfider/fider/app/tasks"
)

// SignInPage renders the sign in page
func SignInPage() web.HandlerFunc {
	return func(c *web.Context) error {

		if c.Tenant().IsPrivate || c.Tenant().Status == enum.TenantLocked {
			return c.Page(web.Props{
				Title:     "Sign in",
				ChunkName: "SignIn.page",
			})
		}

		return c.Redirect(c.BaseURL())
	}
}

// NotInvitedPage renders the not invited page
func NotInvitedPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Render(http.StatusForbidden, "not-invited.html", web.Props{
			Title:       "Not Invited",
			Description: "We couldn't find your account for your email address.",
		})
	}
}

// SignInByEmail sends a new email with verification key
func SignInByEmail() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.SignInByEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.SaveVerificationKey{
			Key:      input.Model.VerificationKey,
			Duration: 30 * time.Minute,
			Request:  input.Model,
		})
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendSignInEmail(input.Model))

		return c.Ok(web.Map{})
	}
}

// VerifySignInKey checks if verify key is correct and sign in user
func VerifySignInKey(kind enum.EmailVerificationKind) web.HandlerFunc {
	return func(c *web.Context) error {
		result, err := validateKey(kind, c)
		if result == nil {
			return err
		}

		var user *models.User
		if kind == enum.EmailVerificationKindSignUp && c.Tenant().Status == enum.TenantPending {
			if err = bus.Dispatch(c, &cmd.ActivateTenant{TenantID: c.Tenant().ID}); err != nil {
				return c.Failure(err)
			}

			user = &models.User{
				Name:   result.Name,
				Email:  result.Email,
				Tenant: c.Tenant(),
				Role:   enum.RoleAdministrator,
			}

			if err = bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
				return c.Failure(err)
			}
		} else if kind == enum.EmailVerificationKindSignIn {
			userByEmail := &query.GetUserByEmail{Email: result.Email}
			err = bus.Dispatch(c, userByEmail)
			user = userByEmail.Result
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					if c.Tenant().IsPrivate {
						return NotInvitedPage()(c)
					}
					return Index()(c)
				}
				return c.Failure(err)
			}
		} else if kind == enum.EmailVerificationKindUserInvitation {
			userByEmail := &query.GetUserByEmail{Email: result.Email}
			err = bus.Dispatch(c, userByEmail)
			user = userByEmail.Result
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					if c.Tenant().IsPrivate {
						return SignInPage()(c)
					}
					return Index()(c)
				}
				return c.Failure(err)
			}
		} else {
			return c.NotFound()
		}

		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: result.Key})
		if err != nil {
			return c.Failure(err)
		}

		webutil.AddAuthUserCookie(c, user)

		return c.Redirect(c.BaseURL())
	}
}

// CompleteSignInProfile handles the action to update user profile
func CompleteSignInProfile() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CompleteProfile)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &query.GetUserByEmail{Email: input.Model.Email})
		if errors.Cause(err) != app.ErrNotFound {
			return c.Ok(web.Map{})
		}

		user := &models.User{
			Name:   input.Model.Name,
			Email:  input.Model.Email,
			Tenant: c.Tenant(),
			Role:   enum.RoleVisitor,
		}
		err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: input.Model.Key})
		if err != nil {
			return c.Failure(err)
		}

		webutil.AddAuthUserCookie(c, user)

		return c.Ok(web.Map{})
	}
}

// SignOut remove auth cookies
func SignOut() web.HandlerFunc {
	return func(c *web.Context) error {
		c.RemoveCookie(web.CookieAuthName)
		return c.Redirect(c.QueryParam("redirect"))
	}
}
