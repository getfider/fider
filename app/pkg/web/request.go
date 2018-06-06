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
	Protocol      string
	Host          string
	Subdomain     string
	Port          string
	Query         url.Values
	Path          string
	URL           string
	FullURL       string
}

// WrapRequest returns Fider wrapper of HTTP Request
func WrapRequest(request *http.Request) Request {
	r := Request{
		instance:      request,
		Method:        request.Method,
		ContentLength: request.ContentLength,
		Body:          request.Body,
	}

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

	r.SetURL(u)
	return r
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

// SetURL change current request URL
func (r *Request) SetURL(u *url.URL) {
	r.Protocol = u.Scheme
	r.Host = u.Hostname()
	r.Port = u.Port()
	r.Query = u.Query()
	r.Path = u.Path
	r.URL = u.RequestURI()
	r.FullURL = u.String()
	r.Subdomain = ExtractSubdomain(r.Host)
}

// ExtractSubdomain returns the Fider subdomain (if available) from given host
func ExtractSubdomain(host string) string {
	domain := env.MultiTenantDomain()
	if domain != "" && strings.Contains(host, domain) {
		return strings.Replace(host, domain, "", -1)
	}
	return ""
}
