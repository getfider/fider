package apiv1_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app/handlers/apiv1"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestCreatePostHandler(t *testing.T) {
	RegisterT(t)

	var newPost *cmd.AddNewPost
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewPost) error {
		newPost = c
		c.Result = &entity.Post{
			ID:          1,
			Title:       c.Title,
			Description: c.Description,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newPost.Title).Equals("My newest post :)")
	Expect(newPost.Description).Equals("")
}

func TestCreatePostHandler_WithoutTitle(t *testing.T) {
	RegisterT(t)

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestGetPostHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecuteAsJSON(apiv1.GetPost())

	Expect(code).Equals(http.StatusOK)
	Expect(query.String("title")).Equals(post.Title)
	Expect(query.String("description")).Equals(post.Description)
}

func TestUpdatePostHandler_TenantStaff(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "With a description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	var updatePost *cmd.UpdatePost
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdatePost) error {
		updatePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(updatePost.Post).Equals(post)
	Expect(updatePost.Title).Equals("the new title")
	Expect(updatePost.Description).Equals("new description")
}

func TestUpdatePostHandler_NonAuthorized(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID: 5,
		Number: 5,
		Title: "My First Post",
		Description: "Such an amazing description",
		User: mock.JonSnow,
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestUpdatePostHandler_IsOwner_AfterGracePeriod(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID: 5,
		Number: 5,
		Title: "My First Post",
		Description: "Such an amazing description",
		User: mock.AryaStark,
		CreatedAt: time.Now().UTC().Add(-2 * time.Hour),
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}


func TestUpdatePostHandler_IsOwner_WithinGracePeriod(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID: 5,
		Number: 5,
		Title: "My First Post",
		Description: "Such an amazing description",
		User: mock.AryaStark,
		CreatedAt: time.Now().UTC(),
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.UploadImages) error { return nil })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.UpdatePost) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestUpdatePostHandler_InvalidTitle(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "", "description": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdatePostHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error { return app.ErrNotFound })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 999).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "This is a good title!", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestUpdatePostHandler_DuplicateTitle(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{ID: 1, Number: 1, Title: "My First Post", Slug: "my-first-post"}
	post2 := &entity.Post{ID: 2, Number: 2, Title: "My Second Post", Slug: "my-second-post"}
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

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		if q.Slug == post1.Slug {
			q.Result = post1
			return nil
		}
		if q.Slug == post2.Slug {
			q.Result = post2
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "My Second Post", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "My First Post", Slug: "my-first-post"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	var setResponse *cmd.SetPostResponse
	bus.AddHandler(func(ctx context.Context, c *cmd.SetPostResponse) error {
		setResponse = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.SetResponse(), fmt.Sprintf(`{ "status": "%s", "text": "Done!" }`, enum.PostCompleted.Name()))

	Expect(code).Equals(http.StatusOK)
	Expect(setResponse.Post).Equals(post)
	Expect(setResponse.Status).Equals(enum.PostCompleted)
	Expect(setResponse.Text).Equals("Done!")
}

func TestSetResponseHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 5).
		ExecutePost(apiv1.SetResponse(), fmt.Sprintf(`{ "status": "%s", "text": "Done!" }`, enum.PostCompleted.Name()))

	Expect(code).Equals(http.StatusForbidden)
}

func TestSetResponseHandler_Duplicate(t *testing.T) {
	RegisterT(t)

	var markAsDuplicate *cmd.MarkPostAsDuplicate
	bus.AddHandler(func(ctx context.Context, c *cmd.MarkPostAsDuplicate) error {
		markAsDuplicate = c
		return nil
	})

	post1 := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	post2 := &entity.Post{ID: 2, Number: 2, Title: "The Post #2", Description: "The Description #2"}
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

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": %d }`, enum.PostDuplicate.Name(), post2.Number)
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusOK)
	Expect(markAsDuplicate.Post).Equals(post1)
	Expect(markAsDuplicate.Original).Equals(post2)
}

func TestSetResponseHandler_Duplicate_NotFound(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post1.Number {
			q.Result = post1
			return nil
		}
		return app.ErrNotFound
	})

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": 9999 }`, enum.PostDuplicate.Name())
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler_Duplicate_Itself(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": %d }`, enum.PostDuplicate.Name(), post.Number)
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddVoteHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var addVote *cmd.AddVote
	bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error {
		addVote = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		Execute(apiv1.AddVote())

	Expect(code).Equals(http.StatusOK)
	Expect(addVote.Post).Equals(post)
	Expect(addVote.User).Equals(mock.AryaStark)
}

func TestAddVoteHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 999).
		Execute(apiv1.AddVote())

	Expect(code).Equals(http.StatusNotFound)
}

func TestRemoveVoteHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var removeVote *cmd.RemoveVote
	bus.AddHandler(func(ctx context.Context, c *cmd.RemoveVote) error {
		removeVote = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.ID).
		Execute(apiv1.RemoveVote())

	Expect(code).Equals(http.StatusOK)
	Expect(removeVote.Post).Equals(post)
	Expect(removeVote.User).Equals(mock.AryaStark)
}

func TestDeletePostHandler_Authorized(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.PostIsReferenced) error {
		q.Result = false
		return nil
	})

	var deletePost *cmd.SetPostResponse
	bus.AddHandler(func(ctx context.Context, c *cmd.SetPostResponse) error {
		deletePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.DeletePost(), `{ }`)

	Expect(code).Equals(http.StatusOK)
	Expect(deletePost.Post).Equals(post)
	Expect(deletePost.Status).Equals(enum.PostDeleted)
	Expect(deletePost.Text).Equals("")
}

func TestPostCommentHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var newComment *cmd.AddNewComment
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewComment) error {
		newComment = c
		c.Result = &entity.Comment{ID: 1, Content: c.Content}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newComment.Post).Equals(post)
	Expect(newComment.Content).Equals("This is a comment!")
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.PostComment(), `{ "content": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.AryaStark}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	var updateComment *cmd.UpdateComment
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateComment) error {
		updateComment = c
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		ExecutePost(apiv1.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(updateComment.Content).Equals("My first comment has been edited")
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.JonSnow}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		ExecutePost(apiv1.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestListCommentHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentsByPost) error {
		q.Result = []*entity.Comment{
			{ID: 1, Content: "First Comment"},
			{ID: 2, Content: "First Comment"},
		}
		return nil
	})

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecuteAsJSON(apiv1.ListComments())

	Expect(code).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(2)
}
