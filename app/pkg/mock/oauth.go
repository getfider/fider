package mock

import "github.com/getfider/fider/app/pkg/oauth"

//OAuthService implements a mocked OAuthService
type OAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (s *OAuthService) GetAuthURL(provider string, redirect string) (string, error) {
	return "http://avengers.test.fider.io/oauth/token?provider=" + provider + "&redirect=" + redirect, nil
}

//GetProfile returns user profile based on provider and code
func (s *OAuthService) GetProfile(provider string, code string) (*oauth.UserProfile, error) {
	if provider == "facebook" && code == "123" {
		return &oauth.UserProfile{
			ID:    "FB123",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	if provider == "facebook" && code == "456" {
		return &oauth.UserProfile{
			ID:    "FB456",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}, nil
	}

	if provider == "facebook" && code == "798" {
		return &oauth.UserProfile{
			ID:    "FB798",
			Name:  "Mark",
			Email: "",
		}, nil
	}

	if provider == "google" && code == "123" {
		return &oauth.UserProfile{
			ID:    "GO123",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	return nil, nil
}

//ListProviders returns a list of all providers for current tenant
func (s *OAuthService) ListProviders() ([]*oauth.ProviderOption, error) {
	return []*oauth.ProviderOption{}, nil
}
