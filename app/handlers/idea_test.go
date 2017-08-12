package handlers_test

import (
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/onsi/gomega"
)

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := setupServer()
	code, _ := server.OnTenant(demoTenant).AsUser(jonSnow).Execute(handlers.Index())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My Idea", "My Idea Description", jonSnow.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		WithParam("number", idea.Number).
		Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	server, _ := setupServer()
	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		WithParam("number", "99").
		Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(404))
}

func TestPostIdeaHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "My newest idea :)" }`)

	idea, err := services.Ideas.GetByID(1)
	Expect(code).To(Equal(200))
	Expect(err).To(BeNil())
	Expect(idea.Title).To(Equal("My newest idea :)"))
	Expect(idea.TotalSupporters).To(Equal(1))
}

func TestPostIdeaHandler_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		ExecutePost(handlers.PostIdea(), `{ "title": "" }`)

	_, err := services.Ideas.GetByID(1)
	Expect(code).To(Equal(400))
	Expect(err).NotTo(BeNil())
}

func TestUpdateIdeaHandler_IdeaOwner(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", aryaStark.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(aryaStark).
		WithParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	idea, _ = services.Ideas.GetByNumber(idea.Number)
	Expect(code).To(Equal(200))
	Expect(idea.Title).To(Equal("the new title"))
	Expect(idea.Description).To(Equal("new description"))
}

func TestUpdateIdeaHandler_TenantStaff(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", aryaStark.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		WithParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	idea, _ = services.Ideas.GetByNumber(idea.Number)
	Expect(code).To(Equal(200))
	Expect(idea.Title).To(Equal("the new title"))
	Expect(idea.Description).To(Equal("new description"))
}

func TestUpdateIdeaHandler_NonAuthorized(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", jonSnow.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(aryaStark).
		WithParam("number", idea.Number).
		ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).To(Equal(401))
}

func TestPostCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", jonSnow.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		WithParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("My First Idea", "With a description", jonSnow.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(jonSnow).
		WithParam("number", idea.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "" }`)

	Expect(code).To(Equal(400))
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	first, _ := services.Ideas.Add("The Idea #1", "The Description #1", jonSnow.ID)
	second, _ := services.Ideas.Add("The Idea #2", "The Description #2", jonSnow.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(aryaStark).
		WithParam("number", second.Number).
		Execute(handlers.AddSupporter())

	first, _ = services.Ideas.GetByNumber(1)
	second, _ = services.Ideas.GetByNumber(2)

	Expect(code).To(Equal(200))
	Expect(first.TotalSupporters).To(Equal(0))
	Expect(second.TotalSupporters).To(Equal(1))
}

func TestAddSupporterHandler_InvalidIdea(t *testing.T) {
	RegisterTestingT(t)

	server, _ := setupServer()
	code, _ := server.
		OnTenant(demoTenant).
		AsUser(aryaStark).
		WithParam("number", 999).
		Execute(handlers.AddSupporter())

	Expect(code).To(Equal(404))
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := setupServer()
	idea, _ := services.Ideas.Add("The Idea #1", "The Description #1", jonSnow.ID)
	services.Ideas.AddSupporter(idea.Number, jonSnow.ID)
	services.Ideas.AddSupporter(idea.Number, aryaStark.ID)

	code, _ := server.
		OnTenant(demoTenant).
		AsUser(aryaStark).
		WithParam("number", idea.ID).
		Execute(handlers.RemoveSupporter())

	idea, _ = services.Ideas.GetByNumber(idea.Number)

	Expect(code).To(Equal(200))
	Expect(idea.TotalSupporters).To(Equal(1))
}
