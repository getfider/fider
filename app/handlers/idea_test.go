package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).Execute(handlers.Index())

	Expect(code).To(Equal(http.StatusOK))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My Idea", "My Idea Description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", "99").
		Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(404))
}

func TestPostIdeaHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "My newest idea :)" }`)

	idea, err := services.Ideas.GetByID(1)
	Expect(code).To(Equal(200))
	Expect(err).To(BeNil())
	Expect(idea.Title).To(Equal("My newest idea :)"))
	Expect(idea.TotalSupporters).To(Equal(1))
}

func TestPostIdeaHandler_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "" }`)

	_, err := services.Ideas.GetByID(1)
	Expect(code).To(Equal(400))
	Expect(err).NotTo(BeNil())
}

func TestUpdateIdeaHandler_TenantStaff(t *testing.T) {
	RegisterTestingT(t)

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
	Expect(code).To(Equal(200))
	Expect(idea.Title).To(Equal("the new title"))
	Expect(idea.Description).To(Equal("new description"))
}

func TestUpdateIdeaHandler_NonAuthorized(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).To(Equal(http.StatusForbidden))
}

func TestUpdateIdeaHandler_InvalidTitle(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "", "description": "" }`)

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestUpdateIdeaHandler_InvalidIdea(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 999).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "This is a good title!", "description": "And description too..." }`)

	Expect(code).To(Equal(http.StatusNotFound))
}

func TestUpdateIdeaHandler_DuplicateTitle(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestPostCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("My First Idea", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "" }`)

	Expect(code).To(Equal(400))
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(200))
	Expect(first.TotalSupporters).To(Equal(0))
	Expect(second.TotalSupporters).To(Equal(1))
}

func TestAddSupporterHandler_InvalidIdea(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 999).
		Execute(handlers.AddSupporter())

	Expect(code).To(Equal(404))
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(200))
	Expect(idea.TotalSupporters).To(Equal(1))
}

func TestSetResponseHandler(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusOK))
	Expect(idea.Status).To(Equal(models.IdeaCompleted))
	Expect(idea.Response.Text).To(Equal("Done!"))
	Expect(idea.Response.User.ID).To(Equal(mock.JonSnow.ID))
}

func TestSetResponseHandler_Unauthorized(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusForbidden))
	Expect(idea.Status).To(Equal(models.IdeaOpen))
}

func TestSetResponseHandler_Duplicate(t *testing.T) {
	RegisterTestingT(t)

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
	Expect(code).To(Equal(http.StatusOK))

	idea1, _ = services.Ideas.GetByNumber(idea1.Number)
	Expect(idea1.Status).To(Equal(models.IdeaDuplicate))

	idea2, _ = services.Ideas.GetByNumber(idea2.Number)
	Expect(idea2.Status).To(Equal(models.IdeaOpen))
}

func TestSetResponseHandler_Duplicate_NotFound(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestSetResponseHandler_Duplicate_Itself(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusBadRequest))
}

func TestAddCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "My first comment" }`)

	Expect(code).To(Equal(http.StatusOK))
	comments, _ := services.Ideas.GetCommentsByIdea(idea)
	Expect(comments).To(HaveLen(1))
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusOK))
	comment, _ := services.Ideas.GetCommentByID(commentId)
	Expect(comment.Content).To(Equal("My first comment has been edited"))
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(code).To(Equal(http.StatusForbidden))
	comment, _ := services.Ideas.GetCommentByID(commentId)
	Expect(comment.Content).To(Equal("My first comment"))
}

func TestDeleteIdeaHandler_Authorized(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", idea.Number).
		ExecutePost(handlers.DeleteIdea(), `{ }`)

	Expect(code).To(Equal(http.StatusOK))
	idea, err := services.Ideas.GetByNumber(idea.Number)
	Expect(idea).To(BeNil())
	Expect(errors.Cause(err)).To(Equal(app.ErrNotFound))
}
