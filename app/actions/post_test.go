package actions_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/actions"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestCreateNewPost_InvalidPostTitles(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		if q.Slug == "my-great-post" {
			q.Result = &entity.Post{Slug: q.Slug}
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
		action := &actions.CreateNewPost{Title: title}
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
		action := &actions.CreateNewPost{Title: title}
		result := action.Validate(context.Background(), nil)
		ExpectSuccess(result)
	}
}

func TestSetResponse_InvalidStatus(t *testing.T) {
	RegisterT(t)

	action := &actions.SetResponse{
		Status: enum.PostDeleted,
		Text:   "Spam!",
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "status")
}

func TestDeletePost_WhenIsBeingReferenced(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{ID: 1, Number: 1, Title: "Post 1"}
	post2 := &entity.Post{ID: 2, Number: 2, Title: "Post 2"}

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

	action := &actions.DeletePost{}
	action.Number = post1.Number
	ExpectSuccess(action.Validate(context.Background(), nil))

	action.Number = post2.Number
	ExpectFailed(action.Validate(context.Background(), nil))
}

func TestDeleteComment(t *testing.T) {
	RegisterT(t)

	author := &entity.User{ID: 1, Role: enum.RoleVisitor}
	notAuthor := &entity.User{ID: 2, Role: enum.RoleVisitor}
	administrator := &entity.User{ID: 3, Role: enum.RoleAdministrator}
	comment := &entity.Comment{
		ID:      1,
		User:    author,
		Content: "Comment #1",
	}

	post := &entity.Post{ID: 1, Number: 1}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		q.Result = &entity.Post{ID: 2, Number: q.Number}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == comment.ID && q.PostID == post.ID {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	action := &actions.DeleteComment{
		PostNumber: post.Number,
		CommentID:  comment.ID,
	}

	authorized := action.IsAuthorized(context.Background(), notAuthor)
	Expect(authorized).IsFalse()

	authorized = action.IsAuthorized(context.Background(), author)
	Expect(authorized).IsTrue()

	authorized = action.IsAuthorized(context.Background(), administrator)
	Expect(authorized).IsTrue()

	// Comment requested through a different post
	action.PostNumber = 2
	Expect(action.IsAuthorized(context.Background(), author)).IsFalse()
	Expect(action.IsAuthorized(context.Background(), administrator)).IsFalse()
}

func TestAddNewComment_TooLongContent(t *testing.T) {
	RegisterT(t)

	action := &actions.AddNewComment{Content: strings.Repeat("a", 4001)}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "content")
}

func TestAddNewComment_AtMaxLength(t *testing.T) {
	RegisterT(t)

	action := &actions.AddNewComment{Content: strings.Repeat("a", 4000)}
	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
}

func TestEditComment_TooLongContent(t *testing.T) {
	RegisterT(t)

	action := &actions.EditComment{Content: strings.Repeat("a", 4001)}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "content")
}

func TestEditComment_AtMaxLength(t *testing.T) {
	RegisterT(t)

	action := &actions.EditComment{Content: strings.Repeat("a", 4000)}
	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
}

func TestCreateNewPost_NullAttachments(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		return app.ErrNotFound
	})

	action := &actions.CreateNewPost{}
	err := json.Unmarshal([]byte(`{"title": "this is my new post", "attachments": [null, null]}`), action)
	Expect(err).IsNil()

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.Attachments).HasLen(0)
}

func TestAddNewComment_NullAttachments(t *testing.T) {
	RegisterT(t)

	action := &actions.AddNewComment{}
	err := json.Unmarshal([]byte(`{"content": "Nice idea", "attachments": [null]}`), action)
	Expect(err).IsNil()

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.Attachments).HasLen(0)
}

func TestEditComment_NullAttachments(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error {
		q.Result = []string{"attachments/a.png"}
		return nil
	})

	action := &actions.EditComment{}
	err := json.Unmarshal([]byte(`{"content": "Nice idea", "attachments": [null, {"bkey": "attachments/a.png", "remove": true}]}`), action)
	Expect(err).IsNil()

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.Attachments).HasLen(1)
	Expect(action.Attachments[0].BlobKey).Equals("attachments/a.png")
}

func executeEditComment(existing []string, attachments ...*dto.ImageUpload) *validate.Result {
	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error {
		q.Result = existing
		return nil
	})

	action := &actions.EditComment{
		Content:     "Nice idea",
		Attachments: attachments,
		Post:        &entity.Post{ID: 1},
		Comment:     &entity.Comment{ID: 1},
	}
	return action.Validate(context.Background(), nil)
}

func newUpload(content []byte) *dto.ImageUpload {
	return &dto.ImageUpload{Upload: &dto.ImageUploadData{FileName: "image.png", ContentType: "image/png", Content: content}}
}

func TestEditComment_ValidAttachment(t *testing.T) {
	RegisterT(t)

	result := executeEditComment([]string{}, newUpload(mock.UniformPNG(100, 100)))
	ExpectSuccess(result)
}

func TestEditComment_UnsupportedAttachment(t *testing.T) {
	RegisterT(t)

	result := executeEditComment([]string{}, newUpload([]byte("<html><script>alert(1)</script></html>")))
	ExpectFailed(result, "attachments")
}

func TestEditComment_DecompressionBombAttachment(t *testing.T) {
	RegisterT(t)

	result := executeEditComment([]string{}, newUpload(mock.UniformPNG(12000, 12000)))
	ExpectFailed(result, "attachments")
}

func TestEditComment_TooManyAttachments(t *testing.T) {
	RegisterT(t)

	// The comment already has 2 attachments, so a new one exceeds the limit...
	existing := []string{"attachments/a.png", "attachments/b.png"}
	result := executeEditComment(existing, newUpload(mock.UniformPNG(100, 100)))
	ExpectFailed(result, "attachments")

	// ...unless one of the existing ones is removed
	result = executeEditComment(existing,
		&dto.ImageUpload{BlobKey: "attachments/a.png", Remove: true},
		newUpload(mock.UniformPNG(100, 100)),
	)
	ExpectSuccess(result)
}
