package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSignUpHandler_MultiTenant(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/signup").
		Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusOK)
}

func TestSignUpHandler_MultiTenant_WrongURL(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signup").
		Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://login.test.fider.io/signup")
}

func TestSignUpHandler_SingleTenant_NoTenants(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		return app.ErrNotFound
	})

	server := mock.NewSingleTenantServer()
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusOK)
}

func TestSignUpHandler_SingleTenant_WithTenants(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		q.Result = &entity.Tenant{ID: 2, Name: "MyCompany"}
		return nil
	})

	server := mock.NewSingleTenantServer()
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

func TestCheckAvailabilityHandler_InvalidSubdomain(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, response := server.AddParam("subdomain", "").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.String("message")).IsNotEmpty()
}

func TestCheckAvailabilityHandler_UnavailableSubdomain(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = false
		return nil
	})

	server := mock.NewServer()
	code, response := server.AddParam("subdomain", "demo").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.String("message")).IsNotEmpty()
}

func TestCheckAvailabilityHandler_ValidSubdomain(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = true
		return nil
	})

	server := mock.NewServer()
	code, response := server.AddParam("subdomain", "mycompany").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.Contains("message")).IsFalse()
}

func TestCreateTenantHandler_EmptyInput(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, _ := server.ExecutePost(handlers.CreateTenant(), `{ }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCreateTenantHandler_WithSocialAccount(t *testing.T) {
	RegisterT(t)

	var newUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		newUser = c.User
		return nil
	})

	var newTenant *cmd.CreateTenant
	bus.AddHandler(func(ctx context.Context, c *cmd.CreateTenant) error {
		newTenant = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = true
		return nil
	})

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.OAuthClaims{
		OAuthID:       "123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
		OAuthProvider: "facebook",
	})
	code, response := server.ExecutePost(
		handlers.CreateTenant(),
		fmt.Sprintf(`{
			"token": "%s",
			"tenantName": "My Company",
			"subdomain": "mycompany",
			"legalAgreement": true
		}`, token),
	)

	Expect(code).Equals(http.StatusOK)

	Expect(newTenant.Name).Equals("My Company")
	Expect(newTenant.Subdomain).Equals("mycompany")
	Expect(newTenant.Status).Equals(enum.TenantActive)

	Expect(newUser.Name).Equals("Jon Snow")
	Expect(newUser.Email).Equals("jon.snow@got.com")
	Expect(newUser.Role).Equals(enum.RoleAdministrator)

	cookie := web.ParseCookie(response.Header().Get("Set-Cookie"))
	Expect(cookie.Name).Equals(web.CookieSignUpAuthName)
	ExpectFiderToken(cookie.Value, newUser)
	Expect(cookie.Domain).Equals("test.fider.io")
	Expect(cookie.HttpOnly).IsTrue()
	Expect(cookie.Path).Equals("/")
	Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(5*time.Minute), 5*time.Second)
}

func TestCreateTenantHandler_SingleHost_WithSocialAccount(t *testing.T) {
	RegisterT(t)

	var newUser *entity.User
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		c.User.ID = 1
		newUser = c.User
		return nil
	})

	var newTenant *cmd.CreateTenant
	bus.AddHandler(func(ctx context.Context, c *cmd.CreateTenant) error {
		newTenant = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = true
		return nil
	})

	server := mock.NewSingleTenantServer()
	token, _ := jwt.Encode(jwt.OAuthClaims{
		OAuthID:       "123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
		OAuthProvider: "facebook",
	})
	code, response := server.ExecutePost(
		handlers.CreateTenant(),
		fmt.Sprintf(`{
			"token": "%s",
			"tenantName": "My Company",
			"legalAgreement": true
		}`, token),
	)

	Expect(code).Equals(http.StatusOK)

	Expect(newTenant.Name).Equals("My Company")
	Expect(newTenant.Subdomain).Equals("default")
	Expect(newTenant.Status).Equals(enum.TenantActive)

	Expect(newUser.Name).Equals("Jon Snow")
	Expect(newUser.Email).Equals("jon.snow@got.com")
	Expect(newUser.Role).Equals(enum.RoleAdministrator)

	ExpectFiderAuthCookie(response, &entity.User{
		ID:    1,
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
	})
}

func TestCreateTenantHandler_WithEmailAndName(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		panic("Should not register any user")
	})

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = true
		return nil
	})

	var newTenant *cmd.CreateTenant
	bus.AddHandler(func(ctx context.Context, c *cmd.CreateTenant) error {
		newTenant = c
		c.Result = &entity.Tenant{ID: 1, Name: c.Name, Subdomain: c.Subdomain, Status: c.Status}
		return nil
	})

	var saveKeyCmd *cmd.SaveVerificationKey
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		saveKeyCmd = c
		return nil
	})

	server := mock.NewServer()
	code, response := server.ExecutePost(
		handlers.CreateTenant(),
		`{
			"name": "Jon Snow",
			"email": "jon.snow@got.com",
			"tenantName": "My Company",
			"subdomain": "mycompany",
			"legalAgreement": true
		}`,
	)

	Expect(code).Equals(http.StatusOK)
	Expect(response.Header().Get("Set-Cookie")).IsEmpty()

	Expect(code).Equals(http.StatusOK)

	Expect(newTenant.Name).Equals("My Company")
	Expect(newTenant.Subdomain).Equals("mycompany")
	Expect(newTenant.Status).Equals(enum.TenantPending)

	Expect(saveKeyCmd.Key).HasLen(64)
	Expect(saveKeyCmd.Request.GetKind()).Equals(enum.EmailVerificationKindSignUp)
	Expect(saveKeyCmd.Request.GetEmail()).Equals("jon.snow@got.com")
	Expect(saveKeyCmd.Request.GetName()).Equals("Jon Snow")
}
