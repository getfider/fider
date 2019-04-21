package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestCreateNewPost_InvalidPostTitles(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		if q.Slug == "my-great-post" {
			q.Result = &models.Post{Slug: q.Slug}
			return nil
		}
		return app.ErrNotFound
	})

	for _, title := range []string{
		"me",
		"",
		"  ",
		"signup",
		"My great great great great great great great great great great great great great great great great great post.",
		"my GREAT post",
	} {
		action := &actions.CreateNewPost{Model: &models.NewPost{Title: title}}
		result := action.Validate(context.Background(), nil)
		ExpectFailed(result, "title")
	}
}

func TestCreateNewPost_ValidPostTitles(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		return app.ErrNotFound
	})

	for _, title := range []string{
		"this is my new post",
		"this post is very descriptive",
	} {
		action := &actions.CreateNewPost{Model: &models.NewPost{Title: title}}
		result := action.Validate(context.Background(), nil)
		ExpectSuccess(result)
	}
}

func TestSetResponse_InvalidStatus(t *testing.T) {
	RegisterT(t)

	action := &actions.SetResponse{Model: &models.SetResponse{
		Status: enum.PostDeleted,
		Text:   "Spam!",
	}}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "status")
}

func TestDeletePost_WhenIsBeingReferenced(t *testing.T) {
	RegisterT(t)

	post1 := &models.Post{ID: 1, Number: 1, Title: "Post 1"}
	post2 := &models.Post{ID: 2, Number: 2, Title: "Post 2"}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post1.Number {
			q.Result = post1
			return nil
		}

		if q.Number == post2.Number {
			q.Result = post2
			return nil
		}

		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.PostIsReferenced) error {
		q.Result = q.PostID == post2.ID
		return nil
	})

	action := &actions.DeletePost{Model: &models.DeletePost{}}
	action.Model.Number = post1.Number
	ExpectSuccess(action.Validate(context.Background(), nil))

	action.Model.Number = post2.Number
	ExpectFailed(action.Validate(context.Background(), nil))
}

func TestDeleteComment(t *testing.T) {
	RegisterT(t)

	author := &models.User{ID: 1, Role: enum.RoleVisitor}
	notAuthor := &models.User{ID: 2, Role: enum.RoleVisitor}
	administrator := &models.User{ID: 3, Role: enum.RoleAdministrator}
	comment := &models.Comment{
		ID:      1,
		User:    author,
		Content: "Comment #1",
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == comment.ID {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	action := &actions.DeleteComment{
		Model: &models.DeleteComment{
			CommentID: comment.ID,
		},
	}

	authorized := action.IsAuthorized(context.Background(), notAuthor)
	Expect(authorized).IsFalse()

	authorized = action.IsAuthorized(context.Background(), author)
	Expect(authorized).IsTrue()

	authorized = action.IsAuthorized(context.Background(), administrator)
	Expect(authorized).IsTrue()
}
