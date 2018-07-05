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
	"github.com/getfider/fider/app/pkg/uuid"
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
		redirectURL, _ := url.ParseRequestURI(c.QueryParam("redirect"))
		redirectURL.ResolveReference(c.Request.URL)

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(redirectURL.String())
		}

		identifier := c.QueryParam("identifier")
		cookie, err := c.Request.Cookie("__oauth_identifier")
		if err != nil {
			return c.Failure(errors.Wrap(err, "failed to get oauth identifier cookie"))
		}

		c.RemoveCookie("__oauth_identifier")
		if identifier != cookie.Value {
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

// OAuthCallback handles OAuth callbacks
func OAuthCallback() web.HandlerFunc {
	return func(c web.Context) error {
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

		//Sign in process
		if redirectURL.Path != "/signup" {
			var query = redirectURL.Query()
			query.Set("code", code)
			query.Set("redirect", redirectURL.RequestURI())
			query.Set("identifier", parts[1])
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
		identifier := uuid.NewV4().String()
		c.AddCookie("__oauth_identifier", identifier, time.Now().Add(5*time.Minute))
		authURL, err := c.Services().OAuth.GetAuthURL(provider, c.QueryParam("redirect"), identifier)
		if err != nil {
			return c.Failure(err)
		}
		return c.Redirect(authURL)
	}
}
