package router_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/mock"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/router"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	domain string
	tenant *model.Tenant
	hosts  []string
}{
	{
		"orange.test.canhearyou.com",
		&model.Tenant{Name: "The Orange Inc."},
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"trishop.test.canhearyou.com",
		&model.Tenant{Name: "The Triathlon Shop"},
		[]string{
			"trishop.test.canhearyou.com",
			"trishop.test.canhearyou.com:1231",
			"trishop.test.canhearyou.com:80",
		},
	},
}

type mockTenantService struct{}

func (svc mockTenantService) GetByDomain(domain string) (*model.Tenant, error) {
	for _, testCase := range testCases {
		if testCase.domain == domain {
			return testCase.tenant, nil
		}
	}
	return nil, service.ErrNotFound
}

var ctx *context.WchyContext = &context.WchyContext{
	Tenant: &mockTenantService{},
}

func TestMultiTenant(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range testCases {
		for _, host := range testCase.hosts {
			var name interface{}

			server := mock.NewServer()
			server.Use(func(c *gin.Context) {
				c.Request.Host = host
			})
			server.Use(router.MultiTenant(ctx))
			server.Register(func(c *gin.Context) {
				c.Status(200)
				name = c.MustGet("Tenant")
			})
			status, _ := server.Request()

			Expect(status).To(Equal(200))
			Expect(name).To(Equal(testCase.tenant))
		}
	}

}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Use(func(c *gin.Context) {
		c.Request.Host = "somedomain.com"
	})
	server.Use(router.MultiTenant(ctx))
	server.Register(func(c *gin.Context) {
		c.Status(200)
	})
	status, _ := server.Request()

	Expect(status).To(Equal(404))
}
