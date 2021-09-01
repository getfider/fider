package middlewares

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/web"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Compress returns a middleware which compresses HTTP response using gzip compression
func Compress() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			if strings.Contains(c.Request.GetHeader("Accept-Encoding"), "gzip") {
				res := c.Response
				res.Header().Set("Content-Encoding", "gzip")
				res.Header().Del("Accept-Encoding")
				res.Header().Add("Vary", "Accept-Encoding")

				gw, _ := gzip.NewWriterLevel(res.Writer, gzip.DefaultCompression)

				c.Response.Writer = &gzipResponseWriter{
					Writer:         gw,
					ResponseWriter: c.Response.Writer,
				}

				err := next(c)
				gw.Close()
				return err
			}
			return next(c)
		}
	}
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)

}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	header := w.ResponseWriter.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", http.DetectContentType(b))
	}
	header.Del("Content-Length")

	return w.Writer.Write(b)
}

func (r gzipResponseWriter) Flush() {
	r.Writer.(http.Flusher).Flush()
}

func (r gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}
