package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/service"
	. "github.com/onsi/gomega"
)

type inMemoryTenantService struct {
	tenants []*model.Tenant
}

func (svc inMemoryTenantService) GetByDomain(domain string) (*model.Tenant, error) {
	for _, tenant := range svc.tenants {
		if tenant.Domain == domain {
			return tenant, nil
		}
	}
	return nil, service.ErrNotFound
}

var ctx *context.WchyContext = &context.WchyContext{
	Tenant: &inMemoryTenantService{tenants: []*model.Tenant{
		&model.Tenant{ID: 1, Name: "Orange Inc.", Domain: "orange"},
		&model.Tenant{ID: 2, Name: "The Triathlon Shop", Domain: "trishop"},
	}},
}

func TestTenantHandler_404(t *testing.T) {
	RegisterTestingT(t)

	server := NewTestServer()
	server.param("domain", "unknown")
	server.register(handler.TenantByDomain(ctx))
	status, _ := server.request()

	Expect(status).To(Equal(404))
}

func TestTenantHandler_200(t *testing.T) {
	RegisterTestingT(t)

	server := NewTestServer()
	server.param("domain", "trishop")
	server.register(handler.TenantByDomain(ctx))
	status, query := server.request()

	Expect(query.Int("id")).To(Equal(2))
	Expect(query.String("name")).To(Equal("The Triathlon Shop"))
	Expect(query.String("domain")).To(Equal("trishop"))
	Expect(status).To(Equal(200))
}
