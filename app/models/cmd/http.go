package cmd

import (
	"io"

	"github.com/getfider/fider/app/models/dto"
)

type HTTPRequest struct {
	URL       string
	Body      io.Reader
	Method    string
	Headers   map[string]string
	BasicAuth *dto.BasicAuth

	//Output
	ResponseBody       []byte
	ResponseStatusCode int
}
