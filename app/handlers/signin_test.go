package handlers_test

import (
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSignOutHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signout?redirect=/").
		AddCookie(web.CookieAuthName, "some-value").
		Execute(handlers.SignOut())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=; Path=/; Expires=")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring("Max-Age=0; HttpOnly")
}

func TestSignInByOAuthHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.Execute(handlers.SignInByOAuth(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/token?provider=facebook&redirect=")
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=abc").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusInternalServerError)
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io")
}

func TestCallbackHandler_SignIn(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/facebook/token?code=123&path=")
}

func TestCallbackHandler_SignIn_WithPath(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewServer()

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io/some-page&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/facebook/token?code=123&path=%2Fsome-page")
}

func TestCallbackHandler_SignUp(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://demo.test.fider.io/signup&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))
	Expect(code).Equals(http.StatusTemporaryRedirect)

	location, _ := url.Parse(response.Header().Get("Location"))
	Expect(location.Host).Equals("demo.test.fider.io")
	Expect(location.Scheme).Equals("http")
	Expect(location.Path).Equals("/signup")
	ExpectOAuthToken(location.Query().Get("token"), &jwt.OAuthClaims{
		OAuthProvider: "facebook",
		OAuthID:       "FB123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
	})
}

func TestOAuthTokenHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=123").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestOAuthTokenHandler_NewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456&path=/hello").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io/hello")

	user, err := services.Users.GetByEmail("some.guy@facebook.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Some Facebook Guy")

	ExpectFiderAuthCookie(response, user)
}

func TestOAuthTokenHandler_NewUserWithoutEmail(t *testing.T) {
	RegisterT(t)

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
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=798").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	user, err := services.Users.GetByID(3)
	Expect(err).IsNil()
	Expect(user.ID).Equals(3)
	Expect(user.Name).Equals("Some Guy")
	Expect(user.Providers).HasLen(1)

	user, err = services.Users.GetByID(4)
	Expect(err).IsNil()
	Expect(user.ID).Equals(4)
	Expect(user.Name).Equals("Mark")
	Expect(user.Providers).HasLen(1)

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, &models.User{
		ID:   4,
		Name: "Mark",
	})
}

func TestOAuthTokenHandler_ExistingUser_WithoutEmail(t *testing.T) {
	RegisterT(t)

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
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)

	_, err := services.Users.GetByID(4)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, &models.User{
		ID:   3,
		Name: "Some Facebook Guy",
	})
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=123").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.GoogleProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)

	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Providers).HasLen(2)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestCallbackHandler_NewUser_PrivateTenant(t *testing.T) {
	RegisterT(t)
	server, services := mock.NewServer()
	mock.AvengersTenant.IsPrivate = true

	code, response := server.
		WithURL("http://ideas.theavengers.com/oauth/facebook/token?code=456").
		OnTenant(mock.AvengersTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail("some.guy@facebook.com")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(user).IsNil()

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://ideas.theavengers.com/not-invited")
	ExpectFiderAuthCookie(response, nil)
}

func TestSignInByEmailHandler_WithoutEmail(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		ExecutePost(handlers.SignInByEmail(), "{ }")

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSignInByEmailHandler_WithEmail(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmail(), `{ "email": "jon.snow@got.com" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignInKeyHandler_UnknownKey(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=unknown").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusNotFound)
}

func TestVerifySignInKeyHandler_UsedKey(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)
	services.Tenants.SetKeyAsVerified("1234567890")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusGone)
}

func TestVerifySignInKeyHandler_ExpiredKey(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 5*time.Minute, e)
	request, _ := services.Tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "1234567890")
	request.ExpiresOn = request.CreatedOn.Add(-6 * time.Minute) //reduce 1 minute

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusGone)
}

func TestVerifySignInKeyHandler_CorrectKey_ExistingUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "jon.snow@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")

	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestVerifySignInKeyHandler_CorrectKey_NewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignInKeyHandler_PrivateTenant_SignInRequest_NonInviteNewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusForbidden)
}

func TestVerifySignInKeyHandler_PrivateTenant_SignInRequest_RegisteredUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	services.Users.Register(&models.User{
		Name:   "Hot Pie",
		Email:  "hot.pie@got.com",
		Tenant: mock.DemoTenant,
	})

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestVerifySignInKeyHandler_PrivateTenant_InviteRequest_ExistingUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	services.Users.Register(&models.User{
		Name:   "Hot Pie",
		Email:  "hot.pie@got.com",
		Tenant: mock.DemoTenant,
	})

	e := &models.UserInvitation{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/invite/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestVerifySignInKeyHandler_PrivateTenant_InviteRequest_NewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	e := &models.UserInvitation{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/invite/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignUpKeyHandler_InactiveTenant(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.CreateTenant{Email: "hot.pie@got.com", Name: "Hot Pie"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)
	mock.DemoTenant.Status = models.TenantInactive

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signup/verify?k=1234567890").
		Execute(handlers.VerifySignUpKey())

	tenant, _ := services.Tenants.GetByDomain("demo")

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(tenant.Status).Equals(models.TenantActive)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestCompleteSignInProfileHandler_UnknownKey(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{  }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCompleteSignInProfileHandler_ExistingUser_CorrectKey(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: mock.JonSnow.Email}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ "name": "Hot Pie", "key": "1234567890" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestCompleteSignInProfileHandler_CorrectKey(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ "name": "Hot Pie", "key": "1234567890" }`)
	Expect(code).Equals(http.StatusOK)

	user, err := services.Users.GetByEmail("hot.pie@got.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Hot Pie")
	Expect(user.Email).Equals("hot.pie@got.com")

	ExpectFiderAuthCookie(response, user)

	request, err := services.Tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "1234567890")
	Expect(err).IsNil()
	Expect(request.VerifiedOn).IsNotNil()
}

func TestSignInPageHandler_AuthenticatedUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithURL("http://demo.test.fider.io/signin").
		Execute(handlers.SignInPage())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestSignInPageHandler_NonPrivateTenant(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin").
		Execute(handlers.SignInPage())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestSignInPageHandler_PrivateTenant_UnauthenticatedUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin").
		Execute(handlers.SignInPage())

	Expect(code).Equals(http.StatusOK)
}

func ExpectFiderAuthCookie(response *httptest.ResponseRecorder, expected *models.User) {
	cookie := web.ParseCookie(response.Header().Get("Set-Cookie"))
	if expected == nil {
		Expect(cookie).IsNil()
	} else {
		Expect(cookie.Name).Equals(web.CookieAuthName)
		ExpectFiderToken(cookie.Value, expected)
		Expect(cookie.HttpOnly).IsTrue()
		Expect(cookie.Path).Equals("/")
		Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(365*24*time.Hour), 5*time.Second)
	}
}

func ExpectFiderToken(token string, expected *models.User) {
	user, err := jwt.DecodeFiderClaims(token)
	Expect(err).IsNil()
	Expect(user.UserID).Equals(expected.ID)
	Expect(user.UserName).Equals(expected.Name)
	Expect(user.UserEmail).Equals(expected.Email)
}

func ExpectOAuthToken(token string, expected *jwt.OAuthClaims) {
	user, err := jwt.DecodeOAuthClaims(token)
	Expect(err).IsNil()
	Expect(user.OAuthID).Equals(expected.OAuthID)
	Expect(user.OAuthName).Equals(expected.OAuthName)
	Expect(user.OAuthEmail).Equals(expected.OAuthEmail)
	Expect(user.OAuthProvider).Equals(expected.OAuthProvider)
}
