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
	tenantStorage storage.Tenant
	oauthService  oauth.Service
	userStorage   storage.User
}

// OAuth creates a new OAuthHandlers
func OAuth(tenantStorage storage.Tenant, oauthService oauth.Service, userStorage storage.User) OAuthHandlers {
	return OAuthHandlers{tenantStorage, oauthService, userStorage}
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

		tenant, err := h.tenantStorage.GetByDomain(stripPort(redirectURL.Host))
		if err != nil {
			return c.Failure(err)
		}

		oauthUser, err := h.oauthService.GetProfile(provider, code)
		if err != nil {
			return c.Failure(err)
		}

		user, err := h.userStorage.GetByEmail(tenant.ID, oauthUser.Email)
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

				err = h.userStorage.Register(user)
				if err != nil {
					return c.Failure(err)
				}
			} else {
				return c.Failure(err)
			}
		} else if !user.HasProvider(provider) {
			err = h.userStorage.RegisterProvider(user.ID, &models.UserProvider{
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
		authURL := h.oauthService.GetAuthURL(provider, c.QueryParam("redirect"))
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
