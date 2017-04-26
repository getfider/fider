package handlers_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/WeCanHearYou/wechy/app/handlers"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

var (
	oauthService  = &oauth.MockOAuthService{}
	tenants       = &inmemory.TenantStorage{}
	users         = &inmemory.UserStorage{}
	oauthHandlers = handlers.OAuth(tenants, oauthService, users)
)

func TestLoginHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_AUTH_ENDPOINT", "http://login.test.canherayou.com:3000")
	code, response := server.ExecuteRaw(oauthHandlers.Login(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")
	code, _ := server.ExecuteRaw(oauthHandlers.Callback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com")
	code, response := server.ExecuteRaw(oauthHandlers.Callback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	tenant := &models.Tenant{ID: 1, Name: "Demonstration", Subdomain: "demo"}
	tenants.Add(tenant)

	users.Register(&models.User{
		ID:     300,
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Tenant: tenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "FB1234", Name: oauth.FacebookProvider},
		},
	})

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://demo.test.canherayou.com/oauth/callback?state=http://demo.test.canhearyou.com&code=123")
	code, response := server.ExecuteRaw(oauthHandlers.Callback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://demo.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.6_dLZrulH37ymBtqy-l7bhCti9hBv0lgEhH8tLm07CI"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)
	tenant := &models.Tenant{ID: 2, Name: "Orange", Subdomain: "orange"}
	tenants.Add(tenant)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com&code=456")
	code, response := server.ExecuteRaw(oauthHandlers.Callback(oauth.FacebookProvider))

	user, err := users.GetByEmail(tenant.ID, "some.guy@facebook.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Facebook Guy"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.PGavs5a6HRotRozfXfNP39JPb0vSus_8LL9MAOeLGDs"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.canhearyou.com&code=123")
	code, response := server.ExecuteRaw(oauthHandlers.Callback(oauth.GoogleProvider))

	tenant, _ := tenants.GetByDomain("demo")
	user, err := users.GetByEmail(tenant.ID, "jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(len(user.Providers)).To(Equal(2))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://demo.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.6_dLZrulH37ymBtqy-l7bhCti9hBv0lgEhH8tLm07CI"))
}
