package httpclientmock

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

type Service struct{}

var RequestsHistory = make([]*http.Request, 0)

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	RequestsHistory = make([]*http.Request, 0)
	bus.AddHandler(s, requestHandler)
}

func requestHandler(ctx context.Context, cmd *httpclient.Request) error {
	req, err := http.NewRequest(cmd.Method, cmd.URL, cmd.Body)
	if err != nil {
		return err
	}

	for k, v := range cmd.Headers {
		req.Header.Set(k, v)
	}
	if cmd.BasicAuth != nil {
		req.SetBasicAuth(cmd.BasicAuth.User, cmd.BasicAuth.Password)
	}

	RequestsHistory = append(RequestsHistory, req)

	cmd.Response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	return nil
}
