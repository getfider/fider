package http

import (
	"context"
	"io"
	"net/http"
	gohttp "net/http"

	"github.com/getfider/fider/app/pkg/bus"
)

func init() {
	bus.Register(&HTTPClientService{})
}

type HTTPClientService struct{}

func (s HTTPClientService) Enabled() bool {
	return true
}

func (s HTTPClientService) Init() {
	bus.AddHandler(s, requestHandler)
}

type BasicAuth struct {
	User     string
	Password string
}

type Request struct {
	URL       string
	Body      io.Reader
	Method    string
	Response  *http.Response
	Headers   map[string]string
	BasicAuth *BasicAuth
}

func requestHandler(ctx context.Context, cmd *Request) error {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	cmd.Response = res
	return nil
}
