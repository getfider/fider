package middlewares_test

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestCompress(t *testing.T) {
	RegisterT(t)

	data := "Hello World\n"
	for i := 0; i <= 500; i++ {
		data += "Hello World\n"
	}

	server := mock.NewServer()
	server.Use(middlewares.Compress())
	handler := func(c *web.Context) error {
		return c.String(http.StatusOK, data)
	}

	status, response := server.
		AddHeader("Accept-Encoding", "gzip").
		Execute(handler)

	reader, _ := gzip.NewReader(response.Body)
	bytes, _ := ioutil.ReadAll(reader)
	Expect(bytes).Equals([]byte(data))
	Expect(status).Equals(http.StatusOK)
	Expect(response.Header().Get("Vary")).Equals("Accept-Encoding")
	Expect(response.Header().Get("Content-Encoding")).Equals("gzip")
}
