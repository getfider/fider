package oauth

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/netguard"
	"golang.org/x/oauth2"
)

func TestTokenExchangeContext(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	testCases := []struct {
		provider     string
		allowPrivate bool
		expected     *http.Client
	}{
		// Custom providers use the guarded client...
		{"_custom", false, netguard.Client},
		// ...unless the instance allows private network targets.
		{"_custom", true, http.DefaultClient},
		// Built-in providers leave the context untouched (oauth2 then uses its default).
		{app.GitHubProvider, false, nil},
		{app.GoogleProvider, false, nil},
		{app.FacebookProvider, false, nil},
	}

	for _, tc := range testCases {
		env.Config.AllowPrivateNetworkTargets = tc.allowPrivate

		ctx := tokenExchangeContext(context.Background(), tc.provider)
		client, _ := ctx.Value(oauth2.HTTPClient).(*http.Client)
		if client != tc.expected {
			t.Errorf("tokenExchangeContext(%q) with AllowPrivateNetworkTargets=%v returned unexpected client", tc.provider, tc.allowPrivate)
		}
	}
}

func TestNewProfileRequest_BlocksPrivateNetworkTargetsForCustomProviders(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		provider string
		block    bool
	}{
		{"_custom", true},
		{"_1234", true},
		{app.GitHubProvider, false},
		{app.GoogleProvider, false},
		{app.FacebookProvider, false},
	}

	for _, tc := range testCases {
		req := newProfileRequest(tc.provider, "https://example.org/me", "TOKEN")
		Expect(req.URL).Equals("https://example.org/me")
		Expect(req.Method).Equals("GET")
		Expect(req.Headers["Authorization"]).Equals("Bearer TOKEN")
		if req.BlockPrivateNetworkTargets != tc.block {
			t.Errorf("newProfileRequest(%q).BlockPrivateNetworkTargets = %v, expected %v", tc.provider, req.BlockPrivateNetworkTargets, tc.block)
		}
	}
}
