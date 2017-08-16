package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
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

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&models.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtGetter())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie("auth", token).
		Execute(func(c web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Body.String()).To(Equal("Jon Snow"))
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&models.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtGetter())
	status, _ := server.
		OnTenant(mock.OrangeTenant).
		AddCookie("auth", token).
		Execute(func(c web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).To(Equal(http.StatusNoContent))
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()

	server.Use(middlewares.JwtSetter())
	status, _ := server.WithURL("http://orange.test.fider.io/abc").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&models.FiderClaims{
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtSetter())
	status, response := server.WithURL("http://orange.test.fider.io/abc?jwt=" + token).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io/abc"))
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&models.FiderClaims{
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtSetter())
	status, response := server.WithURL("http://orange.test.fider.io/abc?jwt=" + token + "&foo=bar").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header().Get("Location")).To(Equal("http://orange.test.fider.io/abc?foo=bar"))
}
