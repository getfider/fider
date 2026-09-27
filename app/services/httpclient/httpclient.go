package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/netguard"
)

func init() {
	http.DefaultClient = &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
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

// clientFor returns the HTTP client to use for the given request. Requests that
// opt in to BlockPrivateNetworkTargets use the netguard client, which rejects
// private/internal addresses at dial time, unless the instance explicitly
// allows private network targets.
func clientFor(c *cmd.HTTPRequest) *http.Client {
	if c.BlockPrivateNetworkTargets && !env.Config.AllowPrivateNetworkTargets {
		return netguard.Client
	}
	return http.DefaultClient
}

func requestHandler(ctx context.Context, c *cmd.HTTPRequest) error {
	req, err := http.NewRequest(c.Method, c.URL, c.Body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	if c.BasicAuth != nil {
		req.SetBasicAuth(c.BasicAuth.User, c.BasicAuth.Password)
	}

	res, err := clientFor(c).Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	respBody, err := io.ReadAll(io.LimitReader(res.Body, 1<<20)) // 1 MB limit
	if err != nil {
		return err
	}

	c.ResponseBody = respBody
	c.ResponseStatusCode = res.StatusCode
	c.ResponseHeader = res.Header
	return nil
}
