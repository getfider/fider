package handlers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/web/util"
)

// OAuthEcho exchanges OAuth Code for a user profile and return directly to the UI, without storing it
func OAuthEcho() web.HandlerFunc {
	return func(c web.Context) error {
		provider := c.Param("provider")

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect("/")
		}

		identifier := c.QueryParam("identifier")
		if identifier == "" || identifier != c.SessionID() {
			c.Logger().Warn("OAuth identifier doesn't match with user session ID. Aborting sign in process.")
			return c.Redirect("/")
		}

		body, err := c.Services().OAuth.GetRawProfile(provider, code)
		if err != nil {
			return c.Page(web.Props{
				Title: "OAuth Test Page",
				Data: web.Map{
					"err": errors.Cause(err).Error(),
				},
			})
		}

		profile, _ := c.Services().OAuth.ParseRawProfile(provider, body)

		return c.Page(web.Props{
			Title: "OAuth Test Page",
			Data: web.Map{
				"body":    body,
				"profile": profile,
			},
		})
	}
}

// OAuthToken exchanges OAuth Code for a user profile
// The user profile is then used to either get an existing user on Fider or creating a new one
// Once Fider user is retrieved/created, an authentication cookie is store in user's browser
func OAuthToken() web.HandlerFunc {
	return func(c web.Context) error {
		provider := c.Param("provider")
		redirectURL, _ := url.ParseRequestURI(c.QueryParam("redirect"))
		redirectURL.ResolveReference(c.Request.URL)

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		identifier := c.QueryParam("identifier")
		if identifier == "" || identifier != c.SessionID() {
			c.Logger().Warn("OAuth identifier doesn't match with user session ID. Aborting sign in process.")
			return c.Redirect(redirectURL.String())
		}

		oauthUser, err := c.Services().OAuth.GetProfile(provider, code)
		if err != nil {
			return c.Failure(err)
		}

		users := c.Services().Users

		user, err := users.GetByProvider(provider, oauthUser.ID)
		if errors.Cause(err) == app.ErrNotFound && oauthUser.Email != "" {
			user, err = users.GetByEmail(oauthUser.Email)
		}
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				if c.Tenant().IsPrivate {
					return c.Redirect("/not-invited")
				}

				user = &models.User{
					Name:   oauthUser.Name,
					Tenant: c.Tenant(),
					Email:  oauthUser.Email,
					Role:   models.RoleVisitor,
					Providers: []*models.UserProvider{
						&models.UserProvider{
							UID:  oauthUser.ID,
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
				UID:  oauthUser.ID,
				Name: provider,
			})
			if err != nil {
				return c.Failure(err)
			}
		}

		webutil.AddAuthUserCookie(c, user)

		return c.Redirect(redirectURL.String())
	}
}

// OAuthCallback handles the redirect back from the OAuth provider
// This callback can run on either Tenant or Login address
// If the request is for a sign in, we redirect the user to the tenant address
// If the request is for a sign up, we exchange the OAuth code and get the user profile
func OAuthCallback() web.HandlerFunc {
	return func(c web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		provider := c.Param("provider")
		state := c.QueryParam("state")
		parts := strings.Split(state, "|")

		redirectURL, err := url.ParseRequestURI(parts[0])
		if err != nil {
			return c.Failure(err)
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		//Test OAuth
		if redirectURL.Path == fmt.Sprintf("/oauth/%s/echo", provider) {
			var query = redirectURL.Query()
			query.Set("code", code)
			query.Set("identifier", parts[1])
			redirectURL.RawQuery = query.Encode()
			return c.Redirect(redirectURL.String())
		}

		//Sign up process
		if redirectURL.Path == "/signup" {
			oauthUser, err := c.Services().OAuth.GetProfile(provider, code)
			if err != nil {
				return c.Failure(err)
			}

			claims := jwt.OAuthClaims{
				OAuthID:       oauthUser.ID,
				OAuthProvider: provider,
				OAuthName:     oauthUser.Name,
				OAuthEmail:    oauthUser.Email,
				Metadata: jwt.Metadata{
					ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
				},
			}

			token, err := jwt.Encode(claims)
			if err != nil {
				return c.Failure(err)
			}

			var query = redirectURL.Query()
			query.Set("token", token)
			redirectURL.RawQuery = query.Encode()
			return c.Redirect(redirectURL.String())
		}

		//Sign in process
		var query = redirectURL.Query()
		query.Set("code", code)
		query.Set("redirect", redirectURL.RequestURI())
		query.Set("identifier", parts[1])
		redirectURL.RawQuery = query.Encode()
		redirectURL.Path = fmt.Sprintf("/oauth/%s/token", provider)
		return c.Redirect(redirectURL.String())
	}
}

// SignInByOAuth is responsible for redirecting the user to the OAuth authorization URL for given provider
// A cookie is stored in user's browser with a random identifier that is later used to verify the authenticity of the request
func SignInByOAuth() web.HandlerFunc {
	return func(c web.Context) error {
		c.Response.Header().Add("X-Robots-Tag", "noindex")

		provider := c.Param("provider")
		redirect := c.QueryParam("redirect")

		if redirect == "" {
			redirect = c.BaseURL()
		}

		redirectURL, _ := url.ParseRequestURI(redirect)
		redirectURL.ResolveReference(c.Request.URL)

		if c.IsAuthenticated() && redirectURL.Path != fmt.Sprintf("/oauth/%s/echo", provider) {
			return c.Redirect(redirect)
		}

		authURL, err := c.Services().OAuth.GetAuthURL(provider, redirect, c.SessionID())
		if err != nil {
			return c.Failure(err)
		}
		return c.Redirect(authURL)
	}
}
