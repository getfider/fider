package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	expected string
	urls     []string
}{
	{
		"Avengers",
		[]string{
			"http://avengers.test.fider.io",
			"http://avengers.test.fider.io:3000",
		},
	},
	{
		"Demonstration",
		[]string{
			"http://demo.test.fider.io",
			"http://demo.test.fider.io:1231",
			"http://demo.test.fider.io:80",
		},
	},
}

func TestMultiTenant(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range testCases {
		for _, url := range testCase.urls {

			server, _ := mock.NewServer()
			server.Use(middlewares.MultiTenant())

			status, response := server.WithURL(url).Execute(func(c web.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})

			Expect(status).To(Equal(http.StatusOK))
			Expect(response.Body.String()).To(Equal(testCase.expected))
		}
	}
}

func TestMultiTenant_SubSubDomain(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://demo.demo.test.fider.io").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestSingleTenant_NoTenants_RedirectToSignUp(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())

	status, response := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.HeaderMap.Get("Location")).To(Equal("/signup"))
}

func TestSingleTenant_WithTenants_ShouldSetFirstToContext(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())
	services.Tenants.Add("MyCompany", "mycompany", models.TenantActive)

	status, response := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Body.String()).Should(Equal("MyCompany"))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()

	server.Use(middlewares.HostChecker("login.test.fider.io"))
	status, _ := server.WithURL("http://login.test.fider.io").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.HostChecker("login.test.fider.io"))
	status, _ := server.WithURL("http://avengers.test.fider.io").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusBadRequest))
}

func TestOnlyActiveTenants_Active(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()

	server.Use(middlewares.OnlyActiveTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestOnlyActiveTenants_Inactive(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.Status = models.TenantInactive

	server.Use(middlewares.OnlyActiveTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestCheckTenantPrivacy_Private_Unauthenticated(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, response := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.HeaderMap.Get("Location")).To(Equal("/signin"))
}

func TestCheckTenantPrivacy_Private_Authenticated(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).To(Equal(http.StatusOK))
}

func TestCheckTenantPrivacy_NotPrivate_Unauthenticated(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = false

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).To(Equal(http.StatusOK))
}
