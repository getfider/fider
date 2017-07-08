package handlers_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

var services = &app.Services{
	Tenants: &inmemory.TenantStorage{},
	Users:   &inmemory.UserStorage{},
	OAuth:   &oauth.MockOAuthService{},
}

func TestLoginHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(services)
	code, response := server.Execute(handlers.Login(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")
	code, _ := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	tenant, _ := services.Tenants.Add("Demonstration", "demo")

	services.Users.Register(&models.User{
		ID:     300,
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Tenant: tenant,
		Providers: []*models.UserProvider{
			{UID: "FB1234", Name: oauth.FacebookProvider},
		},
	})

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://demo.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.n73-Qqhb_6lIALS_H3UzrEuzLTHgtS6nCkG_Ua95v7s"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)
	tenant, _ := services.Tenants.Add("Orange", "orange")

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io&code=456")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail(tenant.ID, "some.guy@facebook.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Facebook Guy"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.QklasoVtsMu_urLZsJv0CDu0VuMFI2CNC78ckx0nSfQ"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.GoogleProvider))

	tenant, _ := services.Tenants.GetByDomain("demo")
	user, err := services.Users.GetByEmail(tenant.ID, "jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(len(user.Providers)).To(Equal(2))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.n73-Qqhb_6lIALS_H3UzrEuzLTHgtS6nCkG_Ua95v7s"))
}
