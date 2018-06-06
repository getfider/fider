package middlewares

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/getfider/fider/app/pkg/web"
)

var (
	minSize = 1024
)

type gzipResponseWriter struct {
	code     int
	response http.ResponseWriter
	writer   *gzip.Writer
	buffer   *bytes.Buffer
}

// Compress returns a middleware which compresses HTTP response using gzip compression
func Compress() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		var pool sync.Pool
		pool.New = func() interface{} {
			return gzip.NewWriter(ioutil.Discard)
		}
		return func(c web.Context) error {
			res := c.Response
			if strings.Contains(c.Request.GetHeader("Accept-Encoding"), "gzip") {
				writer := pool.Get().(*gzip.Writer)
				defer pool.Put(writer)
				gzipResponse := &gzipResponseWriter{response: res, writer: writer, buffer: &bytes.Buffer{}}
				c.Response = gzipResponse
				defer gzipResponse.Flush()
			}
			return next(c)
		}
	}
}
func (w *gzipResponseWriter) Header() http.Header {
	return w.response.Header()
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	w.code = code
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.Header().Add("Vary", "Accept-Encoding")
	if w.buffer.Len() >= minSize {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
		w.response.WriteHeader(w.code)
		w.writer.Reset(w.response)
		w.buffer.WriteTo(w.writer)
		w.writer.Close()
	} else {
		w.response.WriteHeader(w.code)
		w.buffer.WriteTo(w.response)
	}
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.response.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.response.(http.CloseNotifier).CloseNotify()
}
