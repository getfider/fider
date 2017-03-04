package handler

import (
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/labstack/echo"
)

type oauthUserProfile struct {
	Name  string
	ID    string
	Email string
}

type oauthProviderSettings struct {
	profileURL func(*oauth2.Token) string
	config     *oauth2.Config
}

// OAuthHandlers contains multiple oauth HTTP handlers
type OAuthHandlers struct {
	ctx *context.WchyContext
}

// OAuth creates a new OAuthHandlers
func OAuth(ctx *context.WchyContext) OAuthHandlers {
	return OAuthHandlers{ctx: ctx}
}

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback(provider string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var redirectURL *url.URL
		var err error

		redirect := c.QueryParam("state")
		code := c.QueryParam("code")

		if redirectURL, err = url.Parse(redirect); err != nil {
			c.Logger().Errorf("Could not parse url %s", redirect)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		//TODO: Check if code is empty (or other querystring parameter)
		//Because the user can deny access to it

		oauthUser, err := h.ctx.OAuth.GetProfile(provider, code)
		if err != nil {
			c.Logger().Error(err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		user, err := h.ctx.User.GetByEmail(oauthUser.Email)
		if err != nil {
			if err == service.ErrNotFound {
				user = &model.User{
					Name:  oauthUser.Name,
					Email: oauthUser.Email,
					Providers: []*model.UserProvider{
						&model.UserProvider{
							UID:  oauthUser.ID,
							Name: provider,
						},
					},
				}

				err = h.ctx.User.Register(user)
				if err != nil {
					c.Logger().Error(err)
				}
			} else {
				c.Logger().Error(err)
			}
		}

		claims := &auth.WchyClaims{
			UserID:    user.ID,
			UserName:  user.Name,
			UserEmail: user.Email,
		}

		var token string
		if token, err = auth.Encode(claims); err != nil {
			c.Logger().Errorf("Encoding claims failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		var query = redirectURL.Query()
		query.Add("jwt", token)
		redirectURL.RawQuery = query.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
	}
}

// Login handles OAuth logins
func (h OAuthHandlers) Login(provider string) echo.HandlerFunc {
	return func(c echo.Context) error {
		authURL := h.ctx.OAuth.GetAuthURL(provider, c.QueryParam("redirect"))
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
