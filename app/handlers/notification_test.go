package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestTotalUnreadNotificationsHandler(t *testing.T) {
	RegisterT(t)
	bus.AddHandler(func(ctx context.Context, q *query.CountUnreadNotifications) error {
		q.Result = 5
		return nil
	})

	server := mock.NewServer()

	code, query := server.
		ExecuteAsJSON(handlers.TotalUnreadNotifications())

	Expect(code).Equals(http.StatusOK)
	Expect(query.Int32("total")).Equals(5)
}

func TestNotificationsHandler(t *testing.T) {
	RegisterT(t)
	bus.AddHandler(func(ctx context.Context, q *query.GetActiveNotifications) error {
		return nil
	})

	server := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.Notifications())

	Expect(code).Equals(http.StatusOK)
}

func TestReadNotificationHandler(t *testing.T) {
	RegisterT(t)

	not1 := &entity.Notification{ID: 1, Link: "/abc"}
	not2 := &entity.Notification{ID: 2, Link: "/def"}

	bus.AddHandler(func(ctx context.Context, q *query.GetNotificationByID) error {
		if q.ID == 1 {
			q.Result = not1
		} else if q.ID == 2 {
			q.Result = not2
		} else {
			q.Result = nil
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.MarkNotificationAsRead) error {
		if c.ID == 1 {
			not1.Read = true
		} else if c.ID == 2 {
			not2.Read = true
		}
		return nil
	})

	server := mock.NewServer()

	code, resp := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com").
		AsUser(mock.JonSnow).
		AddParam("id", 2).
		Execute(handlers.ReadNotification())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(resp.Header().Get("Location")).Equals("http://example.com/def")
	Expect(not1.Read).IsFalse()
	Expect(not2.Read).IsTrue()
}

func TestReadAllNotificationsHandler(t *testing.T) {
	RegisterT(t)

	called := false
	bus.AddHandler(func(ctx context.Context, c *cmd.MarkAllNotificationsAsRead) error {
		called = true
		return nil
	})

	server := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.ReadAllNotifications())

	Expect(code).Equals(http.StatusOK)
	Expect(called).IsTrue()
}
