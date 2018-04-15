package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

type oauthUserProfile struct {
	Name  string
	ID    string
	Email string
}

// OAuthCallback handles OAuth callbacks
func OAuthCallback(provider string) web.HandlerFunc {
	return func(c web.Context) error {

		redirect := c.QueryParam("state")
		redirectURL, err := url.ParseRequestURI(redirect)
		if err != nil {
			return c.Failure(err)
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(http.StatusTemporaryRedirect, redirect)
		}

		oauthUser, err := c.Services().OAuth.GetProfile(c.AuthEndpoint(), provider, code)
		if err != nil {
			return c.Failure(err)
		}

		var claims jwtgo.Claims
		if redirectURL.Path != "/signup" {

			var (
				tenant *models.Tenant
				err    error
			)
			if env.IsSingleHostMode() {
				tenant, err = c.Services().Tenants.First()
			} else {
				tenant, err = c.Services().Tenants.GetByDomain(stripPort(redirectURL.Host))
			}
			if err != nil {
				return c.Failure(err)
			}

			c.Services().SetCurrentTenant(tenant)
			users := c.Services().Users

			user, err := users.GetByProvider(provider, oauthUser.ID.String())
			if errors.Cause(err) == app.ErrNotFound && oauthUser.Email != "" {
				user, err = users.GetByEmail(oauthUser.Email)
			}
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					if tenant.IsPrivate {
						return c.Redirect(http.StatusTemporaryRedirect, c.TenantBaseURL(tenant)+"/not-invited")
					}

					user = &models.User{
						Name:   oauthUser.Name,
						Tenant: tenant,
						Email:  oauthUser.Email,
						Role:   models.RoleVisitor,
						Providers: []*models.UserProvider{
							&models.UserProvider{
								UID:  oauthUser.ID.String(),
								Name: provider,
							},
						},
					}

					err = users.Register(user)
					if err != nil {
						return c.Failure(err)
					}
				} else {
					return c.Failure(err)
				}
			} else if !user.HasProvider(provider) {
				err = users.RegisterProvider(user.ID, &models.UserProvider{
					UID:  oauthUser.ID.String(),
					Name: provider,
				})
				if err != nil {
					return c.Failure(err)
				}
			}

			claims = models.FiderClaims{
				UserID:    user.ID,
				UserName:  user.Name,
				UserEmail: user.Email,
			}
		} else {
			claims = models.OAuthClaims{
				OAuthID:       oauthUser.ID.String(),
				OAuthProvider: provider,
				OAuthName:     oauthUser.Name,
				OAuthEmail:    oauthUser.Email,
			}
		}

		var token string
		if token, err = jwt.Encode(claims); err != nil {
			return c.Failure(err)
		}

		var query = redirectURL.Query()
		query.Set("token", token)
		redirectURL.RawQuery = query.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
	}
}

// SignInPage renders the sign in page
func SignInPage() web.HandlerFunc {
	return func(c web.Context) error {
		if c.IsAuthenticated() || !c.Tenant().IsPrivate {
			return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}
		return c.Page(web.Map{})
	}
}

// NotInvitedPage renders the not invited page
func NotInvitedPage() web.HandlerFunc {
	return func(c web.Context) error {
		return c.Render(http.StatusForbidden, "not-invited.html", web.Map{})
	}
}

// SignInByOAuth handles OAuth sign in
func SignInByOAuth(provider string) web.HandlerFunc {
	return func(c web.Context) error {
		authURL := c.Services().OAuth.GetAuthURL(c.AuthEndpoint(), provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// SignInByEmail sends a new email with verification key
func SignInByEmail() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.SignInByEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.SaveVerificationKey(input.Model.VerificationKey, 15*time.Minute, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendSignInEmail(input.Model))

		return c.Ok(web.Map{})
	}
}

// VerifySignInKey checks if verify key is correct and sign in user
func VerifySignInKey(kind models.EmailVerificationKind) web.HandlerFunc {
	return func(c web.Context) error {
		result, err := validateKey(kind, c)
		if result == nil {
			return err
		}

		var user *models.User
		if kind == models.EmailVerificationKindSignUp && c.Tenant().Status == models.TenantInactive {
			if err = c.Services().Tenants.Activate(c.Tenant().ID); err != nil {
				return c.Failure(err)
			}

			user = &models.User{
				Name:   result.Name,
				Email:  result.Email,
				Tenant: c.Tenant(),
				Role:   models.RoleAdministrator,
			}

			if err = c.Services().Users.Register(user); err != nil {
				return c.Failure(err)
			}
		} else if kind == models.EmailVerificationKindSignIn {
			user, err = c.Services().Users.GetByEmail(result.Email)
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					if c.Tenant().IsPrivate {
						return NotInvitedPage()(c)
					}

					// This will render a page for /signin/verify URL using same variables as home page
					return Index()(c)
				}
				return c.Failure(err)
			}
		} else if kind == models.EmailVerificationKindUserInvitation {
			user, err = c.Services().Users.GetByEmail(result.Email)
			if err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					// This will render a page for /invite/verify URL using same variables as sign in page
					return SignInPage()(c)
				}
				return c.Failure(err)
			}
		} else {
			return c.NotFound()
		}

		err = c.Services().Tenants.SetKeyAsVerified(result.Key)
		if err != nil {
			return c.Failure(err)
		}

		_, err = c.AddAuthCookie(user)
		if err != nil {
			return c.Failure(err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
	}
}

// CompleteSignInProfile handles the action to update user profile
func CompleteSignInProfile() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CompleteProfile)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		_, err := c.Services().Users.GetByEmail(input.Model.Email)
		if errors.Cause(err) != app.ErrNotFound {
			return c.Ok(web.Map{})
		}

		user := &models.User{
			Name:   input.Model.Name,
			Email:  input.Model.Email,
			Tenant: c.Tenant(),
			Role:   models.RoleVisitor,
		}
		err = c.Services().Users.Register(user)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.SetKeyAsVerified(input.Model.Key)
		if err != nil {
			return c.Failure(err)
		}

		_, err = c.AddAuthCookie(user)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// SignOut remove auth cookies
func SignOut() web.HandlerFunc {
	return func(c web.Context) error {
		c.RemoveCookie(web.CookieAuthName)
		return c.Redirect(http.StatusTemporaryRedirect, c.QueryParam("redirect"))
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
