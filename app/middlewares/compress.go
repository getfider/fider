package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/web"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
}

// Compress returns a middleware which compresses HTTP response using gzip compression
func Compress() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			res := c.Response
			if strings.Contains(c.Request.GetHeader("Accept-Encoding"), "gzip") {

				gw, _ := gzip.NewWriterLevel(res, gzip.DefaultCompression)
				defer gw.Close()

				res.Header().Set("Content-Encoding", "gzip")
				res.Header().Del("Accept-Encoding")
				res.Header().Add("Vary", "Accept-Encoding")

				h, hok := res.(http.Hijacker)
				if !hok { /* w is not Hijacker... oh well... */
					h = nil
				}

				f, fok := res.(http.Flusher)
				if !fok {
					f = nil
				}

				c.Response = &gzipResponseWriter{
					Writer:         gw,
					ResponseWriter: res,
					Hijacker:       h,
					Flusher:        f,
				}
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

type flusher interface {
	Flush() error
}

func (w *gzipResponseWriter) Flush() {
	if f, ok := w.Writer.(flusher); ok {
		f.Flush()
	}
	if w.Flusher != nil {
		w.Flusher.Flush()
	}
}
