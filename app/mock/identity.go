package mock

import (
	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/identity"
)

// TenantService implements a mocked TenantService
type TenantService struct {
}

// GetByDomain returns a tenant based on its domain
func (svc TenantService) GetByDomain(domain string) (*app.Tenant, error) {
	tenant := &app.Tenant{}

	switch domain {
	case "demo":
		tenant.ID = 1
		tenant.Name = "Demonstration"
		tenant.Domain = "demo"
		return tenant, nil
	case "orange.test.canherayou.com":
		tenant.ID = 2
		tenant.Name = "Orange Inc."
		tenant.Domain = "orange."
		return tenant, nil
	}
	return nil, app.ErrNotFound
}

//OAuthService implements a mocked OAuthService
type OAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p OAuthService) GetAuthURL(provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p OAuthService) GetProfile(provider string, code string) (*identity.OAuthUserProfile, error) {
	return nil, nil
}

// UserService implements a mocked UserService
type UserService struct {
}

// GetByEmail returns a user based on given email
func (svc UserService) GetByEmail(email string) (*app.User, error) {
	user := &app.User{}
	return user, nil
}

// Register creates a new user based on given information
func (svc UserService) Register(user *app.User) error {
	return nil
}
