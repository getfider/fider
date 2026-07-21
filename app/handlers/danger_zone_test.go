package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

// ownerIs registers a GetTenantOwner handler that always returns the given user.
func ownerIs(owner *entity.User) {
	bus.AddHandler(func(ctx context.Context, q *query.GetTenantOwner) error {
		q.Result = owner
		return nil
	})
}

func TestRequestTenantDeletion_SingleHostBlocked(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow)

	server := mock.NewSingleTenantServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.RequestTenantDeletion(), `{"subdomain":"demo"}`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestRequestTenantDeletion_NonOwnerForbidden(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow) // owner is Jon, but Arya is making the request

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		ExecutePost(handlers.RequestTenantDeletion(), `{"subdomain":"demo"}`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestRequestTenantDeletion_WrongSubdomain(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow)

	scheduled := false
	bus.AddHandler(func(ctx context.Context, c *cmd.ScheduleTenantDeletion) error {
		scheduled = true
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.RequestTenantDeletion(), `{"subdomain":"not-demo"}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(scheduled).IsFalse()
}

func TestRequestTenantDeletion_OwnerSchedules(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow)

	var scheduled *cmd.ScheduleTenantDeletion
	bus.AddHandler(func(ctx context.Context, c *cmd.ScheduleTenantDeletion) error {
		scheduled = c
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.RequestTenantDeletion(), `{"subdomain":"demo"}`)

	Expect(code).Equals(http.StatusOK)
	Expect(scheduled).IsNotNil()
	Expect(scheduled.TenantID).Equals(mock.DemoTenant.ID)
	Expect(scheduled.RequestedByUserID).Equals(mock.JonSnow.ID)
	Expect(len(scheduled.CancelKey)).Equals(64)
	Expect(scheduled.ScheduledAt.After(time.Now())).IsTrue()
}

func TestCancelTenantDeletionByOwner_NonOwnerForbidden(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow)

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		ExecutePost(handlers.CancelTenantDeletionByOwner(), `{}`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestCancelTenantDeletionByOwner_OwnerCancels(t *testing.T) {
	RegisterT(t)
	ownerIs(mock.JonSnow)

	var cancelled *cmd.CancelTenantDeletion
	bus.AddHandler(func(ctx context.Context, c *cmd.CancelTenantDeletion) error {
		cancelled = c
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.CancelTenantDeletionByOwner(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(cancelled).IsNotNil()
	Expect(cancelled.TenantID).Equals(mock.DemoTenant.ID)
}
