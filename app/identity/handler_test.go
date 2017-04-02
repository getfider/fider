package identity_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/mock"
	. "github.com/onsi/gomega"
)

func handlers() identity.OAuthHandlers {
	oauth := &mock.OAuthService{}
	tenant := &mock.TenantService{}
	user := &mock.UserService{}
	return identity.OAuth(tenant, oauth, user)
}

func TestLoginHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	code, response := server.ExecuteRaw(handlers().Login(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")
	code, _ := server.ExecuteRaw(handlers().Callback(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}
