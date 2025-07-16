package middlewares_test

import (
	"context"
	"net/http"
	"os"
	"testing"

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

func TestBlockFreemiumBillingAccess_WhenFreemiumIsDisabled(t *testing.T) {
	RegisterT(t)

	// Save original value to restore later
	original := os.Getenv("PADDLE_FREEMIUM")
	defer func() {
		os.Setenv("PADDLE_FREEMIUM", original)
		env.Reload()
	}()

	// Disable freemium
	os.Setenv("PADDLE_FREEMIUM", "false")
	env.Reload()
	Expect(env.IsFreemium()).IsFalse()

	// Mock the GetBillingState handler to return a free forever tenant
	bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
		q.Result = &entity.BillingState{
			Status: enum.BillingFreeForever,
		}
		return nil
	})

	server := mock.NewServer()
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(middlewares.BlockFreemiumBillingAccess()(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		}))

	// Should allow access since freemium is disabled
	Expect(status).Equals(http.StatusOK)
}

func TestBlockFreemiumBillingAccess_WhenFreemiumIsEnabled_AndTenantIsFree(t *testing.T) {
	RegisterT(t)

	// Save original value to restore later
	original := os.Getenv("PADDLE_FREEMIUM")
	defer func() {
		os.Setenv("PADDLE_FREEMIUM", original)
		env.Reload()
	}()

	// Enable freemium
	os.Setenv("PADDLE_FREEMIUM", "true")
	env.Reload()
	Expect(env.IsFreemium()).IsTrue()

	// Mock the GetBillingState handler to return a free forever tenant
	bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
		q.Result = &entity.BillingState{
			Status: enum.BillingFreeForever,
		}
		return nil
	})

	server := mock.NewServer()
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(middlewares.BlockFreemiumBillingAccess()(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		}))

	// Should block access since freemium is enabled and tenant is on free plan
	Expect(status).Equals(http.StatusForbidden)
}

func TestBlockFreemiumBillingAccess_WhenFreemiumIsEnabled_AndTenantIsPaid(t *testing.T) {
	RegisterT(t)

	// Save original value to restore later
	original := os.Getenv("PADDLE_FREEMIUM")
	defer func() {
		os.Setenv("PADDLE_FREEMIUM", original)
		env.Reload()
	}()

	// Enable freemium
	os.Setenv("PADDLE_FREEMIUM", "true")
	env.Reload()
	Expect(env.IsFreemium()).IsTrue()

	// Mock the GetBillingState handler to return a paid tenant
	bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
		q.Result = &entity.BillingState{
			Status: enum.BillingActive,
		}
		return nil
	})

	server := mock.NewServer()
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(middlewares.BlockFreemiumBillingAccess()(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		}))

	// Should allow access since tenant is on paid plan
	Expect(status).Equals(http.StatusOK)
}

func TestBlockFreemiumBillingAccess_WhenFreemiumIsEnabled_AndTenantIsOnTrial(t *testing.T) {
	RegisterT(t)

	// Save original value to restore later
	original := os.Getenv("PADDLE_FREEMIUM")
	defer func() {
		os.Setenv("PADDLE_FREEMIUM", original)
		env.Reload()
	}()

	// Enable freemium
	os.Setenv("PADDLE_FREEMIUM", "true")
	env.Reload()
	Expect(env.IsFreemium()).IsTrue()

	// Mock the GetBillingState handler to return a trial tenant
	bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
		q.Result = &entity.BillingState{
			Status: enum.BillingTrial,
		}
		return nil
	})

	server := mock.NewServer()
	status, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(middlewares.BlockFreemiumBillingAccess()(func(c *web.Context) error {
			return c.NoContent(http.StatusOK)
		}))

	// Should allow access since tenant is on trial
	Expect(status).Equals(http.StatusOK)
}
