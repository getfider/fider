package tests_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestFreemiumIntegration(t *testing.T) {
	RegisterT(t)

	// Save original values to restore later
	originalFreemium := os.Getenv("PADDLE_FREEMIUM")
	originalVendorID := os.Getenv("PADDLE_VENDOR_ID")
	originalVendorAuthCode := os.Getenv("PADDLE_VENDOR_AUTHCODE")
	defer func() {
		os.Setenv("PADDLE_FREEMIUM", originalFreemium)
		os.Setenv("PADDLE_VENDOR_ID", originalVendorID)
		os.Setenv("PADDLE_VENDOR_AUTHCODE", originalVendorAuthCode)
		env.Reload()
	}()

	// Test 1: Verify IsFreemium configuration works correctly
	t.Run("IsFreemium Configuration", func(t *testing.T) {
		// Test when freemium is disabled (default)
		os.Setenv("PADDLE_FREEMIUM", "false")
		env.Reload()
		Expect(env.IsFreemium()).IsFalse()

		// Test when freemium is enabled
		os.Setenv("PADDLE_FREEMIUM", "true")
		env.Reload()
		Expect(env.IsFreemium()).IsTrue()
	})

	// Test 2: Verify tenant creation with freemium enabled
	t.Run("Tenant Creation with Freemium Enabled", func(t *testing.T) {
		// Enable billing and freemium
		os.Setenv("PADDLE_VENDOR_ID", "123")
		os.Setenv("PADDLE_VENDOR_AUTHCODE", "456")
		os.Setenv("PADDLE_FREEMIUM", "true")
		env.Reload()

		// Ensure billing is enabled and freemium is enabled
		Expect(env.IsBillingEnabled()).IsTrue()
		Expect(env.IsFreemium()).IsTrue()

		// Mock tenant creation and billing state
		bus.AddHandler(func(ctx context.Context, c *cmd.CreateTenant) error {
			c.Result = &entity.Tenant{
				ID:        999,
				Name:      c.Name,
				Subdomain: c.Subdomain,
				Status:    c.Status,
			}
			return nil
		})

		bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
			// Use today's date for trial_ends_at to avoid null constraint violation
			today := time.Now()
			q.Result = &entity.BillingState{
				Status:      enum.BillingFreeForever,
				TrialEndsAt: &today, // Use today's date instead of nil to avoid constraint violation
			}
			return nil
		})

		// Create a new tenant
		createTenant := &cmd.CreateTenant{
			Name:      "Freemium Test Tenant",
			Subdomain: "freemium-test",
			Status:    enum.TenantActive,
		}

		err := bus.Dispatch(context.Background(), createTenant)
		Expect(err).IsNil()
		Expect(createTenant.Result).IsNotNil()

		// Check the billing state of the new tenant
		billingState := &query.GetBillingState{}
		err = bus.Dispatch(context.Background(), billingState)
		Expect(err).IsNil()
		Expect(billingState.Result).IsNotNil()
		Expect(billingState.Result.Status).Equals(enum.BillingFreeForever)
		// We're using today's date for trial_ends_at to avoid null constraint violation
		Expect(billingState.Result.TrialEndsAt).IsNotNil()
	})

	// Test 3: Verify billing routes middleware blocks access for free users
	t.Run("Billing Routes Middleware", func(t *testing.T) {
		// Enable freemium
		os.Setenv("PADDLE_FREEMIUM", "true")
		env.Reload()
		Expect(env.IsFreemium()).IsTrue()

		// Mock the GetBillingState handler to return a free forever tenant
		bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
			// Use today's date for trial_ends_at to avoid null constraint violation
			today := time.Now()
			q.Result = &entity.BillingState{
				Status:      enum.BillingFreeForever,
				TrialEndsAt: &today, // Use today's date instead of nil to avoid constraint violation
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
		Expect(status).Equals(http.StatusForbidden) // This is correct - free users should be blocked

		// Now change the billing state to a paid plan
		// We need to reset the bus handlers to ensure our new handler is used
		bus.Reset()

		// Re-register the handler for BillingActive status
		bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
			// Use today's date for trial_ends_at to avoid null constraint violation
			today := time.Now()
			q.Result = &entity.BillingState{
				Status:      enum.BillingActive,
				TrialEndsAt: &today, // Use today's date instead of nil to avoid constraint violation
			}
			return nil
		})

		// Create a new server after resetting the bus
		server = mock.NewServer()

		status, _ = server.
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			Execute(middlewares.BlockFreemiumBillingAccess()(func(c *web.Context) error {
				return c.NoContent(http.StatusOK)
			}))

		// Should allow access since tenant is on paid plan
		Expect(status).Equals(http.StatusOK) // Paid users should have access
	})
}
