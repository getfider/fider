package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	expected string
	hosts    []string
}{
	{
		"The Orange Inc.",
		[]string{
			"orange.test.fider.io",
			"orange.test.fider.io:3000",
		},
	},
	{
		"Demonstration",
		[]string{
			"demo.test.fider.io",
			"demo.test.fider.io:1231",
			"demo.test.fider.io:80",
		},
	},
}

func TestMultiTenant(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenants.Add("The Orange Inc.", "orange")
	tenants.Add("Demonstration", "demo")

	for _, testCase := range testCases {
		for _, host := range testCase.hosts {

			server := mock.NewServer()
			server.Context.Request().Host = host
			server.Context.SetServices(&app.Services{
				Tenants: tenants,
			})

			server.Use(middlewares.MultiTenant())
			status, response := server.Execute(func(c web.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})

			Expect(status).To(Equal(http.StatusOK))
			Expect(response.Body.String()).To(Equal(testCase.expected))
		}
	}
}

func TestMultiTenant_AuthDomain(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	server := mock.NewServer()
	server.Context.Request().Host = "login.test.fider.io"
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})

	server.Use(middlewares.MultiTenant())
	status, _ := server.Execute(func(c web.Context) error {
		if c.Tenant() == nil {
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusInternalServerError)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	server := mock.NewServer()
	server.Context.Request().Host = "somedomain.com"
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})

	server.Use(middlewares.MultiTenant())
	status, _ := server.Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestSingleTenant_NoTenants_RedirectToSignUp(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	server := mock.NewServer()
	server.Context.Request().Host = "somedomain.com"
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})

	server.Use(middlewares.SingleTenant())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.HeaderMap.Get("Location")).To(Equal("/signup"))
}

func TestSingleTenant_WithTenants_ShouldSetToContext(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenant, _ := tenants.Add("Some Tenant", "mytenant")

	server := mock.NewServer()
	server.Context.Request().Host = "somedomain.com"
	server.Context.SetServices(&app.Services{
		Tenants: tenants,
	})

	server.Use(middlewares.SingleTenant())
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(server.Context.Tenant()).Should(Equal(tenant))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().Host = "login.test.fider.io"

	server.Use(middlewares.HostChecker("login.test.fider.io"))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.fider.io"
	server.Use(middlewares.HostChecker("login.test.fider.io"))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusBadRequest))
}
