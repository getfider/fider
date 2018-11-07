package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSignUpHandler_MultiTenant(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/signup").
		Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusOK)
}

func TestSignUpHandler_MultiTenant_WrongURL(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signup").
		Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://login.test.fider.io/signup")
}

func TestSignUpHandler_SingleTenant_NoTenants(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewSingleTenantServer()
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusOK)
}

func TestSignUpHandler_SingleTenant_WithTenants(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewSingleTenantServer()
	services.Tenants.Add("Game of Thrones", "got", models.TenantActive)
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

func TestCheckAvailabilityHandler_InvalidSubdomain(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.AddParam("subdomain", "").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.String("message")).IsNotEmpty()
}

func TestCheckAvailabilityHandler_UnavailableSubdomain(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.AddParam("subdomain", "demo").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.String("message")).IsNotEmpty()
}

func TestCheckAvailabilityHandler_ValidSubdomain(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.AddParam("subdomain", "mycompany").ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).Equals(http.StatusOK)
	Expect(response.Contains("message")).IsFalse()
}

func TestCreateTenantHandler_EmptyInput(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.ExecutePost(handlers.CreateTenant(), `{ }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCreateTenantHandler_WithSocialAccount(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
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

	tenant, err := services.Tenants.GetByDomain("mycompany")

	Expect(code).Equals(http.StatusOK)

	Expect(err).IsNil()
	Expect(tenant.Name).Equals("My Company")
	Expect(tenant.Subdomain).Equals("mycompany")
	Expect(tenant.Status).Equals(models.TenantActive)

	services.SetCurrentTenant(tenant)
	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Role).Equals(models.RoleAdministrator)

	cookie := web.ParseCookie(response.Header().Get("Set-Cookie"))
	Expect(cookie.Name).Equals(web.CookieSignUpAuthName)
	ExpectFiderToken(cookie.Value, user)
	Expect(cookie.Domain).Equals("test.fider.io")
	Expect(cookie.HttpOnly).IsTrue()
	Expect(cookie.Path).Equals("/")
	Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(5*time.Minute), 5*time.Second)
}

func TestCreateTenantHandler_SingleHost_WithSocialAccount(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewSingleTenantServer()
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

	tenant, err := services.Tenants.First()

	Expect(code).Equals(http.StatusOK)

	Expect(err).IsNil()
	Expect(tenant.Name).Equals("My Company")
	Expect(tenant.Subdomain).Equals("default")
	Expect(tenant.Status).Equals(models.TenantActive)

	services.SetCurrentTenant(tenant)
	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Role).Equals(models.RoleAdministrator)

	ExpectFiderAuthCookie(response, &models.User{
		ID:    1,
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
	})
}

func TestCreateTenantHandler_WithEmailAndName(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
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

	tenant, err := services.Tenants.GetByDomain("mycompany")

	Expect(code).Equals(http.StatusOK)

	Expect(err).IsNil()
	Expect(tenant.Name).Equals("My Company")
	Expect(tenant.Subdomain).Equals("mycompany")
	Expect(tenant.Status).Equals(models.TenantPending)

	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(user).IsNil()
}
