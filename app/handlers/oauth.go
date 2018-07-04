package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/web/util"
)

type oauthUserProfile struct {
	Name  string
	ID    string
	Email string
}

// OAuthToken exchanges Authorization Code for Authentication Token
func OAuthToken() web.HandlerFunc {
	return func(c web.Context) error {
		provider := c.Param("provider")

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(c.BaseURL())
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
					return c.Redirect(c.BaseURL() + "/not-invited")
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

		redirectURL, _ := url.Parse(c.Request.URL.String())
		var query = redirectURL.Query()
		query.Del("code")
		query.Del("path")
		redirectURL.RawQuery = query.Encode()
		redirectURL.Path = c.QueryParam("path")
		return c.Redirect(redirectURL.String())
	}
}

// OAuthCallback handles OAuth callbacks
func OAuthCallback() web.HandlerFunc {
	return func(c web.Context) error {
		provider := c.Param("provider")
		redirect := c.QueryParam("state")
		redirectURL, err := url.ParseRequestURI(redirect)
		if err != nil {
			return c.Failure(err)
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirect)
		}

		//Sign in process
		if redirectURL.Path != "/signup" {
			var query = redirectURL.Query()
			query.Set("code", code)
			query.Set("path", redirectURL.Path)
			redirectURL.RawQuery = query.Encode()
			redirectURL.Path = fmt.Sprintf("/oauth/%s/token", provider)
			return c.Redirect(redirectURL.String())
		}

		//Sign up process
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
}

// SignInByOAuth handles OAuth sign in
func SignInByOAuth() web.HandlerFunc {
	return func(c web.Context) error {
		provider := c.Param("provider")
		authURL, err := c.Services().OAuth.GetAuthURL(provider, c.QueryParam("redirect"))
		if err != nil {
			return c.Failure(err)
		}
		return c.Redirect(authURL)
	}
}
