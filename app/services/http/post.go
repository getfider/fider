package http

import (
	"context"
	"io"
	"net/http"
	gohttp "net/http"
)

type HTTPPostRequestCommand struct {
	URL       string
	Body      io.Reader
	Response  *http.Response
	Headers   map[string]string
	BasicAuth *BasicAuth
}

func HTTPPostRequest(ctx context.Context, cmd *HTTPPostRequestCommand) error {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	cmd.Response = res
	return nil
}
