package middlewares_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSession_New(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.Session())

	var sessionID string
	status, response := server.Execute(func(c web.Context) error {
		sessionID = c.SessionID()
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
	cookie := web.ParseCookie(response.Header().Get("Set-Cookie"))
	Expect(cookie.Name).Equals(web.CookieSessionName)
	Expect(cookie.Value).Equals(sessionID)
	Expect(cookie.HttpOnly).IsTrue()
	Expect(cookie.Path).Equals("/")
	Expect(cookie.Expires).TemporarilySimilar(time.Now().Add(365*24*time.Hour), 5*time.Second)
}

func TestSession_ExistingSession(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.Session())

	status, response := server.
		AddCookie(web.CookieSessionName, "MY_SESSION_VALUE").
		Execute(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Set-Cookie")).Equals("")
}
