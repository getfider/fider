package handlers_test

import (
	"context"
	"github.com/getfider/fider/app/pkg/env"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func compareGeneratorResponse(input string, fileName string) {
	//os.WriteFile(env.Path(fileName), []byte(input), 0744)
	bytes, err := os.ReadFile(env.Path(fileName))
	Expect(err).IsNil()
	Expect(input).Equals(string(bytes))
}

func TestGlobalFeedHandler(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "First Post",
		Slug:        "first-post",
		Description: "Description of first post",
		CreatedAt:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		User:        &entity.User{ID: 1, Name: "Jon Snow"},
		VotesCount:  5,
	}

	post2 := &entity.Post{
		ID:          2,
		Number:      2,
		Title:       "Second Post",
		Slug:        "second-post",
		Description: "Description of second post",
		CreatedAt:   time.Date(2023, 1, 3, 10, 0, 0, 0, time.UTC),
		User:        &entity.User{ID: 2, Name: "Arya Stark"},
		VotesCount:  2,
	}

	bus.AddHandler(func(ctx context.Context, q *query.SearchPosts) error {
		q.Result = []*entity.Post{post1, post2}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAssignedTags) error {
		q.Result = []*entity.Tag{}
		return nil
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.GlobalFeed())

	Expect(code).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Type")).Equals("application/atom+xml")

	responseBody := response.Body.String()
	compareGeneratorResponse(responseBody, "app/handlers/testdata/global_feed.atom")
}

func TestCommentFeedHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "The Post",
		Slug:        "the-post",
		Description: "Description of the post",
		CreatedAt:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		User:        &entity.User{ID: 1, Name: "Jon Snow"},
		Status:      enum.PostOpen,
	}

	comment1 := &entity.Comment{
		ID:        1,
		Content:   "This is comment 1",
		CreatedAt: time.Date(2023, 1, 2, 10, 0, 0, 0, time.UTC),
		User:      &entity.User{ID: 2, Name: "Arya Stark"},
	}

	comment2 := &entity.Comment{
		ID:        2,
		Content:   "This is comment 2",
		CreatedAt: time.Date(2023, 1, 3, 10, 0, 0, 0, time.UTC),
		User:      &entity.User{ID: 3, Name: "Sansa Stark"},
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentsByPost) error {
		q.Result = []*entity.Comment{comment1, comment2}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAssignedTags) error {
		q.Result = []*entity.Tag{}
		return nil
	})

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("path", "1.atom").
		Execute(handlers.CommentFeed())

	Expect(code).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Type")).Equals("application/atom+xml")

	responseBody := response.Body.String()
	compareGeneratorResponse(responseBody, "app/handlers/testdata/comment_feed.atom")
}

func TestCommentFeedHandler_NotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		// Return not found error for any post number
		return app.ErrNotFound
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentsByPost) error {
		q.Result = []*entity.Comment{}
		return nil
	})

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("path", "999.atom").
		Execute(handlers.CommentFeed())

	Expect(code).Equals(http.StatusNotFound)
}

func TestFeedDisabled(t *testing.T) {
	RegisterT(t)

	disabledTenant := &entity.Tenant{
		ID:            998,
		Name:          "Feeds Disabled",
		Subdomain:     "feeds-disabled",
		Status:        enum.TenantActive,
		IsFeedEnabled: false,
	}

	server := mock.NewServer()
	code, _ := server.
		OnTenant(disabledTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.GlobalFeed())

	Expect(code).Equals(http.StatusNotFound)

	code, _ = server.
		OnTenant(disabledTenant).
		AsUser(mock.JonSnow).
		AddParam("path", "1.atom").
		Execute(handlers.CommentFeed())

	Expect(code).Equals(http.StatusNotFound)
}

func TestPrivacyEnabled(t *testing.T) {
	RegisterT(t)

	disabledTenant := &entity.Tenant{
		ID:            999,
		Name:          "Privacy Enabled",
		Subdomain:     "privacy-enabled",
		Status:        enum.TenantActive,
		IsPrivate:     true,
		IsFeedEnabled: true, // feed should be disabled regardless
	}

	server := mock.NewServer()
	code, _ := server.
		OnTenant(disabledTenant).
		AsUser(mock.JonSnow).
		Execute(handlers.GlobalFeed())

	Expect(code).Equals(http.StatusNotFound)

	code, _ = server.
		OnTenant(disabledTenant).
		AsUser(mock.JonSnow).
		AddParam("path", "1.atom").
		Execute(handlers.CommentFeed())

	Expect(code).Equals(http.StatusNotFound)
}
