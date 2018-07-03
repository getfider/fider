package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type providerSettings struct {
	profileURL string
	config     func(authEndpoint string) *oauth2.Config
}

var (
	systemProviders = map[string]*models.OAuthConfig{
		FacebookProvider: &models.OAuthConfig{
			ProfileURL:   "https://graph.facebook.com/me?fields=name,email",
			ClientID:     os.Getenv("OAUTH_FACEBOOK_APPID"),
			ClientSecret: os.Getenv("OAUTH_FACEBOOK_SECRET"),
			Scope:        "public_profile email",
			AuthorizeURL: facebook.Endpoint.AuthURL,
			TokenURL:     facebook.Endpoint.TokenURL,
		},
		GoogleProvider: &models.OAuthConfig{
			ProfileURL:   "https://www.googleapis.com/oauth2/v2/userinfo",
			ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENTID"),
			ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
			Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
			AuthorizeURL: google.Endpoint.AuthURL,
			TokenURL:     google.Endpoint.TokenURL,
		},
		GitHubProvider: &models.OAuthConfig{
			ProfileURL:   "https://api.github.com/user",
			ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
			ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
			Scope:        "user:email",
			AuthorizeURL: github.Endpoint.AuthURL,
			TokenURL:     github.Endpoint.TokenURL,
		},
	}
)

func doGet(url, accessToken string, v interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer r.Body.Close()
	if r.StatusCode != 200 {
		return errors.New("failed to request GET %s with status code %d", url, r.StatusCode)
	}

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
type HTTPService struct {
	authEndpoint string
}

//NewHTTPService creates a new HTTPService
func NewHTTPService(authEndpoint string) *HTTPService {
	return &HTTPService{
		authEndpoint,
	}
}

//GetAuthURL returns authentication url for given provider
func (p *HTTPService) GetAuthURL(provider string, redirect string) (string, error) {
	config, err := p.getConfig(provider)
	if err != nil {
		return "", err
	}

	authURL, _ := url.Parse(config.AuthorizeURL)
	parameters := url.Values{}
	parameters.Add("client_id", config.ClientID)
	parameters.Add("scope", config.Scope)
	parameters.Add("redirect_uri", fmt.Sprintf("%s/oauth/%s/callback", p.authEndpoint, provider))
	parameters.Add("response_type", "code")
	parameters.Add("state", redirect)
	authURL.RawQuery = parameters.Encode()
	return authURL.String(), nil
}

//GetProfile returns user profile based on provider and code
func (p *HTTPService) GetProfile(provider string, code string) (*UserProfile, error) {
	config, err := p.getConfig(provider)
	if err != nil {
		return nil, err
	}

	exchange := (&oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthorizeURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: fmt.Sprintf("%s/oauth/%s/callback", p.authEndpoint, provider),
	}).Exchange

	oauthToken, err := exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange OAuth2 code with %s", provider)
	}

	profile := &UserProfile{}
	if err = doGet(config.ProfileURL, oauthToken.AccessToken, profile); err != nil {
		return nil, err
	}

	//GitHub allows users to omit name, so we use their login name
	if strings.Trim(profile.Name, " ") == "" {
		profile.Name = profile.Login
	}

	return profile, nil
}

func (p *HTTPService) getConfig(provider string) (*models.OAuthConfig, error) {
	return systemProviders[provider], nil
}
