package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
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

			users := c.Services().Users

			var user *models.User
			if oauthUser.Email == "" {
				user, err = users.GetByProvider(tenant.ID, provider, oauthUser.ID.String())
			} else {
				user, err = users.GetByEmail(tenant.ID, oauthUser.Email)
			}
			if err != nil {
				if err == app.ErrNotFound {
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
		query.Set("jwt", token)
		redirectURL.RawQuery = query.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
	}
}

// SignIn handles OAuth sign in
func SignIn(provider string) web.HandlerFunc {
	return func(c web.Context) error {
		authURL := c.Services().OAuth.GetAuthURL(c.AuthEndpoint(), provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// SignInByEmail sends a new e-mail with verification key
func SignInByEmail() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.SignInByEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.SaveVerificationKey(input.Model.Email, input.Model.VerificationKey)
		if err != nil {
			return c.Failure(err)
		}

		subject := fmt.Sprintf("Sign in to %s", c.Tenant().Name)
		link := fmt.Sprintf("%s/signin/verify?k=%s", c.BaseURL(), input.Model.VerificationKey)
		message := fmt.Sprintf("Click and confirm that you want to sign in. This link will expire in 15 minutes and can only be used once. <br /><br /> <a href='%s'>%s</a>", link, link)
		err = c.Services().Emailer.Send(c.Tenant().Name, input.Model.Email, subject, message)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// VerifySignInKey checks if verify key is correct and sign in user
func VerifySignInKey() web.HandlerFunc {
	return func(c web.Context) error {
		key := c.QueryParam("k")

		result, err := c.Services().Tenants.FindVerificationByKey(key)
		if err != nil {
			if err == app.ErrNotFound {
				return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
			}
			return c.Failure(err)
		}

		//If key has been used, just go back to home page
		if result.VerifiedOn != nil {
			return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}

		//If key expired (15 minutes), go back to home page
		if time.Now().After(result.CreatedOn.Add(15 * time.Minute)) {
			err = c.Services().Tenants.SetKeyAsVerified(key)
			if err != nil {
				return c.Failure(err)
			}
			return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}

		user, err := c.Services().Users.GetByEmail(c.Tenant().ID, result.Email)
		if err != nil {
			if err == app.ErrNotFound {
				// This will render a page for /signin/verify URL using same variables as home page
				return Index()(c)
			}
			return c.Failure(err)
		}

		err = c.Services().Tenants.SetKeyAsVerified(key)
		if err != nil {
			return c.Failure(err)
		}

		err = c.AddAuthCookie(user)
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

		_, err := c.Services().Users.GetByEmail(c.Tenant().ID, input.Model.Email)
		if err != app.ErrNotFound {
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

		err = c.AddAuthCookie(user)
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
