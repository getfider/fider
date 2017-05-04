package middlewares_test

import (
	"net/http"
	"testing"

	"net/url"

	"github.com/WeCanHearYou/wechy/app/middlewares"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	"github.com/WeCanHearYou/wechy/app/pkg/web"
	"github.com/WeCanHearYou/wechy/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Use(middlewares.JwtGetter(nil))
	status, _ := server.Execute(func(c web.Context) error {
		if c.IsAuthenticated() {
			return c.NoContent(http.StatusOK)
		} else {
			return c.NoContent(http.StatusNoContent)
		}
	})

	Expect(status).To(Equal(http.StatusNoContent))
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
	server.Context.SetTenant(tenant)
	server.Context.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	server.Use(middlewares.JwtGetter(users))
	status, response := server.Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.User().Name)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserID:   300,
		UserName: "Jon Snow",
	})

	users := &inmemory.UserStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 400})
	server.Context.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	server.Use(middlewares.JwtGetter(users))
	status, _ := server.Execute(func(c web.Context) error {
		if c.User() == nil {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusNoContent))
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.canhearyou.com"
	server.Context.Request().RequestURI = "/abc"

	server.Use(middlewares.JwtSetter())
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.canhearyou.com"
	server.Context.Request().URL, _ = url.Parse("/abc?jwt=" + token)

	server.Use(middlewares.JwtSetter())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.WechyClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.canhearyou.com"
	server.Context.Request().URL, _ = url.Parse("/abc?jwt=" + token + "&foo=bar")

	server.Use(middlewares.JwtSetter())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.canhearyou.com/abc?foo=bar"))
}
