package httpclient

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

func init() {
	http.DefaultClient.Timeout = 30 * time.Second
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "HTTP"
}

func (s Service) Category() string {
	return "httpclient"
}

func (s Service) Enabled() bool {
	return !env.IsTest()
}

func (s Service) Init() {
	bus.AddHandler(requestHandler)
}

type BasicAuth struct {
	User     string
	Password string
}

type Request struct {
	URL       string
	Body      io.Reader
	Method    string
	Headers   map[string]string
	BasicAuth *BasicAuth

	ResponseBody       []byte
	ResponseStatusCode int
}

func requestHandler(ctx context.Context, cmd *Request) error {
	req, err := http.NewRequest(cmd.Method, cmd.URL, cmd.Body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

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

	defer res.Body.Close()
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cmd.ResponseBody = respBody
	cmd.ResponseStatusCode = res.StatusCode
	return nil
}
