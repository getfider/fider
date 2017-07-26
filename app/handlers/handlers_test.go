package handlers_test

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/inmemory"
)

var demoTenant *models.Tenant
var orangeTenant *models.Tenant

func getEmptyServices() *app.Services {
	return &app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
		OAuth:   &MockOAuthService{},
	}
}

func getServices() *app.Services {
	services := &app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
		OAuth:   &MockOAuthService{},
	}

	demoTenant, _ = services.Tenants.Add("Demonstration", "demo")
	orangeTenant, _ = services.Tenants.Add("Orange Inc.", "orange")

	return services
}

//MockOAuthService implements a mocked OAuthService
type MockOAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p *MockOAuthService) GetAuthURL(authEndpoint string, provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p *MockOAuthService) GetProfile(authEndpoint string, provider string, code string) (*oauth.UserProfile, error) {
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
