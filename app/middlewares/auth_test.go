package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestIsAuthorized_WithAllowedRole(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.IsAuthorized(models.RoleAdministrator, models.RoleCollaborator))
	status, _ := server.AsUser(mock.JonSnow).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestIsAuthorized_WithForbiddenRole(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.IsAuthorized(models.RoleAdministrator, models.RoleCollaborator))
	status, _ := server.AsUser(mock.AryaStark).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusForbidden)
}

func TestIsAuthenticated_WithUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.IsAuthenticated())
	status, _ := server.AsUser(mock.AryaStark).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestIsAuthenticated_WithoutUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.IsAuthenticated())

	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusForbidden)
}

func TestCheckAuthTokenOrigin_WithoutClaims(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.CheckAuthTokenOrigin())

	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckAuthTokenOrigin_WithUIClaims_OnAPIResource(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.CheckAuthTokenOrigin())

	status, _ := server.
		WithURL("http://example.com/api/echo").
		WithClaims(&jwt.FiderClaims{Origin: jwt.FiderClaimsOriginUI}).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckAuthTokenOrigin_WithAPIClaims_OnAPIResource(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.CheckAuthTokenOrigin())

	status, _ := server.
		WithURL("http://example.com/api/echo").
		WithClaims(&jwt.FiderClaims{Origin: jwt.FiderClaimsOriginAPI}).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
}

func TestCheckAuthTokenOrigin_WithAPIClaims_OnUIResource(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.CheckAuthTokenOrigin())

	status, _ := server.
		WithURL("http://example.com/").
		WithClaims(&jwt.FiderClaims{Origin: jwt.FiderClaimsOriginAPI}).
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusForbidden)
}
