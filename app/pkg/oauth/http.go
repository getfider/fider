package oauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getfider/fider/app/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type providerSettings struct {
	profileURL func(token *oauth2.Token) string
	config     func(authEndpoint string) *oauth2.Config
}

var (
	oauthSettings = map[string]*providerSettings{
		FacebookProvider: {
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
		GoogleProvider: {
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
		GitHubProvider: {
			profileURL: func(token *oauth2.Token) string {
				return "https://api.github.com/user?access_token=" + url.QueryEscape(token.AccessToken)
			},
			config: func(authEndpoint string) *oauth2.Config {
				return &oauth2.Config{
					ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
					ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
					RedirectURL:  authEndpoint + "/oauth/github/callback",
					Scopes: []string{
						"user:email",
					},
					Endpoint: github.Endpoint,
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
		return errors.Wrap(err, "failed to request GET %s", url)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return errors.Wrap(err, "failed unmarshal response")
	}

	return nil
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
		return nil, errors.Wrap(err, "failed to exchange OAuth2 code with %s", provider)
	}

	profile := &UserProfile{}
	if err = doGet(settings.profileURL(oauthToken), profile); err != nil {
		return nil, err
	}

	//GitHub allows users to omit name, so we use their login name
	if strings.Trim(profile.Name, " ") == "" {
		profile.Name = profile.Login
	}

	return profile, nil
}
