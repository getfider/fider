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
	ID    string
	Name  string
	Email string
}

//ProviderOption represents an OAuth provider that can be used to authenticate
type ProviderOption struct {
	Provider         string `json:"provider"`
	DisplayName      string `json:"displayName"`
	ClientID         string `json:"clientId"`
	URL              string `json:"url"`
	CallbackURL      string `json:"callbackUrl"`
	LogoURL          string `json:"logoUrl"`
	IsCustomProvider bool   `json:"isCustomProvider"`
}

//Service provides OAuth operations
type Service interface {
	GetAuthURL(provider, redirect, identifier string) (string, error)
	GetProfile(provider string, code string) (*UserProfile, error)
	ListActiveProviders() ([]*ProviderOption, error)
	ListAllProviders() ([]*ProviderOption, error)
}
