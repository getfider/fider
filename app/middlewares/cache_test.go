package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := middlewares.OneYearCache()
	mw(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(rec.Header().Get("Cache-Control")).To(Equal("public, max-age=30672000"))
}
