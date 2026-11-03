package cmd

import (
	"io"
	"net/http"

	"github.com/getfider/fider/app/models/dto"
)

type HTTPRequest struct {
	URL       string
	Body      io.Reader
	Method    string
	Headers   map[string]string
	BasicAuth *dto.BasicAuth

	// BlockPrivateNetworkTargets makes the request refuse to connect to
	// private/internal network addresses (see netguard.IsBlockedIP). The check
	// runs at dial time on the resolved IP, so it is not bypassable via DNS
	// rebinding. Set it for requests to user-configurable URLs (SSRF risk).
	// Ignored when env.Config.AllowPrivateNetworkTargets is true. Guarded
	// requests do not use HTTP(S)_PROXY.
	BlockPrivateNetworkTargets bool

	//Output
	ResponseBody       []byte
	ResponseStatusCode int
	ResponseHeader     http.Header
}
