package handler_test

import (
	"testing"

	"net/http"

	"github.com/WeCanHearYou/wchy/auth"
	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/mock"
	. "github.com/onsi/gomega"
)

//HTTPOAuthService implements real OAuth operations using Golang's oauth2 package
type mockOAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p mockOAuthService) GetAuthURL(provider string, redirect string) string {
	return "http://myapp/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p mockOAuthService) GetProfile(provider string, code string) (*auth.OAuthUserProfile, error) {
	return nil, nil
}

func TestLoginHandlers(t *testing.T) {
	RegisterTestingT(t)

	ctx := &context.WchyContext{
		OAuth: &mockOAuthService{},
	}

	server := mock.NewServer()
	code, response := server.ExecuteRaw(handler.OAuth(ctx).Login(auth.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://myapp/oauth/token?provider=facebook&redirect="))
}
