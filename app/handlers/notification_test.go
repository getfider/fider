package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestTotalUnreadNotificationsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.AryaStark)
	services.Notifications.Insert(mock.JonSnow, "Title", "http://example.com", 1)

	code, query := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.TotalUnreadNotifications())

	Expect(code).Equals(http.StatusOK)
	Expect(query.Int32("total")).Equals(1)
}

func TestNotificationsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.AryaStark)
	services.Notifications.Insert(mock.JonSnow, "Title", "http://example.com", 1)
	services.Notifications.Insert(mock.JonSnow, "Title 2", "http://example.com", 1)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.Notifications())

	Expect(code).Equals(http.StatusOK)
}

func TestReadNotificationHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.JonSnow)
	services.SetCurrentUser(mock.AryaStark)
	not1, _ := services.Notifications.Insert(mock.JonSnow, "Title", "/abc", 1)
	not2, _ := services.Notifications.Insert(mock.JonSnow, "Title 2", "/def", 1)

	code, resp := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com").
		AsUser(mock.JonSnow).
		AddParam("id", not2.ID).
		Execute(handlers.ReadNotification())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(resp.Header().Get("Location")).Equals("http://example.com/def")
	Expect(not1.Read).IsFalse()
	Expect(not2.Read).IsTrue()
}

func TestReadAllNotificationsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.AryaStark)
	not1, _ := services.Notifications.Insert(mock.JonSnow, "Title", "/abc", 1)
	not2, _ := services.Notifications.Insert(mock.JonSnow, "Title 2", "/def", 1)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.ReadAllNotifications())

	Expect(code).Equals(http.StatusOK)
	Expect(not1.Read).IsTrue()
	Expect(not2.Read).IsTrue()
}
