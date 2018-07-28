package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
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

func TestPostCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	first, _ := services.Posts.Add("The Post #1", "The Description #1")
	second, _ := services.Posts.Add("The Post #2", "The Description #2")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", second.Number).
		Execute(handlers.AddSupporter())

	first, _ = services.Posts.GetByNumber(1)
	second, _ = services.Posts.GetByNumber(2)

	Expect(code).Equals(http.StatusOK)
	Expect(first.TotalSupporters).Equals(0)
	Expect(second.TotalSupporters).Equals(1)
}

func TestAddSupporterHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 999).
		Execute(handlers.AddSupporter())

	Expect(code).Equals(http.StatusNotFound)
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	services.Posts.AddSupporter(post, mock.JonSnow)
	services.Posts.AddSupporter(post, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.ID).
		Execute(handlers.RemoveSupporter())

	post, _ = services.Posts.GetByNumber(post.Number)

	Expect(code).Equals(http.StatusOK)
	Expect(post.TotalSupporters).Equals(1)
}

func TestAddCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "My first comment" }`)

	Expect(code).Equals(http.StatusOK)
	comments, _ := services.Posts.GetCommentsByPost(post)
	Expect(comments).HasLen(1)
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	commentId, _ := services.Posts.AddComment(post, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusOK)
	comment, _ := services.Posts.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment has been edited")
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	commentId, _ := services.Posts.AddComment(post, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusForbidden)
	comment, _ := services.Posts.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment")
}

func TestDeletePostHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.DeletePost(), `{ }`)

	Expect(code).Equals(http.StatusOK)
	post, err := services.Posts.GetByNumber(post.Number)
	Expect(post).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
}
