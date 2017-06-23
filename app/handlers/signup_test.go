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
