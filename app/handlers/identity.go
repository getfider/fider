package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
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

			tenant := c.Tenant()
			if tenant == nil {
				// should get from context
				// Single/Multi middleware should handle auth endpoint properly
				// .Set("AuthEndpoint") is not good as well
				tenant, err = c.Services().Tenants.GetByDomain(stripPort(redirectURL.Host))
				if err != nil {
					return c.Failure(err)
				}
			}

			users := c.Services().Users
			user, err := users.GetByEmail(tenant.ID, oauthUser.Email)
			if err != nil {
				if err == app.ErrNotFound {
					user = &models.User{
						Name:   oauthUser.Name,
						Tenant: tenant,
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

			claims = models.FiderClaims{
				UserID:    user.ID,
				UserName:  user.Name,
				UserEmail: user.Email,
			}
		} else {
			claims = models.OAuthClaims{
				ID:       oauthUser.ID,
				Provider: provider,
				Name:     oauthUser.Name,
				Email:    oauthUser.Email,
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

// Logout remove auth cookies
func Logout() web.HandlerFunc {
	return func(c web.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    "auth",
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
		})
		return c.Redirect(http.StatusTemporaryRedirect, c.QueryParam("redirect"))
	}
}

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c web.Context) error {
		subdomain := c.Param("subdomain")
		available, err := c.Services().Tenants.IsSubdomainAvailable(subdomain)

		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(echo.Map{
			"isAvailable": available,
			"subdomain":   subdomain,
		})
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
