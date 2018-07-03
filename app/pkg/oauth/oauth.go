package oauth

import (
	"errors"
	"os"
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

//IsProviderEnabled returns true if provider is enabled
func IsProviderEnabled(name string) bool {
	if name == GoogleProvider {
		return os.Getenv("OAUTH_GOOGLE_CLIENTID") != ""
	} else if name == FacebookProvider {
		return os.Getenv("OAUTH_FACEBOOK_APPID") != ""
	} else if name == GitHubProvider {
		return os.Getenv("OAUTH_GITHUB_CLIENTID") != ""
	}
	return false
}

//Service provides OAuth operations
type Service interface {
	GetAuthURL(provider string, redirect string) (string, error)
	GetProfile(provider string, code string) (*UserProfile, error)
}
