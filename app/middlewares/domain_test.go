package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock2"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	expected string
	hosts    []string
}{
	{
		"The Orange Inc.",
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"Demonstration",
		[]string{
			"demo.test.canhearyou.com",
			"demo.test.canhearyou.com:1231",
			"demo.test.canhearyou.com:80",
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

			server := mock2.NewServer()
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
	server := mock2.NewServer()
	server.Context.Request().Host = "login.test.canhearyou.com"

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
	server := mock2.NewServer()
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
	server := mock2.NewServer()
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

	server := mock2.NewServer()
	server.Context.Request().Host = "login.test.canhearyou.com"

	server.Use(middlewares.HostChecker("login.test.canhearyou.com"))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	server := mock2.NewServer()
	server.Context.Request().Host = "orange.test.canhearyou.com"
	server.Use(middlewares.HostChecker("login.test.canhearyou.com"))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusBadRequest))
}
