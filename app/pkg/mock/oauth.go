package mock

import "github.com/getfider/fider/app/pkg/oauth"

//OAuthService implements a mocked OAuthService
type OAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p *OAuthService) GetAuthURL(authEndpoint string, provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p *OAuthService) GetProfile(authEndpoint string, provider string, code string) (*oauth.UserProfile, error) {
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

	if provider == "google" && code == "456" {
		return &oauth.UserProfile{
			ID:    "GO456",
			Name:  "Bob",
			Email: "",
		}, nil
	}

	return nil, nil
}
