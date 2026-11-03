package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/imagic"
)

func TestServeProcessedImage(t *testing.T) {
	RegisterT(t)

	original := []byte("original")
	result := []byte("result")

	serve := func(err error) (int, http.Header, string) {
		code, response := mock.NewServer().
			OnTenant(mock.DemoTenant).
			Execute(func(c *web.Context) error {
				return serveProcessedImage(c, "image/png", original, result, err)
			})
		return code, response.Header(), response.Body.String()
	}

	code, _, body := serve(nil)
	Expect(code).Equals(http.StatusOK)
	Expect(body).Equals("result")

	// Cancelled requests: nothing is written, and no error is returned (so nothing is logged)
	for _, err := range []error{context.Canceled, context.DeadlineExceeded, fmt.Errorf("wrapped: %w", context.Canceled)} {
		code, _, body = serve(err)
		Expect(code).Equals(http.StatusOK)
		Expect(body).Equals("")
	}

	// Busy: a non cached 503, so the original isn't cached as the resized image
	code, header, _ := serve(validate.ErrDecodeBusy)
	Expect(code).Equals(http.StatusServiceUnavailable)
	Expect(header.Get("Retry-After")).Equals("5")
	Expect(header.Get("Cache-Control")).Equals("no-cache, no-store")

	// Anything else about the image itself (too large, unsupported, corrupt pixel data): the original
	for _, err := range []error{validate.ErrImageTooLarge, imagic.ErrNotSupported, errors.New("gif: missing image data")} {
		code, _, body = serve(err)
		Expect(code).Equals(http.StatusOK)
		Expect(body).Equals("original")
	}
}
