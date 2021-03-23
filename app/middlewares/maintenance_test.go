package middlewares_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestMaintenance_Disabled(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.Maintenance())
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, _ := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
}

func TestMaintenance_Enabled(t *testing.T) {
	RegisterT(t)

	defer func() {
		env.Config.Maintenance.Enabled = false
	}()

	server := mock.NewServer()
	env.Config.Maintenance.Enabled = true
	server.Use(middlewares.ClientCache(30 * time.Hour))
	server.Use(middlewares.Maintenance())
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, response := server.Execute(handler)

	Expect(status).Equals(http.StatusServiceUnavailable)
	Expect(response.Header().Get("Cache-Control")).Equals("no-cache, no-store")
	Expect(response.Header().Get("Retry-After")).Equals("3600")
}
