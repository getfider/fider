package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getfider/fider/app/pkg/env"

	. "github.com/getfider/fider/app/pkg/assert"
)

func newBodyEchoEngine() *Engine {
	e := New()
	e.Post("/echo", func(c *Context) error {
		return c.String(http.StatusOK, c.Request.Body)
	})
	e.NotFound(func(c *Context) error {
		return c.String(http.StatusNotFound, c.Request.Body)
	})
	return e
}

func withMaxBodySize(t *testing.T, size int64) {
	original := env.Config.HTTP.MaxBodySize
	env.Config.HTTP.MaxBodySize = size
	t.Cleanup(func() { env.Config.HTTP.MaxBodySize = original })
}

func TestEngine_MaxBodySize_DefaultFitsLargestUpload(t *testing.T) {
	RegisterT(t)

	// 3 attachments of 5120 KB each, base64-encoded
	largestUpload := int64(3 * 5120 * 1024 * 4 / 3)
	Expect(env.Config.HTTP.MaxBodySize > largestUpload).IsTrue()
}

func TestEngine_BodyWithinLimit(t *testing.T) {
	RegisterT(t)
	withMaxBodySize(t, 10)

	e := newBodyEchoEngine()
	res := httptest.NewRecorder()
	e.mux.ServeHTTP(res, httptest.NewRequest("POST", "/echo", strings.NewReader("0123456789")))

	Expect(res.Code).Equals(http.StatusOK)
	Expect(res.Body.String()).Equals("0123456789")
}

func TestEngine_BodyTooLarge(t *testing.T) {
	RegisterT(t)
	withMaxBodySize(t, 10)

	e := newBodyEchoEngine()
	for _, path := range []string{"/echo", "/unknown-route"} {
		res := httptest.NewRecorder()
		e.mux.ServeHTTP(res, httptest.NewRequest("POST", path, strings.NewReader("01234567890")))

		Expect(res.Code).Equals(http.StatusRequestEntityTooLarge)
		Expect(res.Header().Get("Content-Type")).Equals(UTF8JSONContentType)
		Expect(res.Body.String()).ContainsSubstring(`"errors"`)
	}
}

func TestEngine_ChunkedBodyIsCapped(t *testing.T) {
	RegisterT(t)
	withMaxBodySize(t, 10)

	var readErr error
	e := New()
	e.Post("/raw", func(c *Context) error {
		_, readErr = io.ReadAll(c.Request.instance.Body)
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/raw", strings.NewReader("01234567890"))
	req.ContentLength = -1
	e.mux.ServeHTTP(httptest.NewRecorder(), req)

	Expect(readErr).IsNotNil()
}
