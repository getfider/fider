package httpclientmock

import (
	"context"
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

type Service struct{}

func (s Service) Name() string {
	return "Mock"
}

func (s Service) Category() string {
	return "httpclient"
}

var RequestsHistory = make([]*http.Request, 0)

func (s Service) Enabled() bool {
	return env.IsTest()
}

func (s Service) Init() {
	RequestsHistory = make([]*http.Request, 0)
	bus.AddHandler(requestHandler)
}

func requestHandler(ctx context.Context, c *cmd.HTTPRequest) error {
	req, err := http.NewRequest(c.Method, c.URL, c.Body)
	if err != nil {
		return err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	if c.BasicAuth != nil {
		req.SetBasicAuth(c.BasicAuth.User, c.BasicAuth.Password)
	}

	RequestsHistory = append(RequestsHistory, req)

	c.ResponseStatusCode = http.StatusOK
	c.ResponseBody = []byte("")
	return nil
}
