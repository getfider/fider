package middlewares_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestUser_NoCookie(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.User())
	status, _ := server.Execute(func(c *web.Context) error {
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

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
	Expect(response.Header()["Set-Cookie"]).HasLen(0)
}

func TestUser_Blocked(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.JonSnow.Status = enum.UserBlocked
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusForbidden)
}

func TestUser_LockedTenant_Administrator(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantLocked
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

func TestUser_LockedTenant_Visitor(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantLocked
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.AryaStark.ID,
		UserName: mock.AryaStark.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.AryaStark.ID {
			q.Result = mock.AryaStark
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusForbidden)
}

func TestUser_WithCookie_InvalidUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   999,
		UserName: "Unknown",
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
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

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusNoContent)
}

func TestUser_WithSignUpCookie(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieSignUpAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
	cookies := response.Header()["Set-Cookie"]
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

func TestUser_ValidAPIKey(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "1234567890" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 1234567890").
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

func TestUser_InvalidAPIKey(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, query := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer MY-KEY").
		ExecuteAsJSON(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusBadRequest)
	Expect(query.String("errors[0].message")).Equals("API Key is invalid")
}

func TestUser_ValidAPIKey_Visitor(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, query := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 12345").
		ExecuteAsJSON(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusBadRequest)
	Expect(query.String("errors[0].message")).Equals("API Key is invalid")
}

func TestUser_Impersonation_Collaborator(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "12345" {
			q.Result = &entity.User{
				Name:   "The Collaborator",
				Role:   enum.RoleCollaborator,
				Status: enum.UserActive,
				Tenant: mock.DemoTenant,
			}
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, query := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 12345").
		AddHeader("X-Fider-UserID", strconv.Itoa(mock.JonSnow.ID)).
		ExecuteAsJSON(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusBadRequest)
	Expect(query.String("errors[0].message")).Equals("Only Administrators are allowed to impersonate another user")
}

func TestUser_Impersonation_InvalidUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "1234567890" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, query := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 1234567890").
		AddHeader("X-Fider-UserID", "ABC").
		ExecuteAsJSON(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusBadRequest)
	Expect(query.String("errors[0].message")).Equals("User not found for given impersonate UserID 'ABC'")
}

func TestUser_Impersonation_UserNotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "1234567890" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, query := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 1234567890").
		AddHeader("X-Fider-UserID", "999").
		ExecuteAsJSON(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusBadRequest)
	Expect(query.String("errors[0].message")).Equals("User not found for given impersonate UserID '999'")
}

func TestUser_Impersonation_ValidUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.AryaStark.ID {
			q.Result = mock.AryaStark
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "1234567890" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Bearer 1234567890").
		AddHeader("X-Fider-UserID", strconv.Itoa(mock.AryaStark.ID)).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Arya Stark")
}
