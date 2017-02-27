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
	"golang.org/x/oauth2/google"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/context"
	"github.com/labstack/echo"
)

type facebookUser struct {
	Name  string
	ID    string
	Email string
}

type googleUser struct {
	Name  string
	ID    string
	Email string
}

var (
	fbClientID         = os.Getenv("OAUTH_FACEBOOK_APPID")
	fbClientSecret     = os.Getenv("OAUTH_FACEBOOK_SECRET")
	googleClientID     = os.Getenv("OAUTH_GOOGLE_CLIENTID")
	googleClientSecret = os.Getenv("OAUTH_GOOGLE_SECRET")
	authEndpoint       = os.Getenv("AUTH_ENDPOINT")
)

func newFacebookOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     fbClientID,
		ClientSecret: fbClientSecret,
		RedirectURL:  authEndpoint + "/oauth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

func newGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  authEndpoint + "/oauth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
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

func googleMe(token *oauth2.Token) string {
	return "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken)
}

// FacebookCallback handles OAuth Facebook callbacks
func (h OAuthHandlers) FacebookCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirect := c.QueryParam("state")
		code := c.QueryParam("code")

		config := newFacebookOAuthConfig()
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
				user = &auth.User{
					Name:  fbUser.Name,
					Email: fbUser.Email,
					Providers: []*auth.UserProvider{
						&auth.UserProvider{
							UID:  fbUser.ID,
							Name: auth.OAuthFacebookProvider,
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

// FacebookLogin handles OAuth logins for Facebook
func (h OAuthHandlers) FacebookLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		config := newFacebookOAuthConfig()
		authURL, _ := url.Parse(config.Endpoint.AuthURL)
		parameters := url.Values{}
		parameters.Add("client_id", config.ClientID)
		parameters.Add("scope", strings.Join(config.Scopes, " "))
		parameters.Add("redirect_uri", config.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", c.QueryParam("redirect"))
		authURL.RawQuery = parameters.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, authURL.String())
	}
}

// GoogleCallback handles OAuth Google callbacks
func (h OAuthHandlers) GoogleCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirect := c.QueryParam("state")
		code := c.QueryParam("code")

		config := newGoogleOAuthConfig()
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.Logger().Errorf("googleOAuthConfig.Exchange() failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		gUser := &googleUser{}
		if err = doGet(googleMe(token), gUser); err != nil {
			c.Logger().Errorf("HTTP Google/Me failed with %s", err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		user, err := h.ctx.Auth.GetByEmail(gUser.Email)
		if err != nil {
			if err == auth.ErrUserNotFound {
				user = &auth.User{
					Name:  gUser.Name,
					Email: gUser.Email,
					Providers: []*auth.UserProvider{
						&auth.UserProvider{
							UID:  gUser.ID,
							Name: auth.OAuthGoogleProvider,
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

// GoogleLogin handles OAuth logins for Google
func (h OAuthHandlers) GoogleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		config := newGoogleOAuthConfig()
		authURL, _ := url.Parse(config.Endpoint.AuthURL)
		parameters := url.Values{}
		parameters.Add("client_id", config.ClientID)
		parameters.Add("scope", strings.Join(config.Scopes, " "))
		parameters.Add("redirect_uri", config.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", c.QueryParam("redirect"))
		authURL.RawQuery = parameters.Encode()
		return c.Redirect(http.StatusTemporaryRedirect, authURL.String())
	}
}
