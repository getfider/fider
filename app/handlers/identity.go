package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage"
)

type oauthUserProfile struct {
	Name  string
	ID    string
	Email string
}

// OAuthHandlers contains multiple oauth HTTP handlers
type OAuthHandlers struct {
	tenants storage.Tenant
	oauth   oauth.Service
	users   storage.User
}

// OAuth creates a new OAuthHandlers
func OAuth(tenants storage.Tenant, oauth oauth.Service, users storage.User) OAuthHandlers {
	return OAuthHandlers{tenants, oauth, users}
}

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback(provider string) web.HandlerFunc {
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

		tenant := c.Tenant()
		if tenant == nil {
			// should get from context
			// Single/Multi middleware should handle auth endpoint properly
			// .Set("AuthEndpoint") is not good as well
			tenant, err = h.tenants.GetByDomain(stripPort(redirectURL.Host))
			if err != nil {
				return c.Failure(err)
			}
		}

		oauthUser, err := h.oauth.GetProfile(c.AuthEndpoint(), provider, code)
		if err != nil {
			return c.Failure(err)
		}

		user, err := h.users.GetByEmail(tenant.ID, oauthUser.Email)
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

				err = h.users.Register(user)
				if err != nil {
					return c.Failure(err)
				}
			} else {
				return c.Failure(err)
			}
		} else if !user.HasProvider(provider) {
			err = h.users.RegisterProvider(user.ID, &models.UserProvider{
				UID:  oauthUser.ID,
				Name: provider,
			})
			if err != nil {
				return c.Failure(err)
			}
		}

		claims := &models.WechyClaims{
			UserID:    user.ID,
			UserName:  user.Name,
			UserEmail: user.Email,
		}

		var token string
		if token, err = jwt.Encode(claims); err != nil {
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
func (h OAuthHandlers) Login(provider string) web.HandlerFunc {
	return func(c web.Context) error {
		authURL := h.oauth.GetAuthURL(c.AuthEndpoint(), provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
