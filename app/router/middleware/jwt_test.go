package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/router/middleware"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := middleware.JwtGetter()
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

	token, _ := auth.Encode(&auth.WchyClaims{
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

	mw := middleware.JwtGetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		claims := c.Get("Claims").(*auth.WchyClaims)
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

	mw := middleware.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := auth.Encode(&auth.WchyClaims{
		UserName: "Jon Snow",
	})

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middleware.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := auth.Encode(&auth.WchyClaims{
		UserName: "Jon Snow",
	})

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/abc?jwt="+token+"&foo=bar", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "orange.test.canhearyou.com"

	mw := middleware.JwtSetter()
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}
