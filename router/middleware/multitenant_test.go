package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/router/middleware"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/labstack/echo"
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

			e := echo.New()
			req, _ := http.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Request().Host = host

			cors := middleware.MultiTenant(ctx)
			cors(echo.HandlerFunc(func(c echo.Context) error {
				return c.String(http.StatusOK, c.Get("Tenant").(*model.Tenant).Name)
			}))(c)

			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(Equal(testCase.tenant.Name))

		}
	}

}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "somedomain.com"

	cors := middleware.MultiTenant(ctx)
	cors(echo.HandlerFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, c.Get("Tenant").(*model.Tenant).Name)
	}))(c)

	Expect(rec.Code).To(Equal(404))
}
