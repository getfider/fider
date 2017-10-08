package middlewares

import (
	"bufio"
	"compress/gzip"
	"net"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/web"
)

var (
	minSize = 1024
)

type gzipResponseWriter struct {
	code     int
	response http.ResponseWriter
	buffer   []byte
}

// Compress returns a middleware which compresses HTTP response using gzip compression
func Compress() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			res := c.Response
			res.Header().Add("Vary", "Accept-Encoding")
			if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
				gzipResponse := &gzipResponseWriter{response: res}
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
	w.buffer = append(w.buffer, b...)
	return len(b), nil
}

func (w *gzipResponseWriter) Flush() {
	if len(w.buffer) >= minSize {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
		w.response.WriteHeader(w.code)
		gzipWriter := gzip.NewWriter(w.response)
		gzipWriter.Write(w.buffer)
		gzipWriter.Flush()
		gzipWriter.Close()
	} else {
		w.response.WriteHeader(w.code)
		w.response.Write(w.buffer)
	}
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.response.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.response.(http.CloseNotifier).CloseNotify()
}
