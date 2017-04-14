package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	expected string
	hosts    []string
}{
	{
		"The Orange Inc.",
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"Demonstration",
		[]string{
			"demo.test.canhearyou.com",
			"demo.test.canhearyou.com:1231",
			"demo.test.canhearyou.com:80",
		},
	},
}

func TestMultiTenant(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenants.Add(&models.Tenant{Name: "The Orange Inc.", Subdomain: "orange"})
	tenants.Add(&models.Tenant{Name: "Demonstration", Subdomain: "demo"})

	for _, testCase := range testCases {
		for _, host := range testCase.hosts {

			server := mock.NewServer()
			req, _ := http.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)
			c.Request().Host = host

			mw := middlewares.MultiTenant(tenants)
			mw(func(c web.Context) error {
				return c.String(http.StatusOK, c.Tenant().Name)
			})(c)

			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(Equal(testCase.expected))
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

	mw := middlewares.MultiTenant(&inmemory.TenantStorage{})
	mw(func(c web.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})(c)

	Expect(rec.Code).To(Equal(404))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "login.test.canhearyou.com"

	mw := middlewares.HostChecker("http://login.test.canhearyou.com")
	mw(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middlewares.HostChecker("login.test.canhearyou.com")
	mw(web.HandlerFunc(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}
