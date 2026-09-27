package middlewares_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSecureWithoutCDN(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.Secure())

	var ctxID string
	status, response := server.Execute(func(c *web.Context) error {
		ctxID = c.ContextID()
		return c.NoContent(http.StatusOK)
	})

	expectedPolicy := "base-uri 'self'; default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'nonce-" + ctxID + "' https://www.google-analytics.com; img-src 'self' https: data:; font-src 'self' data:; object-src 'none'; media-src 'none'; connect-src 'self' https://www.google-analytics.com; frame-src 'self'"

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Security-Policy")).Equals(expectedPolicy)
	Expect(response.Header().Get("X-XSS-Protection")).Equals("1; mode=block")
	Expect(response.Header().Get("X-Content-Type-Options")).Equals("nosniff")
	Expect(response.Header().Get("Referrer-Policy")).Equals("no-referrer-when-downgrade")
}

func TestSecureWithCDN(t *testing.T) {
	RegisterT(t)

	env.Config.CDN.Host = "test.fider.io"

	server := mock.NewServer()
	server.Use(middlewares.Secure())

	var ctxID string
	status, response := server.Execute(func(c *web.Context) error {
		ctxID = c.ContextID()
		return c.NoContent(http.StatusOK)
	})

	expectedPolicy := "base-uri 'self'; default-src 'self'; style-src 'self' 'unsafe-inline' *.test.fider.io; script-src 'self' 'nonce-" + ctxID + "' https://www.google-analytics.com *.test.fider.io; img-src 'self' https: data: *.test.fider.io; font-src 'self' data: *.test.fider.io; object-src 'none'; media-src 'none'; connect-src 'self' https://www.google-analytics.com *.test.fider.io; frame-src 'self'"

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Security-Policy")).Equals(expectedPolicy)
	Expect(response.Header().Get("X-XSS-Protection")).Equals("1; mode=block")
	Expect(response.Header().Get("X-Content-Type-Options")).Equals("nosniff")
	Expect(response.Header().Get("Referrer-Policy")).Equals("no-referrer-when-downgrade")
}

func TestSecureWithCDN_SingleHost(t *testing.T) {
	RegisterT(t)

	env.Config.CDN.Host = "test.fider.io"

	server := mock.NewSingleTenantServer()
	server.Use(middlewares.Secure())

	var ctxID string
	status, response := server.WithURL("http://test.fider.io").Execute(func(c *web.Context) error {
		ctxID = c.ContextID()
		return c.NoContent(http.StatusOK)
	})

	expectedPolicy := "base-uri 'self'; default-src 'self'; style-src 'self' 'unsafe-inline' test.fider.io; script-src 'self' 'nonce-" + ctxID + "' https://www.google-analytics.com test.fider.io; img-src 'self' https: data: test.fider.io; font-src 'self' data: test.fider.io; object-src 'none'; media-src 'none'; connect-src 'self' https://www.google-analytics.com test.fider.io; frame-src 'self'"

	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Content-Security-Policy")).Equals(expectedPolicy)
	Expect(response.Header().Get("X-XSS-Protection")).Equals("1; mode=block")
	Expect(response.Header().Get("X-Content-Type-Options")).Equals("nosniff")
	Expect(response.Header().Get("Referrer-Policy")).Equals("no-referrer-when-downgrade")
}

func executeCSRF(method string, headers map[string]string) int {
	server := mock.NewServer()
	server.Use(func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Request.Method = method
			return next(c)
		}
	})
	server.Use(middlewares.CSRF())
	for k, v := range headers {
		server.AddHeader(k, v)
	}
	status, _ := server.Execute(func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	})
	return status
}

