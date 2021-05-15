package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestIndexHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.CountPostPerStatus) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAllTags) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.SearchPosts) error {
		return nil
	})

	server := mock.NewServer()
	code, _ := server.OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.Index())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{Number: 1, Title: "My Post Title", Slug: "my-post-title"}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentsByPost) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.ListPostVotes) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAllTags) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.UserSubscribedTo) error {
		return nil
	})

	server := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		AddParam("slug", post.Slug).
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler_RedirectOnDifferentSlu(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{Number: 1, Title: "My Post Title", Slug: "my-post-title"}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	server := mock.NewServer()

	code, response := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		AddParam("slug", "some-other-slug").
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/posts/1/my-post-title")
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		return app.ErrNotFound
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", "99").
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusNotFound)
}
