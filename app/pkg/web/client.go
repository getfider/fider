package web

import (
	"net/http"
)

// Client is an interface for HTTP Client
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClient is an implementation of HTTP Client
type HTTPClient struct {
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Do sends the request using http.DefaultClient
func (client *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
