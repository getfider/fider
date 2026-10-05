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

func TestUser_WithCookie_IsNotAuthenticatedByAPIKey(t *testing.T) {
	RegisterT(t)

	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		q.Result = mock.JonSnow
		return nil
	})

	// The cookie takes precedence, so the Bearer header is ignored
	server := mock.NewServer()
	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddCookie(web.CookieAuthName, token).
		AddHeader("Authorization", "Bearer 1234567890").
		Execute(func(c *web.Context) error {
			Expect(c.IsAuthenticated()).IsTrue()
			Expect(c.IsAuthenticatedByAPIKey()).IsFalse()
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
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

	Expect(status).Equals(http.StatusUnauthorized)
}

func TestUser_LockedTenant_ShouldAllowSignIn(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	mock.DemoTenant.Status = enum.TenantLocked
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.AryaStark.ID,
		UserName: mock.AryaStark.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		Expect(q.UserID).Equals(mock.AryaStark.ID)
		q.Result = mock.AryaStark
		return nil
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
	Expect(response.Body.String()).Equals("Arya Stark")
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

// TestUser_SignUpCookie_CrossTenant_Denied is the core regression test for the cross-tenant
// signup-cookie takeover. Jon Snow is an Administrator of the demo tenant. His domain-wide
// __signup_auth cookie is presented while the avengers tenant is in context. Because the
// user lookup is now tenant-scoped, no user is installed, the signup token is NOT promoted
// into a host-only auth cookie, and an authenticated route returns 401.
func TestUser_SignUpCookie_CrossTenant_Denied(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID, // belongs to the demo tenant
		UserName: mock.JonSnow.Name,
	})

	// Tenant-aware handler: only resolves the user when the requested tenant matches.
	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID && q.TenantID == mock.JonSnow.Tenant.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	server.Use(middlewares.IsAuthenticated())
	status, response := server.
		OnTenant(mock.AvengersTenant). // victim tenant, different from Jon Snow's
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieSignUpAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	// An authenticated (administrator) route must reject the request.
	Expect(status).Equals(http.StatusUnauthorized)

	// The domain-wide signup token must never become a durable host-only auth cookie on
	// the victim tenant. Any auth cookie emitted here must be a deletion (empty value).
	for _, raw := range response.Header()["Set-Cookie"] {
		c := web.ParseCookie(raw)
		if c.Name == web.CookieAuthName {
			Expect(c.Value).Equals("")
		}
	}
}

