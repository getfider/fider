package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPClient is an mocked implementation of HTTP Client
type HTTPClient struct {
	Requests []*http.Request
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Requests: make([]*http.Request, 0),
	}
}

// Do saves the request internally for later inspection
func (client *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	client.Requests = append(client.Requests, req)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}, nil
}

// Reset the internal request store
func (client *HTTPClient) Reset() {
	client.Requests = make([]*http.Request, 0)
}
