package handlers_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/models"
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
	tenants.Add(&models.Tenant{Name: "Game of Thrones", Subdomain: "got"})
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, _ := server.Execute(handlers.SignUp())

	Expect(code).To(Equal(307))
}

func TestCheckAvailabilityHandler_InvalidSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewSingleTenantServer()
	tenants := &inmemory.TenantStorage{}
	tenants.Add(&models.Tenant{Name: "Game of Thrones", Subdomain: "got"})
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	Expect(response.String("message")).NotTo(BeNil())
}

func TestCheckAvailabilityHandler_UnavailableSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewSingleTenantServer()
	server.Context.SetParamNames("subdomain")
	server.Context.SetParamValues("got")
	tenants := &inmemory.TenantStorage{}
	tenants.Add(&models.Tenant{Name: "Game of Thrones", Subdomain: "got"})
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	Expect(response.String("message")).NotTo(BeNil())
}

func TestCheckAvailabilityHandler_ValidSubdomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewSingleTenantServer()
	server.Context.SetParamNames("subdomain")
	server.Context.SetParamValues("mycompany")
	tenants := &inmemory.TenantStorage{}
	tenants.Add(&models.Tenant{Name: "Game of Thrones", Subdomain: "got"})
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})
	code, response := server.ExecuteAsJSON(handlers.CheckAvailability())

	Expect(code).To(Equal(200))
	_, err := response.String("message")
	Expect(err).NotTo(BeNil())
}
