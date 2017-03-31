package identity_test

import (
	"testing"

	"net/http"

	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/mock"
	. "github.com/onsi/gomega"
)

//HTTPOAuthService implements real OAuth operations using Golang's oauth2 package
type mockOAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p mockOAuthService) GetAuthURL(provider string, redirect string) string {
	return "http://myapp/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p mockOAuthService) GetProfile(provider string, code string) (*identity.OAuthUserProfile, error) {
	return nil, nil
}

func TestLoginHandlers(t *testing.T) {
	RegisterTestingT(t)

	oauth := &mockOAuthService{}

	server := mock.NewServer()
	code, response := server.ExecuteRaw(identity.OAuth(nil, oauth, nil).Login(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://myapp/oauth/token?provider=facebook&redirect="))
}
