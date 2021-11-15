package web

import (
	"bufio"
	"net"
	"net/http"
)

// Response is a wrapper of http.ResponseWriter with access to the StatusCode
type Response struct {
	Writer     http.ResponseWriter
	StatusCode int
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) WriteHeader(code int) {
	r.StatusCode = code
	r.Writer.WriteHeader(code)
}

func (r *Response) Write(b []byte) (int, error) {
	if r.StatusCode == 0 {
		r.WriteHeader(http.StatusOK)
	}

	return r.Writer.Write(b)
}

func (r Response) Flush() {
	r.Writer.(http.Flusher).Flush()
}

func (r Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}
