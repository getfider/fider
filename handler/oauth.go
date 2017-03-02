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
	jwt "github.com/dgrijalva/jwt-go"
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

var (
	authEndpoint  = os.Getenv("AUTH_ENDPOINT")
	jwtSecret     = os.Getenv("JWT_SECRET")
	oauthSettings = map[string]*oauthProviderSettings{
		auth.OAuthFacebookProvider: &oauthProviderSettings{
			profileURL: func(token *oauth2.Token) string {
				return "https://graph.facebook.com/me?fields=name,email&access_token=" + url.QueryEscape(token.AccessToken)
			},
			config: &oauth2.Config{
				ClientID:     os.Getenv("OAUTH_FACEBOOK_APPID"),
				ClientSecret: os.Getenv("OAUTH_FACEBOOK_SECRET"),
				RedirectURL:  authEndpoint + "/oauth/facebook/callback",
				Scopes:       []string{"public_profile", "email"},
				Endpoint:     facebook.Endpoint,
			},
		},
		auth.OAuthGoogleProvider: &oauthProviderSettings{
			profileURL: func(token *oauth2.Token) string {
				return "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken)
			},
			config: &oauth2.Config{
				ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENTID"),
				ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
				RedirectURL:  authEndpoint + "/oauth/google/callback",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email",
				},
				Endpoint: google.Endpoint,
			},
		},
	}
)

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

// Callback handles OAuth callbacks
func (h OAuthHandlers) Callback(provider string) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirect := c.QueryParam("state")
		code := c.QueryParam("code")

		//TODO: Check if code is empty (or other querystring parameter)
		//Because the user can deny access to it

		settings := oauthSettings[provider]
		oauthToken, err := settings.config.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.Logger().Errorf("%s oauthConfig.Exchange() failed with %s", provider, err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		oauthUser := &oauthUserProfile{}
		if err = doGet(settings.profileURL(oauthToken), oauthUser); err != nil {
			c.Logger().Errorf("HTTP Get profile for %s failed with %s", provider, err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		user, err := h.ctx.Auth.GetByEmail(oauthUser.Email)
		if err != nil {
			if err == auth.ErrUserNotFound {
				user = &auth.User{
					Name:  oauthUser.Name,
					Email: oauthUser.Email,
					Providers: []*auth.UserProvider{
						&auth.UserProvider{
							UID:  oauthUser.ID,
							Name: provider,
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

		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
			"UserID":    user.ID,
			"UserName":  user.Name,
			"UserEmail": user.Email,
		})

		var token string
		if token, err = jwtToken.SignedString([]byte(jwtSecret)); err != nil {
			c.Logger().Errorf("%s oauthConfig.Exchange() failed with %s", provider, err)
			return c.Redirect(http.StatusTemporaryRedirect, redirect) //TODO: redirect to some error page
		}

		var redirectURL *url.URL
		if redirectURL, err = url.Parse(redirect); err != nil {
			c.Logger().Errorf("Could not parse url %s", redirect)
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
		config := oauthSettings[provider].config
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