// TestUser_SignUpCookie_SameTenant_Promotes verifies the positive path: the signup cookie
// succeeds on the tenant the user actually belongs to and is promoted to a host-only auth
// cookie for that host.
func TestUser_SignUpCookie_SameTenant_Promotes(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.JonSnow.ID && q.TenantID == mock.JonSnow.Tenant.ID {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})

	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant). // Jon Snow's own tenant
		AddCookie(web.CookieSignUpAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")

	cookies := response.Header()["Set-Cookie"]
	Expect(cookies).HasLen(2)

	// The signup cookie is consumed (deleted)...
	signup := web.ParseCookie(cookies[0])
	Expect(signup.Name).Equals(web.CookieSignUpAuthName)
	Expect(signup.Value).Equals("")

	// ...and promoted to a host-only auth cookie (no Domain) carrying the same token.
	auth := web.ParseCookie(cookies[1])
	Expect(auth.Name).Equals(web.CookieAuthName)
	Expect(auth.Value).Equals(token)
	Expect(auth.Domain).Equals("")
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
			Expect(c.IsAuthenticatedByAPIKey()).IsTrue()
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

func TestUser_ValidAPIKey_LenientBearerParsing(t *testing.T) {
	for _, header := range []string{"bearer 1234567890", "BEARER 1234567890", "Bearer\t1234567890", "Bearer1234567890", "Bearer 1234567890 \r\n"} {
		t.Run(header, func(t *testing.T) {
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
				AddHeader("Authorization", header).
				Execute(func(c *web.Context) error {
					return c.String(http.StatusOK, c.User().Name)
				})

			Expect(status).Equals(http.StatusOK)
			Expect(response.Body.String()).Equals("Jon Snow")
		})
	}
}

func TestUser_OtherAuthorizationSchemes_AreIgnored(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		q.Result = mock.JonSnow
		return nil
	})

	server := mock.NewServer()

	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com/api/v1").
		AddHeader("Authorization", "Basic Bearer 1234567890").
		Execute(func(c *web.Context) error {
			Expect(c.User()).IsNil()
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
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

// TestUser_SecurityStamp_Match verifies that a token whose security stamp
// matches the DB value grants access normally.
func TestUser_SecurityStamp_Match(t *testing.T) {
	RegisterT(t)

	userWithStamp := &entity.User{
		ID:            mock.JonSnow.ID,
		Name:          mock.JonSnow.Name,
		Email:         mock.JonSnow.Email,
		Tenant:        mock.DemoTenant,
		Status:        enum.UserActive,
		Role:          enum.RoleAdministrator,
		SecurityStamp: "stamp-abc123",
		Providers:     mock.JonSnow.Providers,
	}

	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:        userWithStamp.ID,
		UserName:      userWithStamp.Name,
		SecurityStamp: "stamp-abc123",
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == userWithStamp.ID {
			q.Result = userWithStamp
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

// TestUser_SecurityStamp_Mismatch verifies that a browser session is redirected to
// /signin when the stamp in the JWT no longer matches the DB stamp (e.g. after a
// role change or an OAuth provider allowed-roles update).
func TestUser_SecurityStamp_Mismatch(t *testing.T) {
	RegisterT(t)

	userWithStamp := &entity.User{
		ID:            mock.JonSnow.ID,
		Name:          mock.JonSnow.Name,
		Email:         mock.JonSnow.Email,
		Tenant:        mock.DemoTenant,
		Status:        enum.UserActive,
		Role:          enum.RoleAdministrator,
		SecurityStamp: "stamp-new", // DB has been updated
		Providers:     mock.JonSnow.Providers,
	}

	// Token carries the OLD stamp
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:        userWithStamp.ID,
		UserName:      userWithStamp.Name,
		SecurityStamp: "stamp-old",
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == userWithStamp.ID {
			q.Result = userWithStamp
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/settings").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	// Browser request: should redirect to /signin with the current path as return URL
	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).ContainsSubstring("/signin")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=;")
}

// TestUser_SecurityStamp_Mismatch_AJAX verifies that an AJAX request receives a 401
// JSON response (not a redirect) when the security stamp is stale.
func TestUser_SecurityStamp_Mismatch_AJAX(t *testing.T) {
	RegisterT(t)

	userWithStamp := &entity.User{
		ID:            mock.JonSnow.ID,
		Name:          mock.JonSnow.Name,
		Email:         mock.JonSnow.Email,
		Tenant:        mock.DemoTenant,
		Status:        enum.UserActive,
		Role:          enum.RoleAdministrator,
		SecurityStamp: "stamp-new",
		Providers:     mock.JonSnow.Providers,
	}

	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:        userWithStamp.ID,
		UserName:      userWithStamp.Name,
		SecurityStamp: "stamp-old",
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == userWithStamp.ID {
			q.Result = userWithStamp
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.User())
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	// AJAX request: should get a 401 JSON, not a redirect
	Expect(status).Equals(http.StatusUnauthorized)
}

// TestUser_SecurityStamp_EmptyInToken verifies backward compatibility:
// old tokens without a stamp embedded must still be accepted.
func TestUser_SecurityStamp_EmptyInToken(t *testing.T) {
	RegisterT(t)

	userWithStamp := &entity.User{
		ID:            mock.JonSnow.ID,
		Name:          mock.JonSnow.Name,
		Email:         mock.JonSnow.Email,
		Tenant:        mock.DemoTenant,
		Status:        enum.UserActive,
		Role:          enum.RoleAdministrator,
		SecurityStamp: "stamp-in-db",
		Providers:     mock.JonSnow.Providers,
	}

	// Old token: no SecurityStamp field
	token, _ := jwt.Encode(jwt.FiderClaims{
		UserID:   userWithStamp.ID,
		UserName: userWithStamp.Name,
		// SecurityStamp deliberately omitted (simulates pre-stamp tokens)
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == userWithStamp.ID {
			q.Result = userWithStamp
			return nil
		}
		return app.ErrNotFound
	})

	server := mock.NewServer()
	server.Use(middlewares.User())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c *web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	// Old tokens without a stamp must still work (backward compatible)
	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

