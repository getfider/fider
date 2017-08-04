package oauth

import (
	"encoding/json"
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

//UserProfile represents an OAuth user profile
type UserProfile struct {
	ID    json.Number
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
	GetAuthURL(authEndpoint string, provider string, redirect string) string
	GetProfile(authEndpoint string, provider string, code string) (*UserProfile, error)
}
