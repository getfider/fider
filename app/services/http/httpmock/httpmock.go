package httpmock

import (
	"bytes"
	"context"
	"io/ioutil"
	gohttp "net/http"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/http"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := gohttp.DefaultTransport.(*gohttp.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

type ClientService struct{}

var RequestsHistory = make([]*gohttp.Request, 0)

func (s ClientService) Enabled() bool {
	return true
}

func (s ClientService) Init() {
	RequestsHistory = make([]*gohttp.Request, 0)
	bus.AddHandler(s, requestHandler)
}

func requestHandler(ctx context.Context, cmd *http.Request) error {
	req, err := gohttp.NewRequest(cmd.Method, cmd.URL, cmd.Body)
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

	cmd.Response = &gohttp.Response{
		StatusCode: gohttp.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	return nil
}
