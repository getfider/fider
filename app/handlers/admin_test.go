package handlers_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestUpdateSettingsHandler(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenants.Add("Demonstration", "demo")
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Tenants: tenants})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon", Role: models.RoleAdministrator})

	code, _ := server.ExecutePost(handlers.UpdateSettings(), `{ "title": "GoT", "invitation": "Join us!", "welcomeMessage": "Welcome to GoT Feedback Forum" }`)

	tenant, _ := tenants.GetByDomain("demo")
	Expect(code).To(Equal(200))
	Expect(tenant.Name).To(Equal("GoT"))
	Expect(tenant.Invitation).To(Equal("Join us!"))
	Expect(tenant.WelcomeMessage).To(Equal("Welcome to GoT Feedback Forum"))
}
