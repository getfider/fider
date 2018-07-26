package middlewares_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jwt"
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

func TestOnlyActiveTenants_Active(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	server.Use(middlewares.OnlyActiveTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestOnlyActiveTenants_Inactive(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	mock.DemoTenant.Status = models.TenantInactive

	server.Use(middlewares.OnlyActiveTenants())
	status, _ := server.OnTenant(mock.DemoTenant).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusNotFound)
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

func TestUser_NoCookie(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.User())
	status, _ := server.Execute(func(c web.Context) error {
		if c.IsAuthenticated() {
			return c.NoContent(http.StatusOK)
		} else {
			return c.NoContent(http.StatusNoContent)
		}
	})

	Expect(status).Equals(http.StatusNoContent)
}

func TestUser_WithCookie(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
	Expect(response.HeaderMap["Set-Cookie"]).HasLen(0)
}

func TestUser_WithCookie_InvalidUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   999,
		UserName: "Unknown",
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusNoContent)
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=;")
}

func TestUser_WithCookie_DifferentTenant(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusNoContent)
}

func TestUser_WithSignUpCookie(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieSignUpAuthName, token).
		Execute(func(c web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
	cookies := response.HeaderMap["Set-Cookie"]
	Expect(cookies).HasLen(2)

	cookie := web.ParseCookie(cookies[0])
	Expect(cookie.Name).Equals(web.CookieSignUpAuthName)
	Expect(cookie.Value).Equals("")
	Expect(cookie.Domain).Equals("test.fider.io")
	Expect(cookie.HttpOnly).IsTrue()
	Expect(cookie.Path).Equals("/")
	Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(-100*time.Hour), 5*time.Second)

	cookie = web.ParseCookie(cookies[1])
	Expect(cookie.Name).Equals(web.CookieAuthName)
	Expect(cookie.Value).Equals(token)
	Expect(cookie.Domain).Equals("")
	Expect(cookie.HttpOnly).IsTrue()
	Expect(cookie.Path).Equals("/")
	Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(365*24*time.Hour), 5*time.Second)
}
