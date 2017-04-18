package handlers_test

import (
	"testing"

	"strconv"

	"github.com/WeCanHearYou/wechy/app/handlers"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestListHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})

	code, _ := server.Execute(handlers.Handlers(ideas).List())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	idea, _ := ideas.Save(1, 1, "My Idea", "My Idea Description")
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
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
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("99")

	code, _ := server.Execute(handlers.Handlers(ideas).Details())

	Expect(code).To(Equal(404))
}

func TestDetailsHandler_AddIdea(t *testing.T) {
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

func TestDetailsHandler_AddIdea_WithoutTitle(t *testing.T) {
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

func TestDetailsHandler_AddComment(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(ideas).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_AddComment_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	ideas := &inmemory.IdeaStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 1, Name: "Any Tenant"})
	server.Context.SetUser(&models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(ideas).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "" }`)

	Expect(code).To(Equal(400))
}
