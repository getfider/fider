package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/mock"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/service"
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
	server.Param("domain", "unknown")
	server.Register(handler.TenantByDomain(ctx))
	status, _ := server.Request()

	Expect(status).To(Equal(404))
}

func TestTenantHandler_200(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Param("domain", "trishop")
	server.Register(handler.TenantByDomain(ctx))
	status, query := server.Request()

	Expect(query.Int("id")).To(Equal(2))
	Expect(query.String("name")).To(Equal("The Triathlon Shop"))
	Expect(query.String("domain")).To(Equal("trishop"))
	Expect(status).To(Equal(200))
}
