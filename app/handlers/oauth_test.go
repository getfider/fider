package handlers_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/oauth"
	. "github.com/onsi/gomega"
)

func TestLoginHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(getServices())
	code, response := server.Execute(handlers.Login(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(getServices())
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")
	code, _ := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(getServices())
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)
	services := getServices()

	services.Users.Register(&models.User{
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Tenant: demoTenant,
		Providers: []*models.UserProvider{
			{UID: "FB1234", Name: oauth.FacebookProvider},
		},
	})

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://demo.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)
	services := getServices()

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io&code=456")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail(orangeTenant.ID, "some.guy@facebook.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Facebook Guy"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.QklasoVtsMu_urLZsJv0CDu0VuMFI2CNC78ckx0nSfQ"))
}

func TestCallbackHandler_NewUserWithoutEmail(t *testing.T) {
	RegisterTestingT(t)
	services := getServices()
	services.Users.Register(&models.User{
		Name:   "Some Guy",
		Email:  "",
		Tenant: demoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "GO999", Name: oauth.GoogleProvider},
		},
	})

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=798")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByID(1)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Guy"))
	Expect(len(user.Providers)).To(Equal(1))

	user, err = services.Users.GetByID(2)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Mark"))
	Expect(len(user.Providers)).To(Equal(1))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoyLCJ1c2VyL25hbWUiOiJNYXJrIiwidXNlci9lbWFpbCI6IiJ9.NKMh8HaRud9wOtFznIWZuJabsYFm8UWy3jPykJ7U79E"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)
	services := getServices()
	services.Users.Register(&models.User{
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Tenant: demoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "FB123", Name: oauth.FacebookProvider},
		},
	})

	server := mock.NewServer()
	server.Context.SetServices(services)
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.GoogleProvider))

	user, err := services.Users.GetByEmail(demoTenant.ID, "jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(len(user.Providers)).To(Equal(2))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}
