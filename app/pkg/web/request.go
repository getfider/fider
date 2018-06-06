package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

//Request wraps the http request object
type Request struct {
	instance      *http.Request
	Method        string
	ContentLength int64
	Body          io.ReadCloser
	URL           *url.URL
}

// WrapRequest returns Fider wrapper of HTTP Request
func WrapRequest(request *http.Request) Request {
	protocol := "http"
	if request.TLS != nil || request.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}

	host := request.Host
	if request.Header.Get("X-Forwarded-Host") != "" {
		host = request.Header.Get("X-Forwarded-Host")
	}

	fullURL := protocol + "://" + host + request.RequestURI
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Failed to parse url '%s'", fullURL)))
	}

	return Request{
		instance:      request,
		Method:        request.Method,
		ContentLength: request.ContentLength,
		Body:          request.Body,
		URL:           u,
	}
}

// GetHeader returns the value of HTTP header from given key
func (r *Request) GetHeader(key string) string {
	return r.instance.Header.Get(key)
}

// SetHeader updates the value of HTTP header of given key
func (r *Request) SetHeader(key, value string) {
	r.instance.Header.Set(key, value)
}

// Cookie returns the named cookie provided in the request.
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	return r.instance.Cookie(name)
}

// AddCookie adds a cookie
func (r *Request) AddCookie(cookie *http.Cookie) {
	r.instance.AddCookie(cookie)
}

// Subdomain returns the Fider subdomain (if available) from given host
func (r *Request) Subdomain() string {
	return ExtractSubdomain(r.URL.Hostname())
}

// ExtractSubdomain returns the Fider subdomain (if available) from given host
func ExtractSubdomain(host string) string {
	domain := env.MultiTenantDomain()
	if domain != "" && strings.Contains(host, domain) {
		return strings.Replace(host, domain, "", -1)
	}
	return ""
}
