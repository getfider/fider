package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
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

		if c.Tenant().IsPrivate {
			return c.Page(http.StatusOK, web.Props{
				Page:  "SignIn/SignIn.page",
				Title: "Sign in",
			})
		}

		return c.Redirect(c.BaseURL())
	}
}

func LoginEmailSentPage() web.HandlerFunc {
	return func(c *web.Context) error {

		return c.Page(http.StatusOK, web.Props{
			Page:  "SignIn/LoginEmailSent.page",
			Title: "Login email sent",
			Data: web.Map{
				"email": c.QueryParam("email")},
		})

	}
}

// CompleteSignInProfilePage renders the complete profile page for code flow
func CompleteSignInProfilePage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "SignIn/CompleteSignInProfile.page",
			Title: "Complete your profile",
			Data: web.Map{
				"kind": enum.EmailVerificationKindSignIn,
				"k":    c.QueryParam("code"),
			},
		})
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

// SignInByEmail checks if user exists and sends code only for existing users
func SignInByEmail() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewSignInByEmail()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Check if user exists
		userByEmail := &query.GetUserByEmail{Email: action.Email}
		err := bus.Dispatch(c, userByEmail)
		userExists := err == nil

		// Only send code if user exists
		if userExists {
			err := bus.Dispatch(c, &cmd.SaveVerificationKey{
				Key:      action.VerificationCode,
				Duration: 15 * time.Minute,
				Request:  action,
			})
			if err != nil {
				return c.Failure(err)
			}

			c.Enqueue(tasks.SendSignInEmail(action.Email, action.VerificationCode))
		}

		return c.Ok(web.Map{
			"userExists": userExists,
		})
	}
}

// SignInByEmailWithName sends verification code for new users with their name
func SignInByEmailWithName() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewSignInByEmailWithName()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Check that user doesn't already exist
		userByEmail := &query.GetUserByEmail{Email: action.Email}
		err := bus.Dispatch(c, userByEmail)
		if err == nil {
			// User already exists, should use regular sign in
			return c.BadRequest(web.Map{
				"email": "An account with this email already exists. Please sign in.",
			})
		}

		// Check if tenant is private (new users not allowed)
		if c.Tenant().IsPrivate {
			return c.Forbidden()
		}

		// Save verification with name
		err = bus.Dispatch(c, &cmd.SaveVerificationKey{
			Key:      action.VerificationCode,
			Duration: 15 * time.Minute,
			Request:  action,
		})
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendSignInEmail(action.Email, action.VerificationCode))

		return c.Ok(web.Map{})
	}
}

// VerifySignInCode verifies the code entered by the user and signs them in
func VerifySignInCode() web.HandlerFunc {
	return func(c *web.Context) error {
		action := &actions.VerifySignInCode{}
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Get verification by email and code
		verification := &query.GetVerificationByEmailAndCode{
			Email: action.Email,
			Code:  action.Code,
			Kind:  enum.EmailVerificationKindSignIn,
		}
		err := bus.Dispatch(c, verification)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				return c.BadRequest(web.Map{
					"code": "Invalid or expired verification code",
				})
			}
			return c.Failure(err)
		}

		result := verification.Result

		// Check if already verified (with grace period)
		if result.VerifiedAt != nil {
			if time.Since(*result.VerifiedAt) > 5*time.Minute {
				return c.Gone()
			}
		} else {
			// Check if expired
			if time.Now().After(result.ExpiresAt) {
				// Mark as verified to prevent reuse
				_ = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: action.Code})
				return c.Gone()
			}
		}

		// Check if user exists
		userByEmail := &query.GetUserByEmail{Email: result.Email}
		err = bus.Dispatch(c, userByEmail)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				// User doesn't exist
				if c.Tenant().IsPrivate {
					return c.Forbidden()
				}

				// If name is provided in verification, create user account
				if result.Name != "" {
					user := &entity.User{
						Name:   result.Name,
						Email:  result.Email,
						Tenant: c.Tenant(),
						Role:   enum.RoleVisitor,
					}
					err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
					if err != nil {
						return c.Failure(err)
					}

					// Mark code as verified
					err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: action.Code})
					if err != nil {
						return c.Failure(err)
					}

					// Authenticate newly created user
					webutil.AddAuthUserCookie(c, user)
					return c.Ok(web.Map{})
				}

				// Name not provided - shouldn't happen with new flow, but handle legacy case
				return c.Ok(web.Map{
					"showProfileCompletion": true,
				})
			}
			return c.Failure(err)
		}

		// Mark code as verified
		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: action.Code})
		if err != nil {
			return c.Failure(err)
		}

		// Authenticate user
		webutil.AddAuthUserCookie(c, userByEmail.Result)

		return c.Ok(web.Map{})
	}
}

// ResendSignInCode invalidates the previous code and sends a new one
func ResendSignInCode() web.HandlerFunc {
	return func(c *web.Context) error {
		// Create new sign-in action with new code
		action := actions.NewSignInByEmail()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Authorization check
		if !action.IsAuthorized(c, c.User()) {
			return c.Forbidden()
		}

		// Save new verification code
		err := bus.Dispatch(c, &cmd.SaveVerificationKey{
			Key:      action.VerificationCode,
			Duration: 15 * time.Minute,
			Request:  action,
		})
		if err != nil {
			return c.Failure(err)
		}

		// Send new email
		c.Enqueue(tasks.SendSignInEmail(action.Email, action.VerificationCode))

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

				// If name is provided in verification, create user account immediately
				if result.Name != "" {
					user := &entity.User{
						Name:   result.Name,
						Email:  result.Email,
						Tenant: c.Tenant(),
						Role:   enum.RoleVisitor,
					}
					err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
					if err != nil {
						return c.Failure(err)
					}

					err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: key})
					if err != nil {
						return c.Failure(err)
					}

					webutil.AddAuthUserCookie(c, user)
					baseURL := c.BaseURL()
					return c.Redirect(baseURL)
				}

				// Otherwise, show profile completion page
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

		baseURL := c.BaseURL()
		return c.Redirect(baseURL)
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
		return c.Redirect("/")
	}
}
