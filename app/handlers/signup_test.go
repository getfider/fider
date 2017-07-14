package handlers_test

import (
	"testing"

	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestSignUpHandler_MultiTenant(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).To(Equal(200))
}

func TestSignUpHandler_SingleTenant_NoTenants(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewSingleTenantServer()
	tenants := &inmemory.TenantStorage{}
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).To(Equal(200))
}

func TestSignUpHandler_SingleTenant_WithTenants(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewSingleTenantServer()
	tenants := &inmemory.TenantStorage{}
	tenants.Add("Game of Thrones", "got")
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).To(Equal(307))
}

func TestCheckAvailabilityHandler_InvalidSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	tenants := &inmemory.TenantStorage{}
	tenants.Add("Game of Thrones", "got")
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	Expect(response.String("message")).NotTo(BeNil())
}

func TestCheckAvailabilityHandler_UnavailableSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetParams(web.Map{"subdomain": "got"})
	tenants := &inmemory.TenantStorage{}
	tenants.Add("Game of Thrones", "got")
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	Expect(response.String("message")).NotTo(BeNil())
}

func TestCheckAvailabilityHandler_ValidSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetParams(web.Map{"subdomain": "mycompany"})
	tenants := &inmemory.TenantStorage{}
	tenants.Add("Game of Thrones", "got")
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	_, err := response.String("message")
	Expect(err).NotTo(BeNil())
}

func TestCreateTenantHandler_EmptyInput(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(&app.Services{
		Tenants: &inmemory.TenantStorage{},
	})
	code, _ := server.ExecutePost(handlers.CreateTenant(), `{ }`)

	Expect(code).To(Equal(400))
}

func TestCreateTenantHandler_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(&app.Services{
		Tenants: &inmemory.TenantStorage{},
	})
	token, _ := jwt.Encode(&models.OAuthClaims{
		OAuthID:       "123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
		OAuthProvider: "facebook",
	})
	code, _ := server.ExecutePost(handlers.CreateTenant(), fmt.Sprintf(`{ "token": "%s" }`, token))

	Expect(code).To(Equal(400))
}

func TestCreateTenantHandler_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(&app.Services{
		Tenants: &inmemory.TenantStorage{},
	})
	token, _ := jwt.Encode(&models.OAuthClaims{
		OAuthID:       "123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
		OAuthProvider: "facebook",
	})
	code, _ := server.ExecutePost(handlers.CreateTenant(), fmt.Sprintf(`{ "token": "%s", "name": "My Company" }`, token))

	Expect(code).To(Equal(400))
}

func TestCreateTenantHandler_ValidInput(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetServices(&app.Services{
		Tenants: &inmemory.TenantStorage{},
		Users:   &inmemory.UserStorage{},
	})
	token, _ := jwt.Encode(&models.OAuthClaims{
		OAuthID:       "123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
		OAuthProvider: "facebook",
	})
	code, _ := server.ExecutePost(handlers.CreateTenant(), fmt.Sprintf(`{ "token": "%s", "name": "My Company", "subdomain": "mycompany" }`, token))

	Expect(code).To(Equal(200))
}
