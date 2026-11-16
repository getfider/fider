package web

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"syscall"
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

// isClientGone reports whether a response write failed because the client went away
// (closed the page, navigated off, dropped the connection) rather than a server fault.
func isClientGone(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	// net/http's bundled http2 server doesn't export these sentinels
	// (http2errStreamClosed, http2errClientDisconnected), so match their messages.
	msg := err.Error()
	return strings.Contains(msg, "http2: stream closed") || strings.Contains(msg, "client disconnected")
}
