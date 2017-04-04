package infra_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestIsAuthenticated_WithToken(t *testing.T) {
	RegisterTestingT(t)

	token, _ := infra.Encode(&app.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := infra.IsAuthenticated()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestIsAuthenticated_WithoutToken(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := infra.IsAuthenticated()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusForbidden))
}

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := infra.JwtGetter()
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

	token, _ := infra.Encode(&app.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := infra.JwtGetter()
	mw(func(c app.Context) error {
		return c.String(http.StatusOK, c.Claims().UserName)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(rec.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := infra.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := infra.Encode(&app.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token, nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := infra.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := infra.Encode(&app.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token+"&foo=bar", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := infra.JwtSetter()
	mw(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.Request().Host = "login.test.canhearyou.com"

	mw := infra.HostChecker("http://login.test.canhearyou.com")
	mw(func(c app.Context) error {
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

	mw := infra.HostChecker("login.test.canhearyou.com")
	mw(app.HandlerFunc(func(c app.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}
