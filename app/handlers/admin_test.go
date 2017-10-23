package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/handlers"
	. "github.com/onsi/gomega"
)

func TestUpdateSettingsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.UpdateSettings(),
			`{ "title": "GoT", "invitation": "Join us!", "welcomeMessage": "Welcome to GoT Feedback Forum" }`,
		)

	tenant, _ := services.Tenants.GetByDomain("demo")
	Expect(code).To(Equal(http.StatusOK))
	Expect(tenant.Name).To(Equal("GoT"))
	Expect(tenant.Invitation).To(Equal("Join us!"))
	Expect(tenant.WelcomeMessage).To(Equal("Welcome to GoT Feedback Forum"))
}

func TestUManageMembersHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		Execute(
			handlers.ManageMembers(),
		)

	Expect(code).To(Equal(http.StatusOK))
}
