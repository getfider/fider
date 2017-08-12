package handlers_test

import (
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/onsi/gomega"
)

func TestUpdateSettingsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		ExecutePost(
			handlers.UpdateSettings(),
			`{ "title": "GoT", "invitation": "Join us!", "welcomeMessage": "Welcome to GoT Feedback Forum" }`,
		)

	tenant, _ := services.Tenants.GetByDomain("demo")
	Expect(code).To(Equal(200))
	Expect(tenant.Name).To(Equal("GoT"))
	Expect(tenant.Invitation).To(Equal("Join us!"))
	Expect(tenant.WelcomeMessage).To(Equal("Welcome to GoT Feedback Forum"))
}
