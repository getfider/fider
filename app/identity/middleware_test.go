package identity_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	domain string
	tenant *app.Tenant
	hosts  []string
}{
	{
		"orange.test.canhearyou.com",
		&app.Tenant{Name: "The Orange Inc."},
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"trishop.test.canhearyou.com",
		&app.Tenant{Name: "The Triathlon Shop"},
		[]string{
			"trishop.test.canhearyou.com",
			"trishop.test.canhearyou.com:1231",
			"trishop.test.canhearyou.com:80",
		},
	},
}

type mockTenantService struct{}

func (svc mockTenantService) GetByDomain(domain string) (*app.Tenant, error) {
	for _, testCase := range testCases {
		if testCase.domain == domain {
			return testCase.tenant, nil
		}
	}
	return nil, app.ErrNotFound
}

func TestMultiTenant(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range testCases {
		for _, host := range testCase.hosts {

			server := mock.NewServer()
			req, _ := http.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)
			c.Request().Host = host

			mw := identity.MultiTenant(&mockTenantService{})
			mw(func(c app.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})(c)

			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(Equal(testCase.tenant.Name))
		}
	}
}

func TestMultiTenant_UnknownDomain(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "somedomain.com"

	mw := identity.MultiTenant(&mockTenantService{})
	mw(func(c app.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})(c)

	Expect(rec.Code).To(Equal(404))
}