// See GHSA-xjr8-w967-4xjq
func TestCSRF(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		name     string
		method   string
		headers  map[string]string
		expected int
	}{
		{"GET without headers", "GET", map[string]string{}, http.StatusOK},
		{"GET cross-site", "GET", map[string]string{"Sec-Fetch-Site": "cross-site"}, http.StatusOK},
		{"HEAD without headers", "HEAD", map[string]string{}, http.StatusOK},
		{"POST without headers", "POST", map[string]string{}, http.StatusForbidden},
		{"POST with only Accept json", "POST", map[string]string{"Accept": "application/json"}, http.StatusForbidden},
		{"PUT with only Accept json", "PUT", map[string]string{"Accept": "application/json"}, http.StatusForbidden},
		{"DELETE with only Accept json", "DELETE", map[string]string{"Accept": "application/json"}, http.StatusForbidden},
		{"PATCH with only Accept json", "PATCH", map[string]string{"Accept": "application/json"}, http.StatusForbidden},
		{"POST with JSON content type", "POST", map[string]string{"Content-Type": "application/json"}, http.StatusOK},
		{"POST with JSON content type and charset", "POST", map[string]string{"Content-Type": "application/json; charset=utf-8"}, http.StatusOK},
		{"POST with upper-case JSON content type", "POST", map[string]string{"Content-Type": "Application/JSON"}, http.StatusOK},
		{"PUT with JSON content type", "PUT", map[string]string{"Content-Type": "application/json"}, http.StatusOK},
		{"DELETE with JSON content type", "DELETE", map[string]string{"Content-Type": "application/json"}, http.StatusOK},
		{"POST with text/plain smuggling json", "POST", map[string]string{"Content-Type": "text/plain; x=application/json"}, http.StatusForbidden},
		{"POST with text/plain and Accept json", "POST", map[string]string{"Content-Type": "text/plain", "Accept": "application/json"}, http.StatusForbidden},
		{"POST with form content type", "POST", map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, http.StatusForbidden},
		{"POST with multipart content type", "POST", map[string]string{"Content-Type": "multipart/form-data; boundary=x"}, http.StatusForbidden},
		{"POST with json-like content type", "POST", map[string]string{"Content-Type": "application/jsonx"}, http.StatusForbidden},
		{"POST with invalid content type", "POST", map[string]string{"Content-Type": "application/json;;;="}, http.StatusForbidden},
		{"POST with Bearer Authorization", "POST", map[string]string{"Authorization": "Bearer some-api-key"}, http.StatusOK},
		{"DELETE with Bearer Authorization", "DELETE", map[string]string{"Authorization": "Bearer some-api-key"}, http.StatusOK},
		{"POST with Basic Authorization", "POST", map[string]string{"Authorization": "Basic dXNlcjpwYXNz"}, http.StatusForbidden},
		{"POST with empty Bearer Authorization", "POST", map[string]string{"Authorization": "Bearer "}, http.StatusForbidden},
		{"POST cross-site with JSON content type", "POST", map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "cross-site"}, http.StatusForbidden},
		{"POST cross-site with Bearer Authorization", "POST", map[string]string{"Authorization": "Bearer some-api-key", "Sec-Fetch-Site": "cross-site"}, http.StatusForbidden},
		{"POST same-site with JSON content type", "POST", map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-site"}, http.StatusOK},
		{"POST same-origin with JSON content type", "POST", map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}, http.StatusOK},
		{"POST user-initiated with JSON content type", "POST", map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "none"}, http.StatusOK},
	}

	for _, tc := range testCases {
		status := executeCSRF(tc.method, tc.headers)
		if status != tc.expected {
			t.Errorf("%s: expected status %d, got %d", tc.name, tc.expected, status)
		}
	}
}

func TestAuthCookieIsSameSiteLax(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	_, response := server.Execute(func(c *web.Context) error {
		c.AddCookie(web.CookieAuthName, "token", time.Now().Add(time.Hour))
		return c.NoContent(http.StatusOK)
	})

	cookie := web.ParseCookie(response.Header().Get("Set-Cookie"))
	Expect(cookie.Name).Equals(web.CookieAuthName)
	Expect(cookie.SameSite).Equals(http.SameSiteLaxMode)
	Expect(cookie.HttpOnly).IsTrue()
}
