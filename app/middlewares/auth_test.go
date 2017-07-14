package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestIsAuthorized_WithAllowedRole(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetUser(&models.User{
		ID:   1,
		Role: models.RoleMember,
	})

	server.Use(middlewares.IsAuthorized(models.RoleAdministrator, models.RoleMember))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestIsAuthorized_WithForbiddenRole(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetUser(&models.User{
		ID:   1,
		Role: models.RoleVisitor,
	})

	server.Use(middlewares.IsAuthorized(models.RoleAdministrator, models.RoleMember))
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusForbidden))
}

func TestIsAuthenticated_WithUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetUser(&models.User{
		ID: 1,
	})

	server.Use(middlewares.IsAuthenticated())
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestIsAuthenticated_WithoutUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Use(middlewares.IsAuthenticated())

	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusForbidden))
}
