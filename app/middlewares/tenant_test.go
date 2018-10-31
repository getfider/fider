package middlewares_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
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
	RegisterT(t)

	for _, testCase := range testCases {
		for _, url := range testCase.urls {

			server, _ := mock.NewServer()
			server.Use(middlewares.MultiTenant())

			status, response := server.WithURL(url).Execute(func(c web.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})

			Expect(status).Equals(http.StatusOK)
			Expect(response.Body.String()).Equals(testCase.expected)
		}
	}
}

func TestMultiTenant_SubSubDomain(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://demo.demo.test.fider.io").Execute(func(c web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestMultiTenant_CanonicalHeader(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		input  string
		output string
		isAjax bool
	}{
		{
			"http://avengers.test.fider.io/",
			"<http://feedback.theavengers.com/>; rel=\"canonical\"",
			false,
		},
		{
			"http://avengers.test.fider.io/",
			"",
			true,
		},
		{
			"http://feedback.theavengers.com/",
			"",
			false,
		},
		{
			"http://avengers.test.fider.io/posts",
			"<http://feedback.theavengers.com/posts>; rel=\"canonical\"",
			false,
		},
		{
			"http://avengers.test.fider.io/posts?q=1",
			"<http://feedback.theavengers.com/posts?q=1>; rel=\"canonical\"",
			false,
		},
		{
			"http://demo.test.fider.io",
			"",
			false,
		},
	}

	for _, testCase := range testCases {
		server, _ := mock.NewServer()
		server.Use(middlewares.MultiTenant())

		if testCase.isAjax {
			server.AddHeader("Accept", "application/json")
		}
		status, response := server.
			WithURL(testCase.input).
			Execute(func(c web.Context) error {
				return c.Ok(web.Map{})
			})

		Expect(status).Equals(http.StatusOK)
		Expect(response.HeaderMap.Get("Link")).Equals(testCase.output)
	}

}

func TestSingleTenant_NoTenants(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestSingleTenant_WithTenants_ShouldSetFirstToContext(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())
	services.Tenants.Add("MyCompany", "mycompany", models.TenantActive)

	status, response := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("MyCompany")
}

func TestBlockPendingTenants_Active(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.Status = models.TenantActive

	server.Use(middlewares.BlockPendingTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestBlockPendingTenants_Pending(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.Status = models.TenantPending

	server.Use(middlewares.BlockPendingTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusTeapot)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckTenantPrivacy_Private_Unauthenticated(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, response := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.HeaderMap.Get("Location")).Equals("/signin")
}

func TestCheckTenantPrivacy_Private_Authenticated(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckTenantPrivacy_NotPrivate_Unauthenticated(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.IsPrivate = false

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestRequireTenant_MultiHostMode_NoTenants_404(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())
	server.Use(middlewares.RequireTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusNotFound)
}

func TestRequireTenant_MultiHostMode_ValidTenant(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewServer()
	server.Use(middlewares.MultiTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://avengers.test.fider.io").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Avengers")
}

func TestRequireTenant_SingleHostMode_NoTenants_RedirectToSignUp(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://somedomain.com").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.HeaderMap.Get("Location")).Equals("/signup")
}

func TestRequireTenant_SingleHostMode_ValidTenant(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewServer()
	server.Use(middlewares.SingleTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://demo.test.fider.io").Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Demonstration")
}
