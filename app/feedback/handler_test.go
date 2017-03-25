package feedback_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/feedback"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/mock"
	. "github.com/onsi/gomega"
)

type mockIdeaService struct{}

func (svc mockIdeaService) GetByID(tenantID, ideaID int64) (*feedback.Idea, error) {
	return new(feedback.Idea), nil
}

func (svc mockIdeaService) GetAll(tenantID int64) ([]*feedback.Idea, error) {
	return make([]*feedback.Idea, 0), nil
}

func (svc mockIdeaService) GetCommentsByIdeaID(tenantID, ideaID int64) ([]*feedback.Comment, error) {
	return make([]*feedback.Comment, 0), nil
}

func (svc mockIdeaService) Save(tenantID, userID int64, title, description string) (*feedback.Idea, error) {
	return nil, nil
}

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Set("Tenant", &identity.Tenant{ID: 2, Name: "Any Tenant"})
	code, _ := server.Execute(feedback.Index(&mockIdeaService{}).List())

	Expect(code).To(Equal(200))
}
