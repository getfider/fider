package http

import (
	"io"
	"net/http"
	gohttp "net/http"
	"time"

	"github.com/getfider/fider/app/pkg/env"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := gohttp.DefaultTransport.(*gohttp.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

// Client is an interface for HTTP Client
type Client interface {
	Do(req *gohttp.Request) (*gohttp.Response, error)
}

// DefaultClient is the default HTTP Client
type DefaultClient struct {
}

// NewClient creates a new defaultClient
func NewClient() *DefaultClient {
	return &DefaultClient{}
}

// NewRequest creates a new HTTP Request
func NewRequest(method, url string, body io.Reader) (*gohttp.Request, error) {
	return gohttp.NewRequest(method, url, body)
}

// Do sends the request using http.DefaultClient
func (client *DefaultClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
