package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestJwtGetter_NoCookie(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.JwtGetter())
	status, _ := server.Execute(func(c web.Context) error {
		if c.IsAuthenticated() {
			return c.NoContent(http.StatusOK)
		} else {
			return c.NoContent(http.StatusNoContent)
		}
	})

	Expect(status).Equals(http.StatusNoContent)
}

func TestJwtGetter_WithCookie(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtGetter())
	status, response := server.
		OnTenant(mock.DemoTenant).
		AddHeader("Accept", "application/json").
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			return c.String(http.StatusOK, c.User().Name)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Body.String()).Equals("Jon Snow")
}

func TestJwtGetter_WithCookie_InvalidUser(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&jwt.FiderClaims{
		UserID:   999,
		UserName: "Unknown",
	})

	server.Use(middlewares.JwtGetter())
	status, response := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusNoContent)
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=;")
}

func TestJwtGetter_WithCookie_DifferentTenant(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&jwt.FiderClaims{
		UserID:   mock.JonSnow.ID,
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtGetter())
	status, _ := server.
		OnTenant(mock.AvengersTenant).
		AddCookie(web.CookieAuthName, token).
		Execute(func(c web.Context) error {
			if c.User() == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusNoContent)
}

func TestJwtSetter_WithoutJwt(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	server.Use(middlewares.JwtSetter())
	status, _ := server.WithURL("http://avengers.test.fider.io/abc").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestJwtSetter_WithJwt_WithoutParameter(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&jwt.FiderClaims{
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtSetter())
	status, response := server.WithURL("http://avengers.test.fider.io/abc?token=" + token).Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/abc")
}

func TestJwtSetter_WithJwt_WithParameter(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	token, _ := jwt.Encode(&jwt.FiderClaims{
		UserName: mock.JonSnow.Name,
	})

	server.Use(middlewares.JwtSetter())
	status, response := server.WithURL("http://avengers.test.fider.io/abc?token=" + token + "&foo=bar").Execute(func(c web.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/abc?foo=bar")
}
