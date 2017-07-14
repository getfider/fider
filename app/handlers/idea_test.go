package handlers_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})

	code, _ := server.Execute(handlers.Index())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Add("My Idea", "My Idea Description", 1)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": idea.Number})

	code, _ := server.Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": "99"})

	code, _ := server.Execute(handlers.IdeaDetails())

	Expect(code).To(Equal(404))
}

func TestPostIdeaHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	code, _ := server.ExecutePost(handlers.PostIdea(), `{ "title": "My newest idea :)" }`)

	idea, err := ideas.GetByID(1)
	Expect(code).To(Equal(200))
	Expect(err).To(BeNil())
	Expect(idea.Title).To(Equal("My newest idea :)"))
	Expect(idea.TotalSupporters).To(Equal(1))
}

func TestPostIdeaHandler_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	code, _ := server.ExecutePost(handlers.PostIdea(), `{ "title": "" }`)

	_, err := ideas.GetByID(1)
	Expect(code).To(Equal(400))
	Expect(err).NotTo(BeNil())
}

func TestUpdateIdeaHandler_Authorized(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Add("My First Idea", "With a description", 1)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": idea.Number})
	code, _ := server.ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	idea, _ = ideas.GetByNumber(idea.Number)
	Expect(code).To(Equal(200))
	Expect(idea.Title).To(Equal("the new title"))
	Expect(idea.Description).To(Equal("new description"))
}

func TestUpdateIdeaHandler_NonAuthorized(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Add("My First Idea", "With a description", 1)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 2, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": idea.Number})
	code, _ := server.ExecutePost(handlers.UpdateIdea(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).To(Equal(401))
}

func TestPostCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Add("My First Idea", "With a description", 1)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": idea.Number})
	handler := handlers.PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": "1"})
	handler := handlers.PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "" }`)

	Expect(code).To(Equal(400))
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	ideas.Add("The Idea #1", "The Description #1", 1)
	ideas.Add("The Idea #2", "The Description #2", 1)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 2, Name: "Arya"})
	server.Context.SetParams(web.Map{"number": "2"})

	code, _ := server.Execute(handlers.AddSupporter())
	first, _ := ideas.GetByNumber(1)
	second, _ := ideas.GetByNumber(2)

	Expect(code).To(Equal(200))
	Expect(first.TotalSupporters).To(Equal(0))
	Expect(second.TotalSupporters).To(Equal(1))
}

func TestAddSupporterHandler_InvalidIdea(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParams(web.Map{"number": "1"})

	code, _ := server.Execute(handlers.AddSupporter())

	Expect(code).To(Equal(404))
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Add("The Idea #1", "The Description #1", 1)
	ideas.AddSupporter(1, 1)
	ideas.AddSupporter(1, 2)
	ideas.AddSupporter(1, 3)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{Ideas: ideas})
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 2, Name: "Arya"})
	server.Context.SetParams(web.Map{"number": idea.Number})

	code, _ := server.Execute(handlers.RemoveSupporter())
	idea, _ = ideas.GetByNumber(idea.Number)

	Expect(code).To(Equal(200))
	Expect(idea.TotalSupporters).To(Equal(2))
}
