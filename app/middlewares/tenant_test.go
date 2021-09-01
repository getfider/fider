package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
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

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		if q.Domain == "avengers.test.fider.io" {
			q.Result = mock.AvengersTenant
			return nil
		} else if q.Domain == "demo.test.fider.io" {
			q.Result = mock.DemoTenant
			return nil
		}
		return app.ErrNotFound
	})

	for _, testCase := range testCases {
		for _, url := range testCase.urls {

			server := mock.NewServer()
			server.Use(middlewares.MultiTenant())

			status, response := server.WithURL(url).Execute(func(c *web.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})

			Expect(status).Equals(http.StatusOK)
			Expect(response.Body.String()).Equals(testCase.expected)
		}
	}
}

func TestMultiTenant_SubSubDomain(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		if q.Domain == "demo.test.fider.io" {
			q.Result = mock.DemoTenant
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://demo.demo.test.fider.io").Execute(func(c *web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.MultiTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c *web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestMultiTenant_DisabledTenant_ShouldNotSetInContext(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		if q.Domain == "avengers.test.fider.io" {
			q.Result = mock.AvengersTenant
			q.Result.Status = enum.TenantDisabled
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.MultiTenant())
	server.Use(middlewares.RequireTenant())

	status, _ := server.WithURL("http://avengers.test.fider.io").Execute(func(c *web.Context) error {
		return c.Ok(nil)
	})

	Expect(status).Equals(http.StatusNotFound)
}

func TestMultiTenant_CanonicalHeader(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		if q.Domain == "avengers.test.fider.io" {
			q.Result = mock.AvengersTenant
			return nil
		} else if q.Domain == "demo.test.fider.io" {
			q.Result = mock.DemoTenant
			return nil
		}
		return app.ErrNotFound
	})

	var testCases = []struct {
		input  string
		output string
		isAjax bool
	}{
		{
			"http://avengers.test.fider.io/",
			"http://feedback.theavengers.com/",
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
			"http://feedback.theavengers.com/posts",
			false,
		},
		{
			"http://avengers.test.fider.io/posts?q=1",
			"http://feedback.theavengers.com/posts?q=1",
			false,
		},
		{
			"http://demo.test.fider.io",
			"",
			false,
		},
	}

	for _, testCase := range testCases {
		server := mock.NewServer()
		server.Use(middlewares.MultiTenant())

		if testCase.isAjax {
			server.AddHeader("Accept", "application/json")
		}

		var canonicalURL string
		status, _ := server.
			WithURL(testCase.input).
			Execute(func(c *web.Context) error {
				canonicalURL, _ = c.Value("Canonical-URL").(string)
				return c.Ok(web.Map{})
			})

		Expect(status).Equals(http.StatusOK)
		Expect(canonicalURL).Equals(testCase.output)
	}

}

func TestSingleTenant_NoTenants(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		return app.ErrNotFound
	})

	server := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c *web.Context) error {
		if c.Tenant() == nil {
			return c.Ok(web.Map{})
		}
		return c.Failure(errors.New("should not have found tenant"))
	})

	Expect(status).Equals(http.StatusOK)
}

func TestSingleTenant_WithTenants_ShouldSetFirstToContext(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		q.Result = &entity.Tenant{Name: "MyCompany", Subdomain: "mycompany", Status: enum.TenantActive}
		return nil
	})

	server := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())

	status, response := server.WithURL("http://test.fider.io").Execute(func(c *web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("MyCompany")
}

func TestSingleTenant_HostMismatch(t *testing.T) {
	RegisterT(t)
	env.Config.HostDomain = "yoursite.com"

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		q.Result = &entity.Tenant{Name: "MyCompany", Status: enum.TenantActive}
		return nil
	})

	server := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())

	status, _ := server.WithURL("http://someothersite.com").Execute(func(c *web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusNotFound)
}

func TestBlockPendingTenants_Active(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantActive

	server.Use(middlewares.BlockPendingTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestBlockPendingTenants_Pending(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantPending

	server.Use(middlewares.BlockPendingTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusTeapot)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckTenantPrivacy_Private_Unauthenticated(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, response := server.OnTenant(mock.DemoTenant).Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/signin")
}

func TestCheckTenantPrivacy_Private_Authenticated(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = true

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		Execute(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckTenantPrivacy_NotPrivate_Unauthenticated(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.IsPrivate = false

	server.Use(middlewares.CheckTenantPrivacy())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		Execute(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestRequireTenant_MultiHostMode_NoTenants_404(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.MultiTenant())
	server.Use(middlewares.RequireTenant())

	status, _ := server.WithURL("http://somedomain.com").Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusNotFound)
}

func TestRequireTenant_MultiHostMode_ValidTenant(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTenantByDomain) error {
		if q.Domain == "avengers.test.fider.io" {
			q.Result = mock.AvengersTenant
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.MultiTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://avengers.test.fider.io").Execute(func(c *web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Avengers")
}

func TestRequireTenant_SingleHostMode_NoTenants_RedirectToSignUp(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		return app.ErrNotFound
	})

	server := mock.NewSingleTenantServer()
	server.Use(middlewares.SingleTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://somedomain.com").Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/signup")
}

func TestRequireTenant_SingleHostMode_ValidTenant(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetFirstTenant) error {
		q.Result = mock.DemoTenant
		return nil
	})

	server := mock.NewServer()
	server.Use(middlewares.SingleTenant())
	server.Use(middlewares.RequireTenant())

	status, response := server.WithURL("http://test.fider.io").Execute(func(c *web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Demonstration")
}

func TestBlockLockedTenants_ActiveTenant(t *testing.T) {
	RegisterT(t)
	server := mock.NewServer()
	server.Use(middlewares.BlockLockedTenants())

	status, response := server.
		WithURL("http://demo.test.fider.io").
		OnTenant(mock.DemoTenant).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.Tenant().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Demonstration")
}

func TestBlockLockedTenants_LockedTenant(t *testing.T) {
	RegisterT(t)
	server := mock.NewServer()
	server.Use(middlewares.BlockLockedTenants())
	mock.DemoTenant.Status = enum.TenantLocked

	status, response := server.
		WithURL("http://demo.test.fider.io").
		OnTenant(mock.DemoTenant).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.Tenant().Name)
		})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/signin")
}

func TestBlockLockedTenants_LockedTenant_APICall(t *testing.T) {
	RegisterT(t)
	server := mock.NewServer()
	server.Use(middlewares.BlockLockedTenants())
	mock.DemoTenant.Status = enum.TenantLocked

	status, _ := server.
		WithURL("http://demo.test.fider.io/api/v1/posts").
		OnTenant(mock.DemoTenant).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.Tenant().Name)
		})

	Expect(status).Equals(http.StatusLocked)
}
