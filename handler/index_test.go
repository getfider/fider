package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/mock"
	"github.com/WeCanHearYou/wchy/model"
	. "github.com/onsi/gomega"
)

type mockIdeaService struct{}

func (svc mockIdeaService) GetAll(tenantID int64) ([]*model.Idea, error) {
	return make([]*model.Idea, 0), nil
}

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	ctx := &context.WchyContext{
		Idea: &mockIdeaService{},
	}

	server := mock.NewServer()
	server.Context.Set("Tenant", &model.Tenant{ID: 2, Name: "Any Tenant"})
	code, _ := server.Execute(handler.Index(ctx))

	Expect(code).To(Equal(200))
}
