package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestTotalUnreadNotificationsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.JonSnow)
	services.Notifications.Insert(mock.JonSnow, "Title", "http://example.com", 1, mock.AryaStark.ID)

	code, query := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.TotalUnreadNotifications())

	Expect(code).To(Equal(http.StatusOK))
	Expect(query.Int32("total")).To(Equal(1))
}

func TestNotificationsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.JonSnow)
	services.Notifications.Insert(mock.JonSnow, "Title", "http://example.com", 1, mock.AryaStark.ID)
	services.Notifications.Insert(mock.JonSnow, "Title 2", "http://example.com", 1, mock.AryaStark.ID)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.Notifications())

	Expect(code).To(Equal(http.StatusOK))
}

func TestReadNotificationHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.JonSnow)
	not1, _ := services.Notifications.Insert(mock.JonSnow, "Title", "/abc", 1, mock.AryaStark.ID)
	not2, _ := services.Notifications.Insert(mock.JonSnow, "Title 2", "/def", 1, mock.AryaStark.ID)

	code, resp := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://example.com").
		AsUser(mock.JonSnow).
		AddParam("id", not2.ID).
		Execute(handlers.ReadNotification())

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(resp.Header().Get("Location")).To(Equal("http://example.com/def"))
	Expect(not1.Read).To(BeFalse())
	Expect(not2.Read).To(BeTrue())
}

func TestReadAllNotificationsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentUser(mock.JonSnow)
	not1, _ := services.Notifications.Insert(mock.JonSnow, "Title", "/abc", 1, mock.AryaStark.ID)
	not2, _ := services.Notifications.Insert(mock.JonSnow, "Title 2", "/def", 1, mock.AryaStark.ID)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.ReadAllNotifications())

	Expect(code).To(Equal(http.StatusOK))
	Expect(not1.Read).To(BeTrue())
	Expect(not2.Read).To(BeTrue())
}
