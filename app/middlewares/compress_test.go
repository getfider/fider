package middlewares_test

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestCompress(t *testing.T) {
	RegisterTestingT(t)

	data := "Hello World\n"
	for i := 0; i <= 500; i++ {
		data += "Hello World\n"
	}

	server, _ := mock.NewServer()
	server.Use(middlewares.Compress())
	handler := func(c web.Context) error {
		return c.String(http.StatusOK, data)
	}

	status, response := server.
		AddHeader("Accept-Encoding", "gzip").
		Execute(handler)

	reader, _ := gzip.NewReader(response.Body)
	Expect(ioutil.ReadAll(reader)).To(Equal([]byte(data)))
	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Vary")).To(Equal("Accept-Encoding"))
	Expect(response.Header().Get("Content-Encoding")).To(Equal("gzip"))
}

func TestCompress_SmallResponse(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.Compress())
	handler := func(c web.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}

	status, response := server.
		AddHeader("Accept-Encoding", "gzip").
		Execute(handler)

	Expect(ioutil.ReadAll(response.Body)).To(Equal([]byte("Hello World")))
	Expect(status).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Content-Encoding")).To(Equal(""))
}
