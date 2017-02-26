package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

type facebookUser struct {
	Name  string
	ID    string
	Email string
}

var (
	oauthState     = os.Getenv("OAUTH_STATE")
	fbClientID     = os.Getenv("OAUTH_FACEBOOK_APPID")
	fbClientSecret = os.Getenv("OAUTH_FACEBOOK_SECRET")
	authEndpoint   = os.Getenv("AUTH_ENDPOINT")
)

func newFacebookOAuthConfig(redirect string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauthState,
		ClientSecret: fbClientSecret,
		RedirectURL:  authEndpoint + "/oauth/facebook/callback?redirect=" + url.QueryEscape(redirect),
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

// OAuthHandlers contains multiple oauth HTTP handlers
type OAuthHandlers struct {
	ctx *context.WchyContext
}

// OAuth creates a new OAuthHandlers
func OAuth(ctx *context.WchyContext) OAuthHandlers {
	return OAuthHandlers{ctx: ctx}
}

func doGet(url string, v interface{}) error {
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func facebookMe(token *oauth2.Token) string {
	return "https://graph.facebook.com/me?fields=name,email&access_token=" + url.QueryEscape(token.AccessToken)
}

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		state := c.QueryParam("state")
		code := c.QueryParam("code")
		redirect := c.QueryParam("redirect")

		if state != oauthState {
			c.Logger().Errorf("Invalid OAuth state '%s'", state)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		config := newFacebookOAuthConfig(redirect)
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.Logger().Errorf("facebookOAuthConfig.Exchange() failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		fbUser := &facebookUser{}
		if err = doGet(facebookMe(token), fbUser); err != nil {
			c.Logger().Errorf("HTTP Facebook/Me failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		user, err := h.ctx.Auth.GetByEmail(fbUser.Email)
		if err != nil {
			if err == auth.ErrUserNotFound {
				user := &auth.User{
					Name:  fbUser.Name,
					Email: fbUser.Email,
					Providers: []*auth.UserProvider{
						&auth.UserProvider{
							UID:  fbUser.ID,
							Name: auth.OAUTH_FACEBOOK_PROVIDER,
						},
					},
				}
				err = h.ctx.Auth.Register(user)
				if err != nil {
					c.Logger().Error(err)
				}
			} else {
				c.Logger().Error(err)
			}
		}

		c.Logger().Infof("Logged in as %s", user)

		return c.Redirect(http.StatusTemporaryRedirect, redirect)
	}
}

// Login handlers OAuth logins
func (h OAuthHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		config := newFacebookOAuthConfig(c.QueryParam("redirect"))
		authURL, _ := url.Parse(config.Endpoint.AuthURL)
		parameters := url.Values{}
		parameters.Add("client_id", config.ClientID)
		parameters.Add("scope", strings.Join(config.Scopes, " "))
		parameters.Add("redirect_uri", config.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", oauthState)
		authURL.RawQuery = parameters.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, authURL.String())
	}
}
