package middlewares_test

import (
	"net/http"
	"testing"

	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/mock"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Use(middlewares.JwtGetter())
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
	token, _ := jwt.Encode(&models.FiderClaims{
		UserID:   user.ID,
		UserName: user.Name,
	})

	users := &inmemory.UserStorage{}

	users.Register(user)
	server := mock.NewServer()
	server.Context.SetServices(&app.Services{
		Users: users,
	})
	server.Context.SetTenant(tenant)
	server.Context.Request().Header.Add("Accept", "application/json")
	server.Context.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	server.Use(middlewares.JwtGetter())
	status, response := server.Execute(func(c web.Context) error {
		return c.String(http.StatusOK, c.User().Name)
	})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.FiderClaims{
		UserID:   300,
		UserName: "Jon Snow",
	})

	users := &inmemory.UserStorage{}
	server := mock.NewServer()
	server.Context.SetTenant(&models.Tenant{ID: 400})
	server.Context.SetServices(&app.Services{
		Users: users,
	})
	server.Context.Request().AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	server.Use(middlewares.JwtGetter())
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
	server.Context.Request().Host = "orange.test.fider.io"
	server.Context.Request().RequestURI = "/abc"

	server.Use(middlewares.JwtSetter())
	status, _ := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.FiderClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.fider.io"
	server.Context.Request().URL, _ = url.Parse("/abc?jwt=" + token)

	server.Use(middlewares.JwtSetter())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	token, _ := jwt.Encode(&models.FiderClaims{
		UserName: "Jon Snow",
	})

	server := mock.NewServer()
	server.Context.Request().Host = "orange.test.fider.io"
	server.Context.Request().URL, _ = url.Parse("/abc?jwt=" + token + "&foo=bar")

	server.Use(middlewares.JwtSetter())
	status, response := server.Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io/abc?foo=bar"))
}
