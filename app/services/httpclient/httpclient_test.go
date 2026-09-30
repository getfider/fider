package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/netguard"
)

// The httpclient service is disabled under env.IsTest(), so these tests call
// requestHandler directly instead of dispatching through the bus.
func TestRequestHandler_BlockPrivateNetworkTargets(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("internal secret"))
	}))
	t.Cleanup(server.Close)

	testCases := []struct {
		name         string
		block        bool
		allowPrivate bool
		expectBlock  bool
	}{
		{"guarded request to loopback is rejected", true, false, true},
		{"unguarded request to loopback succeeds", false, false, false},
		{"guarded request succeeds when private targets are allowed", true, true, false},
	}

	for _, tc := range testCases {
		env.Config.AllowPrivateNetworkTargets = tc.allowPrivate

		req := &cmd.HTTPRequest{
			URL:                        server.URL,
			Method:                     "GET",
			BlockPrivateNetworkTargets: tc.block,
		}
		err := requestHandler(context.Background(), req)

		if tc.expectBlock {
			if err == nil || !errors.Is(err, netguard.ErrBlockedAddress) {
				t.Errorf("%s: expected ErrBlockedAddress, got %v", tc.name, err)
			}
			Expect(req.ResponseStatusCode).Equals(0)
			continue
		}

		if err != nil {
			t.Errorf("%s: unexpected error %v", tc.name, err)
			continue
		}
		Expect(req.ResponseStatusCode).Equals(http.StatusOK)
		Expect(string(req.ResponseBody)).Equals("internal secret")
	}
}
