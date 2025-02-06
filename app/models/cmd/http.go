package cmd

import (
	"io"
	"net/http"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
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
	ResponseHeader     http.Header
}
