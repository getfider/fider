package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestIndexHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).Execute(handlers.Index())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My Post", "My Post Description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", "99").
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusNotFound)
}
