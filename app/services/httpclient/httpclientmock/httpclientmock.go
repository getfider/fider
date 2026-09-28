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

// CommandsHistory records every dispatched cmd.HTTPRequest, so tests can assert
// on fields that are not part of the http.Request (e.g. BlockPrivateNetworkTargets).
var CommandsHistory = make([]*cmd.HTTPRequest, 0)

func (s Service) Enabled() bool {
	return env.IsTest()
}

// Reset clears the recorded request history. It runs on every Init (i.e. each
// bus.Init that includes this service); tests that assert on the history also
// call it directly so nothing can leak in from earlier tests.
func Reset() {
	RequestsHistory = make([]*http.Request, 0)
	CommandsHistory = make([]*cmd.HTTPRequest, 0)
}

func (s Service) Init() {
	Reset()
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
	CommandsHistory = append(CommandsHistory, c)

	c.ResponseStatusCode = http.StatusOK
	c.ResponseBody = []byte("")
	return nil
}
