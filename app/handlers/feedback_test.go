package handlers_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/handlers"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	. "github.com/onsi/gomega"
)

type mockIdeaStorage struct {
	ideas []*models.Idea
}

func NewMockIdeaStorage() *mockIdeaStorage {
	return &mockIdeaStorage{
		ideas: []*models.Idea{
			&models.Idea{ID: 1, Number: 1, Title: "Idea #1"},
			&models.Idea{ID: 2, Number: 2, Title: "Idea #2"},
		},
	}
}

func (svc mockIdeaStorage) GetByID(tenantID, ideaID int) (*models.Idea, error) {
	for _, idea := range svc.ideas {
		if idea.ID == ideaID {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

func (svc mockIdeaStorage) GetByNumber(tenantID, number int) (*models.Idea, error) {
	for _, idea := range svc.ideas {
		if idea.Number == number {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

func (svc mockIdeaStorage) GetAll(tenantID int) ([]*models.Idea, error) {
	return svc.ideas, nil
}

func (svc mockIdeaStorage) GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error) {
	return make([]*models.Comment, 0), nil
}

func (svc mockIdeaStorage) Save(tenantID, userID int, title, description string) (*models.Idea, error) {
	return nil, nil
}

func (svc mockIdeaStorage) AddComment(userID, ideaID int, content string) (int, error) {
	return 0, nil
}

func TestListHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	code, _ := server.Execute(handlers.Handlers(NewMockIdeaStorage()).List())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")

	code, _ := server.Execute(handlers.Handlers(NewMockIdeaStorage()).Details())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("99")

	code, _ := server.Execute(handlers.Handlers(NewMockIdeaStorage()).Details())

	Expect(code).To(Equal(404))
}

func TestDetailsHandler_AddIdea(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &models.User{ID: 1, Name: "Jon"})
	handler := handlers.Handlers(NewMockIdeaStorage()).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "My newest idea :)" }`)

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_AddIdea_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &models.User{ID: 1, Name: "Jon"})
	handler := handlers.Handlers(NewMockIdeaStorage()).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "" }`)

	Expect(code).To(Equal(400))
}

func TestDetailsHandler_AddComment(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(NewMockIdeaStorage()).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_AddComment_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &models.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &models.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := handlers.Handlers(NewMockIdeaStorage()).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "" }`)

	Expect(code).To(Equal(400))
}
