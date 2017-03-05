package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app/context"
	"github.com/WeCanHearYou/wchy/app/handler"
	"github.com/WeCanHearYou/wchy/app/mock"
	"github.com/WeCanHearYou/wchy/app/model"
	"github.com/WeCanHearYou/wchy/app/service"
	. "github.com/onsi/gomega"
)

type mockTenantService struct{}

func (svc mockTenantService) GetByDomain(domain string) (*model.Tenant, error) {
	if domain == "trishop" {
		return &model.Tenant{ID: 2, Name: "The Triathlon Shop", Domain: "trishop"}, nil
	}
	return nil, service.ErrNotFound
}

var ctx *context.WchyContext = &context.WchyContext{
	Tenant: &mockTenantService{},
}

func TestTenantHandler_404(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetParamNames("domain")
	server.Context.SetParamValues("unknown")
	code, _ := server.Execute(handler.TenantByDomain(ctx))

	Expect(code).To(Equal(404))
}

func TestTenantHandler_200(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.SetParamNames("domain")
	server.Context.SetParamValues("trishop")
	code, query := server.Execute(handler.TenantByDomain(ctx))

	Expect(query.Int("id")).To(Equal(2))
	Expect(query.String("name")).To(Equal("The Triathlon Shop"))
	Expect(query.String("domain")).To(Equal("trishop"))
	Expect(code).To(Equal(200))
}
