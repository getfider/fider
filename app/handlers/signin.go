package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"

	"github.com/Spicy-Bush/fider-tarkov-community/app"
	"github.com/Spicy-Bush/fider-tarkov-community/app/actions"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/errors"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	webutil "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web/util"
	"github.com/Spicy-Bush/fider-tarkov-community/app/tasks"
)

// SignInPage renders the sign in page
func SignInPage() web.HandlerFunc {
	return func(c *web.Context) error {

		if c.Tenant().IsPrivate {
			return c.Page(http.StatusOK, web.Props{
				Page:  "SignIn/SignIn.page",
				Title: "Sign in",
			})
		}

		return c.Redirect(c.BaseURL())
	}
}

// NotInvitedPage renders the not invited page
func NotInvitedPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusForbidden, web.Props{
			Page:        "Error/NotInvited.page",
			Title:       "Not Invited",
			Description: "We couldn't find an account for your email address.",
		})
	}
}

// SignInByEmail sends a new email with verification key
func SignInByEmail() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewSignInByEmail()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.SaveVerificationKey{
			Key:      action.VerificationKey,
			Duration: 30 * time.Minute,
			Request:  action,
		})
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendSignInEmail(action.Email, action.VerificationKey))

		return c.Ok(web.Map{})
	}
}

// VerifySignInKey checks if verify key is correct and sign in user
func VerifySignInKey(kind enum.EmailVerificationKind) web.HandlerFunc {
	return func(c *web.Context) error {
		key := c.QueryParam("k")
		result, err := validateKey(kind, key, c)
		if result == nil {
			return err
		}

		userByEmail := &query.GetUserByEmail{Email: result.Email}
		err = bus.Dispatch(c, userByEmail)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				if kind == enum.EmailVerificationKindSignIn && c.Tenant().IsPrivate {
					return NotInvitedPage()(c)
				}

				return c.Page(http.StatusOK, web.Props{
					Page:  "SignIn/CompleteSignInProfile.page",
					Title: "Complete Sign In Profile",
					Data: web.Map{
						"kind": kind,
						"k":    key,
					},
				})
			}
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: key})
		if err != nil {
			return c.Failure(err)
		}

		webutil.AddAuthUserCookie(c, userByEmail.Result)

		return c.Redirect(c.BaseURL())
	}
}

// CompleteSignInProfile handles the action to update user profile
func CompleteSignInProfile() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.CompleteProfile)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		result, err := validateKey(action.Kind, action.Key, c)
		if result == nil {
			return err
		}

		err = bus.Dispatch(c, &query.GetUserByEmail{Email: result.Email})
		if errors.Cause(err) != app.ErrNotFound {
			// Not possible to create user that already exists
			return c.BadRequest(web.Map{})
		}

		user := &entity.User{
			Name:   action.Name,
			Email:  result.Email,
			Tenant: c.Tenant(),
			Role:   enum.RoleVisitor,
		}
		err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: action.Key})
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
		redirect := c.QueryParam("redirect")

		u, err := url.Parse(redirect)
		// Check if parsing failed, or if the URL is absolute (has a scheme) or includes a host.
		if err != nil || u.IsAbs() || u.Host != "" {
			redirect = "/"
		} else if !strings.HasPrefix(u.Path, "/") {
			// Ensure that relative paths always start with "/"
			redirect = "/" + u.Path
		}

		return c.Redirect(redirect)
	}
}
