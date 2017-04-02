package identity_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/WeCanHearYou/wechy/app/dbx"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

//OAuthService implements a mocked OAuthService
type OAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p OAuthService) GetAuthURL(provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p OAuthService) GetProfile(provider string, code string) (*identity.OAuthUserProfile, error) {
	if provider == "facebook" && code == "123" {
		return &identity.OAuthUserProfile{
			ID:    "FB1234",
			Name:  "Jon Snow",
			Email: "jon.sno@got.com",
		}, nil
	}

	if provider == "facebook" && code == "456" {
		return &identity.OAuthUserProfile{
			ID:    "FB5678",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}, nil
	}

	return nil, nil
}

func handlers() identity.OAuthHandlers {
	db, _ := dbx.New()

	oauth := &OAuthService{}
	tenant := &postgres.TenantService{DB: db}
	user := &postgres.UserService{DB: db}
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

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com")
	code, response := server.ExecuteRaw(handlers().Callback(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com&code=123")
	code, response := server.ExecuteRaw(handlers().Callback(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vQGdvdC5jb20iLCJ0ZW5hbnQvaWQiOjQwMH0.nxu7QHFXTeYh_ObpKYV6e3p0kk1mkGbNGplITZImSs8"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com&code=456")
	code, response := server.ExecuteRaw(handlers().Callback(identity.OAuthFacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20iLCJ0ZW5hbnQvaWQiOjQwMH0.iWkpvh11QrUDJQnTk9PtqW6m48DnaHbsj-lbl6feK-Q"))
}
