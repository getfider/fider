package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app"

	"net/http"
	"net/http/httptest"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

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
	request.ExpiresAt = request.CreatedAt.Add(-6 * time.Minute) //reduce 1 minute

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusGone)
}

func TestVerifySignInKeyHandler_CorrectKey_ExistingUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == mock.JonSnow.Email {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, &models.SignInByEmail{
		Email: "jon.snow@got.com",
	})

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

	bus.AddHandler(func(ctx context.Context, q *query.SearchPosts) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.CountPostPerStatus) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAllTags) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

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

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

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

	user := &models.User{
		Name:   "Hot Pie",
		Email:  "hot.pie@got.com",
		Tenant: mock.DemoTenant,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == user.Email {
			q.Result = user
			return nil
		}
		return app.ErrNotFound
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

	user := &models.User{
		Name:   "Hot Pie",
		Email:  "hot.pie@got.com",
		Tenant: mock.DemoTenant,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == user.Email {
			q.Result = user
			return nil
		}
		return app.ErrNotFound
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

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	e := &models.UserInvitation{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/invite/verify?k=1234567890").
		Execute(handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignUpKeyHandler_PendingTenant(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()

	var newUser *models.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		newUser = c.User
		return nil
	})

	activated := false
	bus.AddHandler(func(ctx context.Context, c *cmd.ActivateTenant) error {
		activated = true
		return nil
	})

	e := &models.CreateTenant{Email: "hot.pie@got.com", Name: "Hot Pie"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)
	mock.DemoTenant.Status = models.TenantPending

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signup/verify?k=1234567890").
		Execute(handlers.VerifySignUpKey())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	Expect(newUser.Name).Equals("Hot Pie")
	Expect(newUser.Email).Equals("hot.pie@got.com")
	Expect(activated).IsTrue()
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

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == mock.JonSnow.Email {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

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

	var newUser *models.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		newUser = c.User
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	e := &models.SignInByEmail{Email: "hot.pie@got.com"}
	services.Tenants.SaveVerificationKey("1234567890", 15*time.Minute, e)

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ "name": "Hot Pie", "key": "1234567890" }`)
	Expect(code).Equals(http.StatusOK)

	Expect(newUser.Name).Equals("Hot Pie")
	Expect(newUser.Email).Equals("hot.pie@got.com")

	ExpectFiderAuthCookie(response, newUser)

	request, err := services.Tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "1234567890")
	Expect(err).IsNil()
	Expect(request.VerifiedAt).IsNotNil()
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
	cookies := response.HeaderMap["Set-Cookie"]
	if expected == nil {
		for _, c := range cookies {
			cookie := web.ParseCookie(c)
			Expect(cookie.Name).NotEquals(web.CookieAuthName)
		}
	} else {
		for _, c := range cookies {
			cookie := web.ParseCookie(c)
			if cookie.Name == web.CookieAuthName {
				Expect(cookie.Name).Equals(web.CookieAuthName)
				ExpectFiderToken(cookie.Value, expected)
				Expect(cookie.HttpOnly).IsTrue()
				Expect(cookie.Path).Equals("/")
				Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(365*24*time.Hour), 5*time.Second)
				return
			}
		}
		panic("Cookie not found...")
	}
}

func ExpectFiderToken(token string, expected *models.User) {
	user, err := jwt.DecodeFiderClaims(token)
	Expect(err).IsNil()
	Expect(user.UserID).Equals(expected.ID)
	Expect(user.UserName).Equals(expected.Name)
	Expect(user.UserEmail).Equals(expected.Email)
}
