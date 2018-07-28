package apiv1_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers/apiv1"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestCreatePostHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)" }`)

	post, err := services.Posts.GetByID(1)
	Expect(code).Equals(http.StatusOK)
	Expect(err).IsNil()
	Expect(post.Title).Equals("My newest post :)")
	Expect(post.TotalSupporters).Equals(1)
}

func TestCreatePostHandler_WithoutTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "" }`)

	_, err := services.Posts.GetByID(1)
	Expect(code).Equals(http.StatusBadRequest)
	Expect(err).IsNotNil()
}
