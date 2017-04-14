package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)

	mw := middlewares.JwtGetter(nil)
	mw(func(c web.Context) error {
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

	tenant := &models.Tenant{ID: 300}
	user := &models.User{
		ID:     300,
		Name:   "Jon Snow",
		Tenant: tenant,
	}
	token, _ := jwt.Encode(&models.WechyClaims{
		UserID:   user.ID,
		UserName: user.Name,
	})

	users := &inmemory.UserStorage{}
	users.Register(user)
	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetTenant(tenant)
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := middlewares.JwtGetter(users)
	mw(func(c web.Context) error {
		return c.String(http.StatusOK, c.User().Name)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(rec.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserID:   300,
		UserName: "Jon Snow",
	})

	users := &inmemory.UserStorage{}
	server := mock.NewServer()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetTenant(&models.Tenant{ID: 400})
	c.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	mw := middlewares.JwtGetter(users)
	mw(func(c web.Context) error {
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
	mw(func(c web.Context) error {
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
	mw(func(c web.Context) error {
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
	mw(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})(c)

	Expect(rec.Code).To(Equal(http.StatusTemporaryRedirect))
	Expect(rec.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}
