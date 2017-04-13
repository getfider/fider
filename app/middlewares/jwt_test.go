package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

var testCases = []struct {
	domain string
	tenant *models.Tenant
	hosts  []string
}{
	{
		"orange.test.canhearyou.com",
		&models.Tenant{Name: "The Orange Inc."},
		[]string{
			"orange.test.canhearyou.com",
			"orange.test.canhearyou.com:3000",
		},
	},
	{
		"trishop.test.canhearyou.com",
		&models.Tenant{Name: "The Triathlon Shop"},
		[]string{
			"trishop.test.canhearyou.com",
			"trishop.test.canhearyou.com:1231",
			"trishop.test.canhearyou.com:80",
		},
	},
}

var (
	db            *dbx.Database
	TenantStorage *postgres.TenantStorage
	UserStorage   *postgres.UserStorage
)

func setup() {
	db, _ = dbx.New()

	TenantStorage = &postgres.TenantStorage{DB: db}
	UserStorage = &postgres.UserStorage{DB: db}
}

type mockTenantStorage struct{}

func (svc mockTenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
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

			mw := middlewares.MultiTenant(&mockTenantStorage{})
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

	mw := middlewares.MultiTenant(&mockTenantStorage{})
	mw(func(c app.Context) error {
		return c.String(http.StatusOK, c.Tenant().Name)
	})(c)

	Expect(rec.Code).To(Equal(404))
}

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)
	setup()

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := middlewares.JwtGetter(UserStorage)
	mw(func(c app.Context) error {
		if c.IsAuthenticated() {
			return c.NoContent(http.StatusOK)
		} else {
			return c.NoContent(http.StatusNoContent)
		}
	})(c)

	Expect(rec.Code).To(Equal(http.StatusNoContent))
}

func TestJwtGetter_WithCookie(t *testing.T) {
	RegisterTestingT(t)
	setup()

	token, _ := jwt.Encode(&models.WechyClaims{
		UserID:   300,
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetTenant(&models.Tenant{ID: 300})
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := middlewares.JwtGetter(UserStorage)
	mw(func(c app.Context) error {
		return c.String(http.StatusOK, c.User().Name)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(rec.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterTestingT(t)
	setup()

	token, _ := jwt.Encode(&models.WechyClaims{
		UserID:   300,
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetTenant(&models.Tenant{ID: 400})
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := middlewares.JwtGetter(UserStorage)
	mw(func(c app.Context) error {
		if c.User() == nil {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusNoContent))
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middlewares.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token, nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middlewares.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token+"&foo=bar", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middlewares.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}
