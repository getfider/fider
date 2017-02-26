package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

var (
	oauthState = os.Getenv("OAUTH_STATE")
)

func newFacebookOAuthConfig(redirect string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_FACEBOOK_APPID"),
		ClientSecret: os.Getenv("OAUTH_FACEBOOK_SECRET"),
		RedirectURL:  os.Getenv("AUTH_ENDPOINT") + "/oauth/callback?redirect=" + url.QueryEscape(redirect),
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

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		state := c.QueryParam("state")
		code := c.QueryParam("code")
		redirect := c.QueryParam("redirect")

		if state != oauthState {
			c.Logger().Infof("Invalid OAuth state '%s'", state)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		config := newFacebookOAuthConfig(redirect)
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.Logger().Errorf("facebookOAuthConfig.Exchange() failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		resp, err := http.Get("https://graph.facebook.com/me?fields=name,email&access_token=" + url.QueryEscape(token.AccessToken))
		defer resp.Body.Close()
		response, err := ioutil.ReadAll(resp.Body)
		log.Printf("parseResponseBody: %s\n", string(response))

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
