package handlers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/getfider/fider/app"

	"net/http"
	"net/http/httptest"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSignInByEmailHandler_WithoutEmail(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmail(), "{ }")

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSignInByEmailHandler_ExistingUser(t *testing.T) {
	RegisterT(t)

	var saveKeyCmd *cmd.SaveVerificationKey
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		saveKeyCmd = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == "jon.snow@got.com" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmail(), `{ "email": "jon.snow@got.com" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(saveKeyCmd.Key).HasLen(6)
	Expect(saveKeyCmd.Request.GetKind()).Equals(enum.EmailVerificationKindSignIn)
	Expect(saveKeyCmd.Request.GetEmail()).Equals("jon.snow@got.com")
	Expect(response.Body.String()).ContainsSubstring(`"userExists":true`)
}

func TestSignInByEmailHandler_NewUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmail(), `{ "email": "new.user@got.com" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(response.Body.String()).ContainsSubstring(`"userExists":false`)
}

func TestSignInByEmailWithNameHandler_NewUser(t *testing.T) {
	RegisterT(t)

	var saveKeyCmd *cmd.SaveVerificationKey
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		saveKeyCmd = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmailWithName(), `{ "email": "new.user@got.com", "name": "New User" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(saveKeyCmd.Key).HasLen(6)
	Expect(saveKeyCmd.Request.GetKind()).Equals(enum.EmailVerificationKindSignIn)
	Expect(saveKeyCmd.Request.GetEmail()).Equals("new.user@got.com")
	Expect(saveKeyCmd.Request.GetName()).Equals("New User")
}

func TestSignInByEmailWithNameHandler_ExistingUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == "jon.snow@got.com" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.SignInByEmailWithName(), `{ "email": "jon.snow@got.com", "name": "Jon Snow" }`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.Body.String()).ContainsSubstring("already exists")
}

func TestSignInByEmailWithNameHandler_PrivateTenant(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	privateTenant := &entity.Tenant{
		ID:        1,
		Name:      "Private Tenant",
		Subdomain: "private",
		IsPrivate: true,
	}

	server := mock.NewServer()
	code, _ := server.
		OnTenant(privateTenant).
		ExecutePost(handlers.SignInByEmailWithName(), `{ "email": "new.user@got.com", "name": "New User" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestVerifySignInKeyHandler_UnknownKey(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=unknown").
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusNotFound)
}

func TestVerifySignInKeyHandler_ExpiredKey(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignIn {
			now := time.Now().Add(-5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: now,
				Email:     "jon.snow@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	verified := false
	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		if c.Key == key {
			verified = true
		}
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusGone)
	Expect(verified).IsTrue()
}

func TestVerifySignInKeyHandler_KeyUsedLongAgo_ShouldReturnGone(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		Expect(q.Key).Equals(key)
		Expect(q.Kind).Equals(enum.EmailVerificationKindSignIn)

		verifiedAt := time.Now().Add(-1 * time.Hour)
		q.Result = &entity.EmailVerification{
			Key:        q.Key,
			Kind:       q.Kind,
			VerifiedAt: &verifiedAt,
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Email:      "jon.snow@got.com",
		}
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusGone)
}

func TestVerifySignInKeyHandler_CorrectKey_ExistingUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		Expect(q.Key).Equals(key)
		Expect(q.Kind).Equals(enum.EmailVerificationKindSignIn)

		expiresAt := time.Now().Add(5 * time.Minute)
		q.Result = &entity.EmailVerification{
			Key:       q.Key,
			Kind:      q.Kind,
			ExpiresAt: expiresAt,
			Email:     "jon.snow@got.com",
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		Expect(q.Email).Equals(mock.JonSnow.Email)
		q.Result = mock.JonSnow
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		Expect(c.Key).Equals(key)
		return nil
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")

	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestVerifySignInKeyHandler_RecentlyUsedKey_ShouldAllowReuse(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		Expect(q.Key).Equals(key)
		Expect(q.Kind).Equals(enum.EmailVerificationKindSignIn)

		verifiedAt := time.Now().Add(-4 * time.Minute)
		q.Result = &entity.EmailVerification{
			Key:        q.Key,
			Kind:       q.Kind,
			VerifiedAt: &verifiedAt,
			ExpiresAt:  time.Now().Add(1 * time.Hour),
			Email:      "jon.snow@got.com",
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		Expect(c.Key).Equals(key)
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		Expect(q.Email).Equals(mock.JonSnow.Email)
		q.Result = mock.JonSnow
		return nil
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

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

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignIn {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignInKeyHandler_PrivateTenant_SignInRequest_NonInviteNewUser(t *testing.T) {
	RegisterT(t)

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignIn {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusForbidden)
}

func TestVerifySignInKeyHandler_PrivateTenant_SignInRequest_RegisteredUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignIn {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	user := &entity.User{
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

	verified := false
	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		if c.Key == key {
			verified = true
		}
		return nil
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	Expect(verified).IsTrue()
}

func TestVerifySignInKeyHandler_PrivateTenant_InviteRequest_ExistingUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	user := &entity.User{
		Name:   "Hot Pie",
		Email:  "hot.pie@got.com",
		Tenant: mock.DemoTenant,
	}

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindUserInvitation {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	verified := false
	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		if c.Key == key {
			verified = true
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == user.Email {
			q.Result = user
			return nil
		}
		return app.ErrNotFound
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/invite/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindUserInvitation))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	Expect(verified).IsTrue()
}

func TestVerifySignInKeyHandler_PrivateTenant_InviteRequest_NewUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindUserInvitation {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/invite/verify?k=" + key).
		Execute(handlers.VerifySignInKey(enum.EmailVerificationKindUserInvitation))

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignUpKeyHandler_PendingTenant(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantPending

	var newUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		newUser = c.User
		return nil
	})

	activated := false
	bus.AddHandler(func(ctx context.Context, c *cmd.ActivateTenant) error {
		activated = true
		return nil
	})

	key := "1234567890"
	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		if q.Key == key && q.Kind == enum.EmailVerificationKindSignUp {
			expiresAt := time.Now().Add(5 * time.Minute)
			q.Result = &entity.EmailVerification{
				Key:       q.Key,
				Kind:      q.Kind,
				ExpiresAt: expiresAt,
				Name:      "Hot Pie",
				Email:     "hot.pie@got.com",
			}
			return nil
		}
		return app.ErrNotFound
	})

	verified := false
	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		if c.Key == key {
			verified = true
		}
		return nil
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signup/verify?k=" + key).
		Execute(handlers.VerifySignUpKey())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	Expect(newUser.Name).Equals("Hot Pie")
	Expect(newUser.Email).Equals("hot.pie@got.com")
	Expect(activated).IsTrue()
	Expect(verified).IsTrue()
}

func TestCompleteSignInProfileHandler_UnknownKey(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), `{ }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCompleteSignInProfileHandler_ExistingUser_CorrectKey(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	key := "1234567890"

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		Expect(q.Email).Equals(mock.JonSnow.Email)
		q.Result = mock.JonSnow
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		Expect(q.Key).Equals(key)
		Expect(q.Kind).Equals(enum.EmailVerificationKindSignIn)

		expiresAt := time.Now().Add(5 * time.Minute)
		q.Result = &entity.EmailVerification{
			Key:       q.Key,
			Kind:      q.Kind,
			ExpiresAt: expiresAt,
			Email:     mock.JonSnow.Email,
		}
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), fmt.Sprintf(`
		{
			"name": "Hot Pie",
			"kind": %d,
			"key": "%s"
		}`, enum.EmailVerificationKindSignIn, key))

	Expect(code).Equals(http.StatusBadRequest)
	ExpectHandler(&query.GetVerificationByKey{}).CalledOnce()
	ExpectHandler(&query.GetUserByEmail{}).CalledOnce()
}

func TestCompleteSignInProfileHandler_CorrectKey(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	key := "1234567890"

	var newUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		Expect(c.User.Name).Equals("Hot Pie")
		Expect(c.User.Email).Equals("hot.pie@got.com")
		newUser = c.User
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByKey) error {
		Expect(q.Key).Equals(key)
		Expect(q.Kind).Equals(enum.EmailVerificationKindSignIn)

		expiresAt := time.Now().Add(5 * time.Minute)
		q.Result = &entity.EmailVerification{
			Key:       q.Key,
			Kind:      q.Kind,
			ExpiresAt: expiresAt,
			Email:     "hot.pie@got.com",
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		Expect(c.Key).Equals(key)
		return nil
	})

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin/complete").
		ExecutePost(handlers.CompleteSignInProfile(), fmt.Sprintf(`
		{
			"name": "Hot Pie",
			"kind": %d,
			"key": "%s"
		}`, enum.EmailVerificationKindSignIn, key))

	Expect(code).Equals(http.StatusOK)
	ExpectHandler(&cmd.SetKeyAsVerified{}).CalledOnce()
	ExpectHandler(&query.GetVerificationByKey{}).CalledOnce()
	ExpectHandler(&query.GetUserByEmail{}).CalledOnce()
	ExpectFiderAuthCookie(response, newUser)
}

func TestSignInPageHandler_AuthenticatedUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
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

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin").
		Execute(handlers.SignInPage())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
}

func TestSignInPageHandler_PrivateTenant_UnauthenticatedUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/signin").
		Execute(handlers.SignInPage())

	Expect(code).Equals(http.StatusOK)
}

func TestVerifySignInCodeHandler_InvalidCode(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "jon.snow@got.com", "code": "999999" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestVerifySignInCodeHandler_ExpiredCode(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		q.Result = &entity.EmailVerification{
			Email:     "jon.snow@got.com",
			Key:       "123456",
			CreatedAt: time.Now().Add(-20 * time.Minute),
			ExpiresAt: time.Now().Add(-5 * time.Minute),
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "jon.snow@got.com", "code": "123456" }`)

	Expect(code).Equals(http.StatusGone)
}

func TestVerifySignInCodeHandler_CorrectCode_ExistingUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		q.Result = &entity.EmailVerification{
			Email:     "jon.snow@got.com",
			Key:       "123456",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		q.Result = mock.JonSnow
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		return nil
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "jon.snow@got.com", "code": "123456" }`)

	Expect(code).Equals(http.StatusOK)
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestVerifySignInCodeHandler_CorrectCode_NewUser_WithoutName(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		q.Result = &entity.EmailVerification{
			Email:     "new.user@got.com",
			Key:       "123456",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Name:      "", // No name stored (legacy flow)
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "new.user@got.com", "code": "123456" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(response.Body.String()).ContainsSubstring(`"showProfileCompletion":true`)
}

func TestVerifySignInCodeHandler_CorrectCode_NewUser_WithName(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		q.Result = &entity.EmailVerification{
			Email:     "new.user@got.com",
			Key:       "123456",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Name:      "New User", // Name stored (new flow)
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	userRegistered := false
	var registeredUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		userRegistered = true
		registeredUser = c.User
		registeredUser.ID = 999
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetKeyAsVerified) error {
		return nil
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "new.user@got.com", "code": "123456" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(userRegistered).IsTrue()
	Expect(registeredUser.Name).Equals("New User")
	Expect(registeredUser.Email).Equals("new.user@got.com")
	ExpectFiderAuthCookie(response, registeredUser)
}

func TestVerifySignInCodeHandler_CorrectCode_NewUser_PrivateTenant(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	bus.AddHandler(func(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
		q.Result = &entity.EmailVerification{
			Email:     "new.user@got.com",
			Key:       "123456",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.VerifySignInCode(), `{ "email": "new.user@got.com", "code": "123456" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestResendSignInCodeHandler_ValidEmail(t *testing.T) {
	RegisterT(t)

	var saveKeyCmd *cmd.SaveVerificationKey
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		saveKeyCmd = c
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.ResendSignInCode(), `{ "email": "jon.snow@got.com" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(saveKeyCmd.Key).HasLen(6)
	Expect(saveKeyCmd.Request.GetKind()).Equals(enum.EmailVerificationKindSignIn)
	Expect(saveKeyCmd.Request.GetEmail()).Equals("jon.snow@got.com")
}

func TestResendSignInCodeHandler_InvalidEmail(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		ExecutePost(handlers.ResendSignInCode(), `{ "email": "invalid" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func ExpectFiderAuthCookie(response *httptest.ResponseRecorder, expected *entity.User) {
	cookies := response.Header()["Set-Cookie"]
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

func ExpectFiderToken(token string, expected *entity.User) {
	user, err := jwt.DecodeFiderClaims(token)
	Expect(err).IsNil()
	Expect(user.UserID).Equals(expected.ID)
	Expect(user.UserName).Equals(expected.Name)
	Expect(user.UserEmail).Equals(expected.Email)
}
