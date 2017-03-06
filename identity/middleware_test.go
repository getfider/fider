package identity_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/identity"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	domain string
	tenant *identity.Tenant
	hosts  []string
}{
	{
		"orange.test.canhearyou.com",
		&identity.Tenant{Name: "The Orange Inc."},
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"trishop.test.canhearyou.com",
		&identity.Tenant{Name: "The Triathlon Shop"},
		[]string{
			"trishop.test.canhearyou.com",
			"trishop.test.canhearyou.com:1231",
			"trishop.test.canhearyou.com:80",
		},
	},
}

type mockTenantService struct{}

func (svc mockTenantService) GetByDomain(domain string) (*identity.Tenant, error) {
	for _, testCase := range testCases {
		if testCase.domain == domain {
			return testCase.tenant, nil
		}
	}
	return nil, identity.ErrNotFound
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

			mw := identity.MultiTenant(&mockTenantService{})
			mw(echo.HandlerFunc(func(c echo.Context) error {
				return c.String(http.StatusOK, c.Get("Tenant").(*identity.Tenant).Name)
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

	mw := identity.MultiTenant(&mockTenantService{})
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, c.Get("Tenant").(*identity.Tenant).Name)
	}))(c)

	Expect(rec.Code).To(Equal(404))
}

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := identity.JwtGetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		if c.Get("Claims") == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			return c.NoContent(http.StatusOK)
		}
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusNoContent))
}

func TestJwtGetter_WithCookie(t *testing.T) {
	RegisterTestingT(t)

	token, _ := identity.Encode(&identity.WchyClaims{
		UserName: "Jon Snow",
	})

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := identity.JwtGetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		claims := c.Get("Claims").(*identity.WchyClaims)
		return c.String(http.StatusOK, claims.UserName)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(rec.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := identity.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := identity.Encode(&identity.WchyClaims{
		UserName: "Jon Snow",
	})

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := identity.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := identity.Encode(&identity.WchyClaims{
		UserName: "Jon Snow",
	})

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token+"&foo=bar", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := identity.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "login.test.canhearyou.com"

	mw := identity.HostChecker("http://login.test.canhearyou.com")
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestHostChecker_DifferentHost(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := identity.HostChecker("login.test.canhearyou.com")
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}
