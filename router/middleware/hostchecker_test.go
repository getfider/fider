package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/router/middleware"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestHostChecker(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Host = "login.test.canhearyou.com"

	mw := middleware.HostChecker("login.test.canhearyou.com")
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

	mw := middleware.HostChecker("login.test.canhearyou.com")
	mw(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))(c)

	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}
