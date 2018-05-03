package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"

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
	idea, _ := services.Ideas.Add("My Idea", "My Idea Description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		Execute(handlers.IdeaDetails())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", "99").
		Execute(handlers.IdeaDetails())

	Expect(code).Equals(http.StatusNotFound)
}

func TestPostIdeaHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "My newest idea :)" }`)

	idea, err := services.Ideas.GetByID(1)
	Expect(code).Equals(http.StatusOK)
	Expect(err).IsNil()
	Expect(idea.Title).Equals("My newest idea :)")
	Expect(idea.TotalSupporters).Equals(1)
}

func TestPostIdeaHandler_WithoutTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "" }`)

	_, err := services.Ideas.GetByID(1)
	Expect(code).Equals(http.StatusBadRequest)
	Expect(err).IsNotNil()
}

func TestUpdateIdeaHandler_TenantStaff(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	idea, _ = services.Ideas.GetByNumber(idea.Number)
	Expect(code).Equals(http.StatusOK)
	Expect(idea.Title).Equals("the new title")
	Expect(idea.Description).Equals("new description")
}

func TestUpdateIdeaHandler_NonAuthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestUpdateIdeaHandler_InvalidTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "", "description": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdateIdeaHandler_InvalidIdea(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 999).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "This is a good title!", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestUpdateIdeaHandler_DuplicateTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea1, _ := services.Ideas.Add("My First Idea", "With a description")
	services.Ideas.Add("My Second Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea1.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "My Second Idea", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestPostCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	first, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	second, _ := services.Ideas.Add("The Idea #2", "The Description #2")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", second.Number).
		Execute(handlers.AddSupporter())

	first, _ = services.Ideas.GetByNumber(1)
	second, _ = services.Ideas.GetByNumber(2)

	Expect(code).Equals(http.StatusOK)
	Expect(first.TotalSupporters).Equals(0)
	Expect(second.TotalSupporters).Equals(1)
}

func TestAddSupporterHandler_InvalidIdea(t *testing.T) {
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
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	services.Ideas.AddSupporter(idea, mock.JonSnow)
	services.Ideas.AddSupporter(idea, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.ID).
		Execute(handlers.RemoveSupporter())

	idea, _ = services.Ideas.GetByNumber(idea.Number)

	Expect(code).Equals(http.StatusOK)
	Expect(idea.TotalSupporters).Equals(1)
}

func TestSetResponseHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	services.Ideas.AddSupporter(idea, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.ID).
		ExecutePost(handlers.SetResponse(), fmt.Sprintf(`{ "status": %d, "text": "Done!" }`, models.IdeaCompleted))

	idea, _ = services.Ideas.GetByNumber(idea.Number)

	Expect(code).Equals(http.StatusOK)
	Expect(idea.Status).Equals(models.IdeaCompleted)
	Expect(idea.Response.Text).Equals("Done!")
	Expect(idea.Response.User.ID).Equals(mock.JonSnow.ID)
}

func TestSetResponseHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	services.Ideas.AddSupporter(idea, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.ID).
		ExecutePost(handlers.SetResponse(), fmt.Sprintf(`{ "status": %d, "text": "Done!" }`, models.IdeaCompleted))

	idea, _ = services.Ideas.GetByNumber(idea.Number)

	Expect(code).Equals(http.StatusForbidden)
	Expect(idea.Status).Equals(models.IdeaOpen)
}

func TestSetResponseHandler_Duplicate(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea1, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	idea2, _ := services.Ideas.Add("The Idea #2", "The Description #2")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": %d }`, models.IdeaDuplicate, idea2.Number)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea1.ID).
		ExecutePost(handlers.SetResponse(), body)
	Expect(code).Equals(http.StatusOK)

	idea1, _ = services.Ideas.GetByNumber(idea1.Number)
	Expect(idea1.Status).Equals(models.IdeaDuplicate)

	idea2, _ = services.Ideas.GetByNumber(idea2.Number)
	Expect(idea2.Status).Equals(models.IdeaOpen)
}

func TestSetResponseHandler_Duplicate_NotFound(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea1, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": 9999 }`, models.IdeaDuplicate)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea1.ID).
		ExecutePost(handlers.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler_Duplicate_Itself(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": %d }`, models.IdeaDuplicate, idea.Number)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.ID).
		ExecutePost(handlers.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "My first comment" }`)

	Expect(code).Equals(http.StatusOK)
	comments, _ := services.Ideas.GetCommentsByIdea(idea)
	Expect(comments).HasLen(1)
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	commentId, _ := services.Ideas.AddComment(idea, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusOK)
	comment, _ := services.Ideas.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment has been edited")
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")
	commentId, _ := services.Ideas.AddComment(idea, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusForbidden)
	comment, _ := services.Ideas.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment")
}

func TestDeleteIdeaHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.DeleteIdea(), `{ }`)

	Expect(code).Equals(http.StatusOK)
	idea, err := services.Ideas.GetByNumber(idea.Number)
	Expect(idea).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
}
