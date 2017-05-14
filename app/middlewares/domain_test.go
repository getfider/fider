package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/models"
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
	tenants.Add(&models.Tenant{Name: "The Orange Inc.", Subdomain: "orange"})
	tenants.Add(&models.Tenant{Name: "Demonstration", Subdomain: "demo"})

	for _, testCase := range testCases {
		for _, host := range testCase.hosts {

			server := mock.NewServer()
			server.Context.Request().Host = host

			server.Use(middlewares.MultiTenant(tenants))
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

	server.Use(middlewares.MultiTenant(tenants))
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

	server.Use(middlewares.MultiTenant(tenants))
	status, _ := server.Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestSingleTenant(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	server := mock.NewServer()
	server.Context.Request().Host = "somedomain.com"

	server.Use(middlewares.SingleTenant(tenants))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tenant, err := tenants.First()

	Expect(status).To(Equal(http.StatusOK))
	Expect(err).To(BeNil())
	Expect(tenant.Name).To(Equal("Default"))
	Expect(tenant.Subdomain).To(Equal("default"))
	Expect(server.Context.Tenant()).To(Equal(tenant))
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
