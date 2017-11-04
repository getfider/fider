package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	code, _ := server.OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).Execute(handlers.Index())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	idea, _ := services.Ideas.Add("My Idea", "My Idea Description", mock.JonSnow.ID)

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
	idea, _ := services.Ideas.Add("My First Idea", "With a description", mock.AryaStark.ID)

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
	idea, _ := services.Ideas.Add("My First Idea", "With a description", mock.JonSnow.ID)

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
	idea, _ := services.Ideas.Add("My First Idea", "With a description", mock.JonSnow.ID)

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

func TestPostCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", mock.JonSnow.ID)

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
	idea, _ := services.Ideas.Add("My First Idea", "With a description", mock.JonSnow.ID)

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
	first, _ := services.Ideas.Add("The Idea #1", "The Description #1", mock.JonSnow.ID)
	second, _ := services.Ideas.Add("The Idea #2", "The Description #2", mock.JonSnow.ID)

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
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1", mock.JonSnow.ID)
	services.Ideas.AddSupporter(idea.Number, mock.JonSnow.ID)
	services.Ideas.AddSupporter(idea.Number, mock.AryaStark.ID)

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
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1", mock.AryaStark.ID)
	services.Ideas.AddSupporter(idea.Number, mock.AryaStark.ID)

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
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1", mock.AryaStark.ID)
	services.Ideas.AddSupporter(idea.Number, mock.AryaStark.ID)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", idea.ID).
		ExecutePost(handlers.SetResponse(), fmt.Sprintf(`{ "status": %d, "text": "Done!" }`, models.IdeaCompleted))

	idea, _ = services.Ideas.GetByNumber(idea.Number)

	Expect(code).To(Equal(http.StatusForbidden))
	Expect(idea.Status).To(Equal(models.IdeaNew))
}
