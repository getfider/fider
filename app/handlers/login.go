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
			c.Logger().Errorf("Encoding claims failed with %s", err)
			return c.Failure(err)
		}

		var query = redirectURL.Query()
		query.Set("jwt", token)
		redirectURL.RawQuery = query.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
	}
}

// Login handles OAuth logins
func Login(provider string) web.HandlerFunc {
	return func(c web.Context) error {
		authURL := c.Services().OAuth.GetAuthURL(c.AuthEndpoint(), provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// LoginByEmail sends a new e-mail with verification key
func LoginByEmail() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.LoginByEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.SaveVerificationKey(input.Model.Email, input.Model.VerificationKey)
		if err != nil {
			return c.Failure(err)
		}

		subject := fmt.Sprintf("Log in to %s", c.Tenant().Name)
		message := fmt.Sprintf("%s/login/verify?k=%s", c.BaseURL(), input.Model.VerificationKey)
		err = c.Services().Emailer.Send(c.Tenant().Name, input.Model.Email, subject, message)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// VerifyLoginKey checks if verify key is correct and log in user
func VerifyLoginKey() web.HandlerFunc {
	return func(c web.Context) error {
		key := c.QueryParam("k")

		//TODO: If not found, go to main page
		//TODO: check if key didn't expire
		result, err := c.Services().Tenants.FindVerificationByKey(key)
		if err != nil {
			return c.Failure(err)
		}

		//TODO: If not found, request name and register user
		user, err := c.Services().Users.GetByEmail(c.Tenant().ID, result.Email)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.SetKeyAsVerified(key)
		if err != nil {
			return c.Failure(err)
		}

		//TODO: do something better with time
		claims := models.FiderClaims{
			UserID:    user.ID,
			UserName:  user.Name,
			UserEmail: user.Email,
		}

		token, err := jwt.Encode(claims)
		if err != nil {
			return c.Failure(err)
		}

		c.SetCookie(&http.Cookie{
			Name:     "auth",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
		})

		return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
	}
}

// Logout remove auth cookies
func Logout() web.HandlerFunc {
	return func(c web.Context) error {
		c.RemoveCookie("auth")
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
