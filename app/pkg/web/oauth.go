package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage"

	"github.com/getfider/fider/app/pkg/jsonq"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	systemProviders = []*models.OAuthConfig{
		&models.OAuthConfig{
			Provider:          oauth.FacebookProvider,
			DisplayName:       "Facebook",
			ProfileURL:        "https://graph.facebook.com/me?fields=name,email",
			ClientID:          os.Getenv("OAUTH_FACEBOOK_APPID"),
			ClientSecret:      os.Getenv("OAUTH_FACEBOOK_SECRET"),
			Scope:             "public_profile email",
			AuthorizeURL:      facebook.Endpoint.AuthURL,
			TokenURL:          facebook.Endpoint.TokenURL,
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "name",
			JSONUserEmailPath: "email",
		},
		&models.OAuthConfig{
			Provider:          oauth.GoogleProvider,
			DisplayName:       "Google",
			ProfileURL:        "https://www.googleapis.com/plus/v1/people/me",
			ClientID:          os.Getenv("OAUTH_GOOGLE_CLIENTID"),
			ClientSecret:      os.Getenv("OAUTH_GOOGLE_SECRET"),
			Scope:             "profile email",
			AuthorizeURL:      google.Endpoint.AuthURL,
			TokenURL:          google.Endpoint.TokenURL,
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "displayName",
			JSONUserEmailPath: "emails[0].value",
		},
		&models.OAuthConfig{
			Provider:          oauth.GitHubProvider,
			DisplayName:       "GitHub",
			ProfileURL:        "https://api.github.com/user",
			ClientID:          os.Getenv("OAUTH_GITHUB_CLIENTID"),
			ClientSecret:      os.Getenv("OAUTH_GITHUB_SECRET"),
			Scope:             "user:email",
			AuthorizeURL:      github.Endpoint.AuthURL,
			TokenURL:          github.Endpoint.TokenURL,
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "name, login",
			JSONUserEmailPath: "email",
		},
	}
)

//OAuthService implements real OAuth operations using Golang's oauth2 package
type OAuthService struct {
	oauthBaseURL  string
	tenantStorage storage.Tenant
}

//NewOAuthService creates a new OAuthService
func NewOAuthService(oauthBaseURL string, tenantStorage storage.Tenant) *OAuthService {
	return &OAuthService{
		oauthBaseURL,
		tenantStorage,
	}
}

//GetAuthURL returns authentication url for given provider
func (s *OAuthService) GetAuthURL(provider, redirect, identifier string) (string, error) {
	config, err := s.getConfig(provider)
	if err != nil {
		return "", err
	}

	authURL, _ := url.Parse(config.AuthorizeURL)
	parameters := url.Values{}
	parameters.Add("client_id", config.ClientID)
	parameters.Add("scope", config.Scope)
	parameters.Add("redirect_uri", fmt.Sprintf("%s/oauth/%s/callback", s.oauthBaseURL, provider))
	parameters.Add("response_type", "code")
	parameters.Add("state", redirect+"|"+identifier)
	authURL.RawQuery = parameters.Encode()
	return authURL.String(), nil
}

//GetProfile returns user profile based on provider and code
func (s *OAuthService) GetProfile(provider string, code string) (*oauth.UserProfile, error) {
	config, err := s.getConfig(provider)
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
		RedirectURL: fmt.Sprintf("%s/oauth/%s/callback", s.oauthBaseURL, provider),
	}).Exchange

	oauthToken, err := exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange OAuth2 code with %s", provider)
	}

	bytes, err := s.doGet(config.ProfileURL, oauthToken.AccessToken)
	if err != nil {
		return nil, err
	}

	return s.ParseProfileResponse(string(bytes), config)
}

//ParseProfileResponse parses profile response into UserProfile model
func (s *OAuthService) ParseProfileResponse(body string, config *models.OAuthConfig) (*oauth.UserProfile, error) {
	query := jsonq.New(body)
	profile := &oauth.UserProfile{
		ID:    strings.TrimSpace(query.String(config.JSONUserIDPath)),
		Name:  strings.TrimSpace(query.String(config.JSONUserNamePath)),
		Email: strings.TrimSpace(query.String(config.JSONUserEmailPath)),
	}

	if profile.ID == "" {
		return nil, oauth.ErrUserIDRequired
	}

	if profile.Name == "" && profile.Email != "" {
		parts := strings.Split(profile.Email, "@")
		profile.Name = parts[0]
	}

	if profile.Name == "" {
		profile.Name = "Anonymous"
	}

	if !validate.Email(profile.Email).Ok {
		profile.Email = ""
	}

	return profile, nil
}

//ListProviders returns a list of all available providers for current tenant
func (s *OAuthService) ListProviders() ([]*oauth.ProviderOption, error) {
	providers, err := s.listProvidersWithConfiguration()
	if err != nil {
		return nil, err
	}

	list := make([]*oauth.ProviderOption, 0)

	for _, p := range providers {
		if p.ClientID != "" {
			list = append(list, &oauth.ProviderOption{
				Provider:    p.Provider,
				DisplayName: p.DisplayName,
				URL:         fmt.Sprintf("/oauth/%s", p.Provider),
			})
		}
	}

	return list, nil
}

func (s *OAuthService) listProvidersWithConfiguration() ([]*models.OAuthConfig, error) {
	list := make([]*models.OAuthConfig, 0)

	for _, p := range systemProviders {
		if p.ClientID != "" {
			list = append(list, p)
		}
	}

	customProviders, err := s.tenantStorage.ListOAuthConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of custom OAuth providers")
	}

	for _, p := range customProviders {
		if p.ClientID != "" {
			list = append(list, p)
		}
	}

	return list, nil
}

func (s *OAuthService) doGet(url, accessToken string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request GET %s", url)
	}

	if r.StatusCode != 200 {
		return nil, errors.New("failed to request GET %s with status code %d and response %s", url, r.StatusCode, string(bytes))
	}

	return bytes, nil
}

func (s *OAuthService) getConfig(provider string) (*models.OAuthConfig, error) {
	for _, config := range systemProviders {
		if config.ClientID != "" && config.Provider == provider {
			return config, nil
		}
	}

	return s.tenantStorage.GetOAuthConfigByProvider(provider)
}
