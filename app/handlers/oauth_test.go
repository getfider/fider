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

	server, _ := mock.NewServer()
	code, response := server.Execute(handlers.Login(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")

	Expect(func() {
		server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))
	}).To(Panic())
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://demo.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.fider.io&code=456")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail(mock.OrangeTenant.ID, "some.guy@facebook.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Facebook Guy"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.ydyGDIZZHbJ-mgvpAsXTZnbs1rBH6cTVjHctZEACUOo"))
}

func TestCallbackHandler_NewUserWithoutEmail(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.Users.Register(&models.User{
		Name:   "Some Guy",
		Email:  "",
		Tenant: mock.DemoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "GO999", Name: oauth.GoogleProvider},
		},
	})

	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=798")
	code, response := server.Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByID(3)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Guy"))
	Expect(len(user.Providers)).To(Equal(1))

	user, err = services.Users.GetByID(4)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Mark"))
	Expect(len(user.Providers)).To(Equal(1))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjo0LCJ1c2VyL25hbWUiOiJNYXJrIiwidXNlci9lbWFpbCI6IiJ9.G93ZTFcDuHiIlYbDvMnjhoDebeJZMifWbv9v0rayQOI"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.fider.io&code=123")
	code, response := server.Execute(handlers.OAuthCallback(oauth.GoogleProvider))

	user, err := services.Users.GetByEmail(mock.DemoTenant.ID, "jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(len(user.Providers)).To(Equal(2))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}
