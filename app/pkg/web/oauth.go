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
			LogoURL:           "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHg9IjBweCIgeT0iMHB4IgogICAgIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIKICAgICB2aWV3Qm94PSIwIDAgNDggNDgiCiAgICAgc3R5bGU9ImZpbGw6IzAwMDAwMDsiPjxnIGlkPSJzdXJmYWNlMSI+PHBhdGggc3R5bGU9IiBmaWxsOiMzRjUxQjU7IiBkPSJNIDQyIDM3IEMgNDIgMzkuNzYxNzE5IDM5Ljc2MTcxOSA0MiAzNyA0MiBMIDExIDQyIEMgOC4yMzgyODEgNDIgNiAzOS43NjE3MTkgNiAzNyBMIDYgMTEgQyA2IDguMjM4MjgxIDguMjM4MjgxIDYgMTEgNiBMIDM3IDYgQyAzOS43NjE3MTkgNiA0MiA4LjIzODI4MSA0MiAxMSBaICI+PC9wYXRoPjxwYXRoIHN0eWxlPSIgZmlsbDojRkZGRkZGOyIgZD0iTSAzNC4zNjcxODggMjUgTCAzMSAyNSBMIDMxIDM4IEwgMjYgMzggTCAyNiAyNSBMIDIzIDI1IEwgMjMgMjEgTCAyNiAyMSBMIDI2IDE4LjU4OTg0NCBDIDI2LjAwMzkwNiAxNS4wODIwMzEgMjcuNDYwOTM4IDEzIDMxLjU5Mzc1IDEzIEwgMzUgMTMgTCAzNSAxNyBMIDMyLjcxNDg0NCAxNyBDIDMxLjEwNTQ2OSAxNyAzMSAxNy42MDE1NjMgMzEgMTguNzIyNjU2IEwgMzEgMjEgTCAzNSAyMSBaICI+PC9wYXRoPjwvZz48L3N2Zz4=",
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
			LogoURL:           "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHg9IjBweCIgeT0iMHB4IgogICAgIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIKICAgICB2aWV3Qm94PSIwIDAgNDggNDgiCiAgICAgc3R5bGU9ImZpbGw6IzAwMDAwMDsiPjxnIGlkPSJzdXJmYWNlMSI+PHBhdGggc3R5bGU9IiBmaWxsOiNGRkMxMDc7IiBkPSJNIDQzLjYwOTM3NSAyMC4wODIwMzEgTCA0MiAyMC4wODIwMzEgTCA0MiAyMCBMIDI0IDIwIEwgMjQgMjggTCAzNS4zMDQ2ODggMjggQyAzMy42NTIzNDQgMzIuNjU2MjUgMjkuMjIyNjU2IDM2IDI0IDM2IEMgMTcuMzcxMDk0IDM2IDEyIDMwLjYyODkwNiAxMiAyNCBDIDEyIDE3LjM3MTA5NCAxNy4zNzEwOTQgMTIgMjQgMTIgQyAyNy4wNTg1OTQgMTIgMjkuODQzNzUgMTMuMTUyMzQ0IDMxLjk2MDkzOCAxNS4wMzkwNjMgTCAzNy42MTcxODggOS4zODI4MTMgQyAzNC4wNDY4NzUgNi4wNTQ2ODggMjkuMjY5NTMxIDQgMjQgNCBDIDEyLjk1MzEyNSA0IDQgMTIuOTUzMTI1IDQgMjQgQyA0IDM1LjA0Njg3NSAxMi45NTMxMjUgNDQgMjQgNDQgQyAzNS4wNDY4NzUgNDQgNDQgMzUuMDQ2ODc1IDQ0IDI0IEMgNDQgMjIuNjYwMTU2IDQzLjg2MzI4MSAyMS4zNTE1NjMgNDMuNjA5Mzc1IDIwLjA4MjAzMSBaICI+PC9wYXRoPjxwYXRoIHN0eWxlPSIgZmlsbDojRkYzRDAwOyIgZD0iTSA2LjMwNDY4OCAxNC42OTE0MDYgTCAxMi44Nzg5MDYgMTkuNTExNzE5IEMgMTQuNjU2MjUgMTUuMTA5Mzc1IDE4Ljk2MDkzOCAxMiAyNCAxMiBDIDI3LjA1ODU5NCAxMiAyOS44NDM3NSAxMy4xNTIzNDQgMzEuOTYwOTM4IDE1LjAzOTA2MyBMIDM3LjYxNzE4OCA5LjM4MjgxMyBDIDM0LjA0Njg3NSA2LjA1NDY4OCAyOS4yNjk1MzEgNCAyNCA0IEMgMTYuMzE2NDA2IDQgOS42NTYyNSA4LjMzNTkzOCA2LjMwNDY4OCAxNC42OTE0MDYgWiAiPjwvcGF0aD48cGF0aCBzdHlsZT0iIGZpbGw6IzRDQUY1MDsiIGQ9Ik0gMjQgNDQgQyAyOS4xNjQwNjMgNDQgMzMuODU5Mzc1IDQyLjAyMzQzOCAzNy40MTAxNTYgMzguODA4NTk0IEwgMzEuMjE4NzUgMzMuNTcwMzEzIEMgMjkuMjEwOTM4IDM1LjA4OTg0NCAyNi43MTQ4NDQgMzYgMjQgMzYgQyAxOC43OTY4NzUgMzYgMTQuMzgyODEzIDMyLjY4MzU5NCAxMi43MTg3NSAyOC4wNTQ2ODggTCA2LjE5NTMxMyAzMy4wNzgxMjUgQyA5LjUwMzkwNiAzOS41NTQ2ODggMTYuMjI2NTYzIDQ0IDI0IDQ0IFogIj48L3BhdGg+PHBhdGggc3R5bGU9IiBmaWxsOiMxOTc2RDI7IiBkPSJNIDQzLjYwOTM3NSAyMC4wODIwMzEgTCA0MiAyMC4wODIwMzEgTCA0MiAyMCBMIDI0IDIwIEwgMjQgMjggTCAzNS4zMDQ2ODggMjggQyAzNC41MTE3MTkgMzAuMjM4MjgxIDMzLjA3MDMxMyAzMi4xNjQwNjMgMzEuMjE0ODQ0IDMzLjU3MDMxMyBDIDMxLjIxODc1IDMzLjU3MDMxMyAzMS4yMTg3NSAzMy41NzAzMTMgMzEuMjE4NzUgMzMuNTcwMzEzIEwgMzcuNDEwMTU2IDM4LjgwODU5NCBDIDM2Ljk3MjY1NiAzOS4yMDMxMjUgNDQgMzQgNDQgMjQgQyA0NCAyMi42NjAxNTYgNDMuODYzMjgxIDIxLjM1MTU2MyA0My42MDkzNzUgMjAuMDgyMDMxIFogIj48L3BhdGg+PC9nPjwvc3ZnPg==",
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
			LogoURL:           "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHg9IjBweCIgeT0iMHB4IgogICAgIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIKICAgICB2aWV3Qm94PSIwIDAgMzIgMzIiCiAgICAgc3R5bGU9ImZpbGw6IzAwMDAwMDsiPjxnIGlkPSJzdXJmYWNlMSI+PHBhdGggc3R5bGU9IiBmaWxsLXJ1bGU6ZXZlbm9kZDsiIGQ9Ik0gMTYgNCBDIDkuMzcxMDk0IDQgNCA5LjM3MTA5NCA0IDE2IEMgNCAyMS4zMDA3ODEgNy40Mzc1IDI1LjgwMDc4MSAxMi4yMDcwMzEgMjcuMzg2NzE5IEMgMTIuODA4NTk0IDI3LjQ5NjA5NCAxMy4wMjczNDQgMjcuMTI4OTA2IDEzLjAyNzM0NCAyNi44MDg1OTQgQyAxMy4wMjczNDQgMjYuNTIzNDM4IDEzLjAxNTYyNSAyNS43Njk1MzEgMTMuMDExNzE5IDI0Ljc2OTUzMSBDIDkuNjcxODc1IDI1LjQ5MjE4OCA4Ljk2ODc1IDIzLjE2MDE1NiA4Ljk2ODc1IDIzLjE2MDE1NiBDIDguNDIxODc1IDIxLjc3MzQzOCA3LjYzNjcxOSAyMS40MDIzNDQgNy42MzY3MTkgMjEuNDAyMzQ0IEMgNi41NDY4NzUgMjAuNjYwMTU2IDcuNzE4NzUgMjAuNjc1NzgxIDcuNzE4NzUgMjAuNjc1NzgxIEMgOC45MjE4NzUgMjAuNzYxNzE5IDkuNTU0Njg4IDIxLjkxMDE1NiA5LjU1NDY4OCAyMS45MTAxNTYgQyAxMC42MjUgMjMuNzQ2MDk0IDEyLjM2MzI4MSAyMy4yMTQ4NDQgMTMuMDQ2ODc1IDIyLjkxMDE1NiBDIDEzLjE1NjI1IDIyLjEzMjgxMyAxMy40Njg3NSAyMS42MDU0NjkgMTMuODA4NTk0IDIxLjMwNDY4OCBDIDExLjE0NDUzMSAyMS4wMDM5MDYgOC4zNDM3NSAxOS45NzI2NTYgOC4zNDM3NSAxNS4zNzUgQyA4LjM0Mzc1IDE0LjA2MjUgOC44MTI1IDEyLjk5MjE4OCA5LjU3ODEyNSAxMi4xNTIzNDQgQyA5LjQ1NzAzMSAxMS44NTE1NjMgOS4wNDI5NjkgMTAuNjI4OTA2IDkuNjk1MzEzIDguOTc2NTYzIEMgOS42OTUzMTMgOC45NzY1NjMgMTAuNzAzMTI1IDguNjU2MjUgMTIuOTk2MDk0IDEwLjIwNzAzMSBDIDEzLjk1MzEyNSA5Ljk0MTQwNiAxNC45ODA0NjkgOS44MDg1OTQgMTYgOS44MDQ2ODggQyAxNy4wMTk1MzEgOS44MDg1OTQgMTguMDQ2ODc1IDkuOTQxNDA2IDE5LjAwMzkwNiAxMC4yMDcwMzEgQyAyMS4yOTY4NzUgOC42NTYyNSAyMi4zMDA3ODEgOC45NzY1NjMgMjIuMzAwNzgxIDguOTc2NTYzIEMgMjIuOTU3MDMxIDEwLjYyODkwNiAyMi41NDY4NzUgMTEuODUxNTYzIDIyLjQyMTg3NSAxMi4xNTIzNDQgQyAyMy4xOTE0MDYgMTIuOTkyMTg4IDIzLjY1MjM0NCAxNC4wNjI1IDIzLjY1MjM0NCAxNS4zNzUgQyAyMy42NTIzNDQgMTkuOTg0Mzc1IDIwLjg0NzY1NiAyMC45OTYwOTQgMTguMTc1NzgxIDIxLjI5Njg3NSBDIDE4LjYwNTQ2OSAyMS42NjQwNjMgMTguOTg4MjgxIDIyLjM5ODQzOCAxOC45ODgyODEgMjMuNTE1NjI1IEMgMTguOTg4MjgxIDI1LjEyMTA5NCAxOC45NzY1NjMgMjYuNDE0MDYzIDE4Ljk3NjU2MyAyNi44MDg1OTQgQyAxOC45NzY1NjMgMjcuMTI4OTA2IDE5LjE5MTQwNiAyNy41MDM5MDYgMTkuODAwNzgxIDI3LjM4NjcxOSBDIDI0LjU2NjQwNiAyNS43OTY4NzUgMjggMjEuMzAwNzgxIDI4IDE2IEMgMjggOS4zNzEwOTQgMjIuNjI4OTA2IDQgMTYgNCBaICI+PC9wYXRoPjwvZz48L3N2Zz4=",
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

	statusCode, body, err := s.GetRawProfile(provider, code)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New("Failed to request GET profile. Status Code: %d. Body: %s", statusCode, body)
	}

	return s.ParseRawProfile(provider, body)
}

//GetRawProfile returns raw JSON response from Profile API
func (s *OAuthService) GetRawProfile(provider string, code string) (int, string, error) {
	config, err := s.getConfig(provider)
	if err != nil {
		return 0, "", err
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
		rErr := err.(*oauth2.RetrieveError)
		return rErr.Response.StatusCode, string(rErr.Body), nil
	}

	return s.doGet(config.ProfileURL, oauthToken.AccessToken)
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

	if !validate.Email(profile.Email).Ok {
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
			LogoURL:          p.LogoURL,
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
