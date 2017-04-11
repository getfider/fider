package feedback_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/feedback"
	"github.com/WeCanHearYou/wechy/app/mock"
	. "github.com/onsi/gomega"
)

type mockIdeaService struct {
	ideas []*feedback.Idea
}

func NewMockIdeaService() *mockIdeaService {
	return &mockIdeaService{
		ideas: []*feedback.Idea{
			&feedback.Idea{ID: 1, Number: 1, Title: "Idea #1"},
			&feedback.Idea{ID: 2, Number: 2, Title: "Idea #2"},
		},
	}
}

func (svc mockIdeaService) GetByID(tenantID, ideaID int) (*feedback.Idea, error) {
	for _, idea := range svc.ideas {
		if idea.ID == ideaID {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

func (svc mockIdeaService) GetByNumber(tenantID, number int) (*feedback.Idea, error) {
	for _, idea := range svc.ideas {
		if idea.Number == number {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

func (svc mockIdeaService) GetAll(tenantID int) ([]*feedback.Idea, error) {
	return svc.ideas, nil
}

func (svc mockIdeaService) GetCommentsByIdeaID(tenantID, ideaID int) ([]*feedback.Comment, error) {
	return make([]*feedback.Comment, 0), nil
}

func (svc mockIdeaService) Save(tenantID, userID int, title, description string) (*feedback.Idea, error) {
	return nil, nil
}

func (svc mockIdeaService) AddComment(userID, ideaID int, content string) (int, error) {
	return 0, nil
}

func TestListHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	code, _ := server.Execute(feedback.Handlers(NewMockIdeaService()).List())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("1")

	code, _ := server.Execute(feedback.Handlers(NewMockIdeaService()).Details())

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.SetParamNames("number")
	server.Context.SetParamValues("99")

	code, _ := server.Execute(feedback.Handlers(NewMockIdeaService()).Details())

	Expect(code).To(Equal(404))
}

func TestDetailsHandler_AddIdea(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &app.User{ID: 1, Name: "Jon"})
	handler := feedback.Handlers(NewMockIdeaService()).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "My newest idea :)" }`)

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_AddIdea_WithoutTitle(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &app.User{ID: 1, Name: "Jon"})
	handler := feedback.Handlers(NewMockIdeaService()).PostIdea()
	code, _ := server.ExecutePost(handler, `{ "title": "" }`)

	Expect(code).To(Equal(400))
}

func TestDetailsHandler_AddComment(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &app.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := feedback.Handlers(NewMockIdeaService()).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "This is a comment!" }`)

	Expect(code).To(Equal(200))
}

func TestDetailsHandler_AddComment_WithoutContent(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("__CTX_TENANT", &app.Tenant{ID: 2, Name: "Any Tenant"})
	server.Context.Set("__CTX_USER", &app.User{ID: 1, Name: "Jon"})
	server.Context.SetParamNames("id")
	server.Context.SetParamValues("1")
	handler := feedback.Handlers(NewMockIdeaService()).PostComment()
	code, _ := server.ExecutePost(handler, `{ "content": "" }`)

	Expect(code).To(Equal(400))
}
