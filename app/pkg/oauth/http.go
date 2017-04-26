package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type providerSettings struct {
	profileURL func(token *oauth2.Token) string
	config     func(authEndpoint string) *oauth2.Config
}

var (
	oauthSettings = map[string]*providerSettings{
		FacebookProvider: &providerSettings{
			profileURL: func(token *oauth2.Token) string {
				return "https://graph.facebook.com/me?fields=name,email&access_token=" + url.QueryEscape(token.AccessToken)
			},
			config: func(authEndpoint string) *oauth2.Config {
				return &oauth2.Config{
					ClientID:     os.Getenv("OAUTH_FACEBOOK_APPID"),
					ClientSecret: os.Getenv("OAUTH_FACEBOOK_SECRET"),
					RedirectURL:  authEndpoint + "/oauth/facebook/callback",
					Scopes:       []string{"public_profile", "email"},
					Endpoint:     facebook.Endpoint,
				}
			},
		},
		GoogleProvider: &providerSettings{
			profileURL: func(token *oauth2.Token) string {
				return "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken)
			},
			config: func(authEndpoint string) *oauth2.Config {
				return &oauth2.Config{
					ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENTID"),
					ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
					RedirectURL:  authEndpoint + "/oauth/google/callback",
					Scopes: []string{
						"https://www.googleapis.com/auth/userinfo.profile",
						"https://www.googleapis.com/auth/userinfo.email",
					},
					Endpoint: google.Endpoint,
				}
			},
		},
	}
)

func doGet(url string, v interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

//HTTPService implements real OAuth operations using Golang's oauth2 package
type HTTPService struct{}

//GetAuthURL returns authentication url for given provider
func (p *HTTPService) GetAuthURL(authEndpoint string, provider string, redirect string) string {
	config := oauthSettings[provider].config(authEndpoint)

	authURL, _ := url.Parse(config.Endpoint.AuthURL)
	parameters := url.Values{}
	parameters.Add("client_id", config.ClientID)
	parameters.Add("scope", strings.Join(config.Scopes, " "))
	parameters.Add("redirect_uri", config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", redirect)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

//GetProfile returns user profile based on provider and code
func (p *HTTPService) GetProfile(authEndpoint string, provider string, code string) (*UserProfile, error) {
	settings := oauthSettings[provider]
	oauthToken, err := settings.config(authEndpoint).Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("OAuth code Exchange for %s failed with %s", provider, err)
	}

	profile := &UserProfile{}
	if err = doGet(settings.profileURL(oauthToken), profile); err != nil {
		return nil, fmt.Errorf("HTTP Get profile for %s failed with %s", provider, err)
	}

	return profile, nil
}
