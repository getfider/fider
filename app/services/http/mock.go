package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	gohttp "net/http"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := gohttp.DefaultTransport.(*gohttp.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

type MockHTTPClientService struct{}

var MockHTTPRequestsHistory = make([]*gohttp.Request, 0)

func (s MockHTTPClientService) Enabled() bool {
	return true
}

func (s MockHTTPClientService) Init() {
	MockHTTPRequestsHistory = make([]*gohttp.Request, 0)
	bus.AddHandler(s, MockHTTPGetRequest)
	bus.AddHandler(s, MockHTTPPostRequest)
}

func MockHTTPPostRequest(ctx context.Context, cmd *HTTPPostRequestCommand) error {
	req, err := gohttp.NewRequest("POST", cmd.URL, cmd.Body)
	if err != nil {
		return err
	}

	for k, v := range cmd.Headers {
		req.Header.Set(k, v)
	}
	if cmd.BasicAuth != nil {
		req.SetBasicAuth(cmd.BasicAuth.User, cmd.BasicAuth.Password)
	}

	MockHTTPRequestsHistory = append(MockHTTPRequestsHistory, req)

	cmd.Response = &gohttp.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	return nil
}

func MockHTTPGetRequest(ctx context.Context, cmd *HTTPGetRequestCommand) error {
	req, err := gohttp.NewRequest("GET", cmd.URL, nil)
	if err != nil {
		return err
	}

	for k, v := range cmd.Headers {
		req.Header.Set(k, v)
	}
	if cmd.BasicAuth != nil {
		req.SetBasicAuth(cmd.BasicAuth.User, cmd.BasicAuth.Password)
	}

	MockHTTPRequestsHistory = append(MockHTTPRequestsHistory, req)

	cmd.Response = &gohttp.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	return nil
}
