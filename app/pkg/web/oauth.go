package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
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

func getProviderStatus(envKey string) int {
	if os.Getenv(envKey) == "" {
		return models.OAuthConfigDisabled
	}
	return models.OAuthConfigEnabled
}

var (
	systemProviders = []*models.OAuthConfig{
		&models.OAuthConfig{
			Provider:          oauth.FacebookProvider,
			DisplayName:       "Facebook",
			ProfileURL:        "https://graph.facebook.com/me?fields=name,email",
			Status:            getProviderStatus("OAUTH_FACEBOOK_APPID"),
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
			Status:            getProviderStatus("OAUTH_GOOGLE_CLIENTID"),
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
			Status:            getProviderStatus("OAUTH_GITHUB_CLIENTID"),
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
	configCache   map[string]*models.OAuthConfig
}

//NewOAuthService creates a new OAuthService
func NewOAuthService(oauthBaseURL string, tenantStorage storage.Tenant) *OAuthService {
	configCache := make(map[string]*models.OAuthConfig, 0)

	return &OAuthService{
		oauthBaseURL,
		tenantStorage,
		configCache,
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

//GetProfile returns user profile based on provider (only if enabled) and code
func (s *OAuthService) GetProfile(provider string, code string) (*oauth.UserProfile, error) {
	config, err := s.getConfig(provider)
	if err != nil {
		return nil, err
	}

	if config.Status == models.OAuthConfigDisabled {
		return nil, errors.New("Provider %s is disabled", provider)
	}

	body, err := s.GetRawProfile(provider, code)
	if err != nil {
		return nil, err
	}

	return s.ParseRawProfile(provider, body)
}

//GetRawProfile returns raw JSON response from Profile API
func (s *OAuthService) GetRawProfile(provider string, code string) (string, error) {
	config, err := s.getConfig(provider)
	if err != nil {
		return "", err
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
		return "", err
	}

	if config.ProfileURL == "" {
		parts := strings.Split(oauthToken.AccessToken, ".")
		if len(parts) != 3 {
			return "", errors.New("AccessToken is not JWT")
		}

		body, _ := jwt.DecodeSegment(parts[1])
		return string(body), nil
	}

	statusCode, body, err := s.doGet(config.ProfileURL, oauthToken.AccessToken)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", errors.New("Failed to request profile. Status Code: %d. Body: %s", statusCode, body)
	}

	return body, nil
}

//ParseRawProfile parses raw profile response into UserProfile model
func (s *OAuthService) ParseRawProfile(provider, body string) (*oauth.UserProfile, error) {
	config, err := s.getConfig(provider)
	if err != nil {
		return nil, err
	}

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

	if len(validate.Email(profile.Email)) != 0 {
		profile.Email = ""
	}

	return profile, nil
}

//ListActiveProviders returns a list of all enabled providers for current tenant
func (s *OAuthService) ListActiveProviders() ([]*oauth.ProviderOption, error) {
	list := make([]*oauth.ProviderOption, 0)
	providers, err := s.ListAllProviders()
	if err != nil {
		return nil, err
	}

	for _, p := range providers {
		if p.IsEnabled {
			list = append(list, p)
		}
	}
	return list, nil
}

//ListAllProviders returns a list of all providers for current tenant
func (s *OAuthService) ListAllProviders() ([]*oauth.ProviderOption, error) {
	providers, err := s.allOAuthConfigs()
	if err != nil {
		return nil, err
	}

	list := make([]*oauth.ProviderOption, 0)

	for _, p := range providers {
		list = append(list, &oauth.ProviderOption{
			Provider:         p.Provider,
			DisplayName:      p.DisplayName,
			ClientID:         p.ClientID,
			URL:              fmt.Sprintf("/oauth/%s", p.Provider),
			CallbackURL:      fmt.Sprintf("%s/oauth/%s/callback", s.oauthBaseURL, p.Provider),
			IsCustomProvider: string(p.Provider[0]) == "_",
			LogoID:           p.LogoID,
			IsEnabled:        p.Status == models.OAuthConfigEnabled,
		})
	}

	return list, nil
}

func (s *OAuthService) allOAuthConfigs() ([]*models.OAuthConfig, error) {
	list := make([]*models.OAuthConfig, 0)

	for _, p := range systemProviders {
		list = append(list, p)
	}

	customProviders, err := s.tenantStorage.ListOAuthConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of custom OAuth providers")
	}

	for _, p := range customProviders {
		list = append(list, p)
	}

	return list, nil
}

func (s *OAuthService) doGet(url, accessToken string) (int, string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, "", err
	}

	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	return r.StatusCode, string(bytes), err
}

func (s *OAuthService) getConfig(provider string) (*models.OAuthConfig, error) {
	for _, config := range systemProviders {
		if config.Status == models.OAuthConfigEnabled && config.Provider == provider {
			return config, nil
		}
	}

	var err error
	config, ok := s.configCache[provider]
	if !ok {
		config, err = s.tenantStorage.GetOAuthConfigByProvider(provider)
		if err != nil {
			return nil, err
		}

		s.configCache[provider] = config
	}

	return config, nil
}
