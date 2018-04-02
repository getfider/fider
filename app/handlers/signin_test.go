package handlers_test

import (
	"fmt"
	"testing"
	"time"

	"net/http"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestSignOutHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signout?redirect=/").
		AddCookie(web.CookieAuthName, "some-value").
		Execute(handlers.SignOut())

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("/"))
	Expect(response.Header().Get("Set-Cookie")).To(ContainSubstring(web.CookieAuthName + "=;"))
}

func TestSignInHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.Execute(handlers.SignIn(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://avengers.test.fider.io/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=abc").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://avengers.test.fider.io"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/callback?state=http://demo.test.fider.io&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}

func TestCallbackHandler_SignUp(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/callback?state=http://demo.test.fider.io/signup&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io/signup?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvYXV0aC9pZCI6IkZCMTIzIiwib2F1dGgvcHJvdmlkZXIiOiJmYWNlYm9vayIsIm9hdXRoL25hbWUiOiJKb24gU25vdyIsIm9hdXRoL2VtYWlsIjoiam9uLnNub3dAZ290LmNvbSJ9.AFMEtMLxd2nAxPzMVeGVrJomn-n56WdhLvgYSmSq008"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io&code=456").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail("some.guy@facebook.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Facebook Guy"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://avengers.test.fider.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.ydyGDIZZHbJ-mgvpAsXTZnbs1rBH6cTVjHctZEACUOo"))
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

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://demo.test.fider.io&code=798").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	user, err := services.Users.GetByID(3)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Some Guy"))
	Expect(len(user.Providers)).To(Equal(1))

	user, err = services.Users.GetByID(4)
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Mark"))
	Expect(len(user.Providers)).To(Equal(1))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjo0LCJ1c2VyL25hbWUiOiJNYXJrIiwidXNlci9lbWFpbCI6IiJ9.G93ZTFcDuHiIlYbDvMnjhoDebeJZMifWbv9v0rayQOI"))
}

func TestCallbackHandler_ExistingUser_WithoutEmail(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.Users.Register(&models.User{
		Name:   "Some Facebook Guy",
		Email:  "",
		Tenant: mock.DemoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "FB456", Name: oauth.FacebookProvider},
		},
	})

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://demo.test.fider.io&code=456").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	_, err := services.Users.GetByID(4)
	Expect(errors.Cause(err)).To(Equal(app.ErrNotFound))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiIifQ.DBcMmKSrGxiXGVCq9QT716xfw3M4kzP-Njazsk6vaFI"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://demo.test.fider.io&code=123").
		Execute(handlers.OAuthCallback(oauth.GoogleProvider))

	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(len(user.Providers)).To(Equal(2))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJKb24gU25vdyIsInVzZXIvZW1haWwiOiJqb24uc25vd0Bnb3QuY29tIn0.S7P8zTU0rVovmchNbwamBewYbO96GdJcOygn7tbsikw"))
}

func TestSignInByEmailHandler_WithoutEmail(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		ExecutePost(handlers.SignInByEmail(), "{ }")

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestSignInByEmailHandler_WithEmail(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmail(), `{ "email": "jon.snow@got.com" }`)

	Expect(code).To(Equal(http.StatusOK))
}

func TestVerifySignInKeyHandler_UnknownKey(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=unknown").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io"))
}

func TestVerifySignInKeyHandler_UsedKey(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)
	services.Tenants.SetKeyAsVerified("1234567890")

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io"))
}

func TestVerifySignInKeyHandler_ExpiredKey(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 5*time.Minute, e)
	request, _ := services.Tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "1234567890")
	request.ExpiresOn = request.CreatedOn.Add(-6 * time.Minute) //reduce 1 minute

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io"))
}

func TestVerifySignInKeyHandler_CorrectKey_ExistingUser(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	token, _ := jwt.Encode(models.FiderClaims{
		UserID:    mock.JonSnow.ID,
		UserName:  mock.JonSnow.Name,
		UserEmail: mock.JonSnow.Email,
	})

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io"))
	Expect(response.Header().Get("Set-Cookie")).To(ContainSubstring(fmt.Sprintf("%s=%s;", web.CookieAuthName, token)))
}

func TestVerifySignInKeyHandler_CorrectKey_NewUser(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).To(Equal(http.StatusOK))
}

func TestVerifySignUpKeyHandler_InactiveTenant(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.CreateTenant{Email: "hot.pie@got.com", Name: "Hot Pie"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)
	mock.DemoTenant.Status = models.TenantInactive

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signup/verify?k=1234567890").
		Execute(handlers.VerifySignUpKey())

	tenant, _ := services.Tenants.GetByDomain("demo")

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(tenant.Status).To(Equal(models.TenantActive))
	Expect(response.Header().Get("Location")).To(Equal("http://demo.test.fider.io"))
}

func TestCompleteSignInProfileHandler_UnknownKey(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{  }`)

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestCompleteSignInProfileHandler_ExistingUser_CorrectKey(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: mock.JonSnow.Email}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ "name": "Hot Pie", "key": "1234567890" }`)

	Expect(code).To(Equal(http.StatusOK))
}

func TestCompleteSignInProfileHandler_CorrectKey(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ "name": "Hot Pie", "key": "1234567890" }`)
	Expect(code).To(Equal(http.StatusOK))

	user, err := services.Users.GetByEmail("hot.pie@got.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Hot Pie"))
	Expect(user.Email).To(Equal("hot.pie@got.com"))

	token, _ := jwt.Encode(models.FiderClaims{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
	})
	Expect(response.Header().Get("Set-Cookie")).To(ContainSubstring(fmt.Sprintf("%s=%s;", web.CookieAuthName, token)))

	request, err := services.Tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "1234567890")
	Expect(err).To(BeNil())
	Expect(request.VerifiedOn).NotTo(BeNil())
}
