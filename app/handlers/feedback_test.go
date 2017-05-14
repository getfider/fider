package handlers_test

import (
	"testing"

	"strconv"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestListHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})

	code, _ := server.Execute(handlers.Handlers(ideas).List())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Save(1, 1, "My Idea", "My Idea Description")
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues(strconv.Itoa(idea.Number))

	code, _ := server.Execute(handlers.Handlers(ideas).Details())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("99")

	code, _ := server.Execute(handlers.Handlers(ideas).Details())

	Expect(code).To(Equal(404))
}

func TestPostIdeaHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	handler := handlers.Handlers(ideas).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "My newest idea :)" }`)

	idea, err := ideas.GetByID(1, 1)
	Expect(code).To(Equal(200))
	Expect(err).To(BeNil())
	Expect(idea.Title).To(Equal("My newest idea :)"))
	Expect(idea.TotalSupporters).To(Equal(1))
}

func TestPostIdeaHandler_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	handler := handlers.Handlers(ideas).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "" }`)

	_, err := ideas.GetByID(1, 1)
	Expect(code).To(Equal(400))
	Expect(err).NotTo(BeNil())
}

func TestPostCommentHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	ideas.Save(1, 1, "Title", "Description")
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(ideas).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(ideas).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "" }`)

	Expect(code).To(Equal(400))
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	ideas.Save(1, 1, "The Idea #1", "The Description #1")
	ideas.Save(1, 1, "The Idea #2", "The Description #2")
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 2, Name: "Arya"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("2")

	code, _ := server.Execute(handlers.Handlers(ideas).AddSupporter())
	first, _ := ideas.GetByNumber(1, 1)
	second, _ := ideas.GetByNumber(1, 2)

	Expect(code).To(Equal(200))
	Expect(first.TotalSupporters).To(Equal(0))
	Expect(second.TotalSupporters).To(Equal(1))
}

func TestAddSupporterHandler_InvalidIdea(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")

	code, _ := server.Execute(handlers.Handlers(ideas).AddSupporter())

	Expect(code).To(Equal(404))
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	ideas.Save(1, 1, "The Idea #1", "The Description #1")
	ideas.AddSupporter(1, 1)
	ideas.AddSupporter(2, 1)
	ideas.AddSupporter(3, 1)
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 2, Name: "Arya"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")

	code, _ := server.Execute(handlers.Handlers(ideas).RemoveSupporter())
	idea, _ := ideas.GetByNumber(1, 1)

	Expect(code).To(Equal(200))
	Expect(idea.TotalSupporters).To(Equal(2))
}
