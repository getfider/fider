package http

import (
	"context"
	"net/http"
	gohttp "net/http"
)

type HTTPGetRequestCommand struct {
	URL       string
	Response  *http.Response
	Headers   map[string]string
	BasicAuth *BasicAuth
}

func HTTPGetRequest(ctx context.Context, cmd *HTTPGetRequestCommand) error {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	cmd.Response = res
	return nil
}
