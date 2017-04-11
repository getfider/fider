package identity

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/infra"
)

type oauthUserProfile struct {
	Name  string
	ID    string
	Email string
}

// OAuthHandlers contains multiple oauth HTTP handlers
type OAuthHandlers struct {
	tenantService TenantService
	oauthService  OAuthService
	userService   UserService
}

// OAuth creates a new OAuthHandlers
func OAuth(tenantService TenantService, oauthService OAuthService, userService UserService) OAuthHandlers {
	return OAuthHandlers{tenantService, oauthService, userService}
}

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback(provider string) app.HandlerFunc {
	return func(c app.Context) error {

		redirect := c.QueryParam("state")
		redirectURL, err := url.ParseRequestURI(redirect)
		if err != nil {
			return c.Failure(err)
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.Redirect(http.StatusTemporaryRedirect, redirect)
		}

		tenant, err := h.tenantService.GetByDomain(stripPort(redirectURL.Host))
		if err != nil {
			return c.Failure(err)
		}

		oauthUser, err := h.oauthService.GetProfile(provider, code)
		if err != nil {
			return c.Failure(err)
		}

		user, err := h.userService.GetByEmail(tenant.ID, oauthUser.Email)
		if err != nil {
			if err == app.ErrNotFound {
				user = &app.User{
					Name:   oauthUser.Name,
					Tenant: tenant,
					Email:  oauthUser.Email,
					Role:   app.RoleVisitor,
					Providers: []*app.UserProvider{
						&app.UserProvider{
							UID:  oauthUser.ID,
							Name: provider,
						},
					},
				}

				err = h.userService.Register(user)
				if err != nil {
					return c.Failure(err)
				}
			} else {
				return c.Failure(err)
			}
		} else if !user.HasProvider(provider) {
			err = h.userService.RegisterProvider(user.ID, &app.UserProvider{
				UID:  oauthUser.ID,
				Name: provider,
			})
			if err != nil {
				return c.Failure(err)
			}
		}

		claims := &app.WechyClaims{
			UserID:    user.ID,
			UserName:  user.Name,
			UserEmail: user.Email,
		}

		var token string
		if token, err = infra.Encode(claims); err != nil {
			c.Logger().Errorf("Encoding claims failed with %s", err)
			return c.Failure(err)
		}

		var query = redirectURL.Query()
		query.Add("jwt", token)
		redirectURL.RawQuery = query.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
	}
}

// Login handles OAuth logins
func (h OAuthHandlers) Login(provider string) app.HandlerFunc {
	return func(c app.Context) error {
		authURL := h.oauthService.GetAuthURL(provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// Logout remove auth cookies
func (h OAuthHandlers) Logout() app.HandlerFunc {
	return func(c app.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    "auth",
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
		})
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
