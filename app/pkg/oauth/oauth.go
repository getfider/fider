package oauth

import (
	"errors"
)

const (
	//FacebookProvider is const for 'facebook'
	FacebookProvider = "facebook"
	//GoogleProvider is const for 'google'
	GoogleProvider = "google"
	//GitHubProvider is const for 'github'
	GitHubProvider = "github"
)

//ErrUserIDRequired is used when OAuth integration returns an empty user ID
var ErrUserIDRequired = errors.New("UserID is required during OAuth integration")

//UserProfile represents an OAuth user profile
type UserProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//ProviderOption represents an OAuth provider that can be used to authenticate
type ProviderOption struct {
	Provider         string `json:"provider"`
	DisplayName      string `json:"displayName"`
	ClientID         string `json:"clientID"`
	URL              string `json:"url"`
	CallbackURL      string `json:"callbackURL"`
	LogoID           int    `json:"logoID"`
	IsCustomProvider bool   `json:"isCustomProvider"`
	IsEnabled        bool   `json:"isEnabled"`
}

//Service provides OAuth operations
type Service interface {
	GetAuthURL(provider, redirect, identifier string) (string, error)
	GetProfile(provider string, code string) (*UserProfile, error)
	GetRawProfile(provider string, code string) (string, error)
	ParseRawProfile(provider, body string) (*UserProfile, error)
	ListActiveProviders() ([]*ProviderOption, error)
	ListAllProviders() ([]*ProviderOption, error)
}
