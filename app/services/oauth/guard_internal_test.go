package oauth

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/netguard"
	"golang.org/x/oauth2"
)

func systemProvider(provider string) *entity.OAuthConfig {
	for _, p := range systemProviders {
		if p.Provider == provider {
			return p
		}
	}
	panic("unknown system provider " + provider)
}

func TestRequiresSSRFGuard(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		name     string
		provider string
		config   *entity.OAuthConfig
		guard    bool
	}{
		{"custom provider", "_custom", &entity.OAuthConfig{Provider: "_custom"}, true},
		{"built-in GitHub", app.GitHubProvider, systemProvider(app.GitHubProvider), false},
		{"built-in Google", app.GoogleProvider, systemProvider(app.GoogleProvider), false},
		{"built-in Facebook", app.FacebookProvider, systemProvider(app.FacebookProvider), false},
		// getConfig falls through to the tenant's custom config when a built-in
		// provider is not enabled at instance level: fail closed.
		{"tenant config named like a built-in", app.GitHubProvider, &entity.OAuthConfig{Provider: app.GitHubProvider}, true},
		{"unknown provider without prefix", "custom", &entity.OAuthConfig{Provider: "custom"}, true},
		{"nil config", app.GitHubProvider, nil, true},
		// A "_" key is guarded even if it somehow resolved to a system config.
		{"prefixed key", "_github", systemProvider(app.GitHubProvider), true},
	}

	for _, tc := range testCases {
		if got := requiresSSRFGuard(tc.provider, tc.config); got != tc.guard {
			t.Errorf("%s: requiresSSRFGuard(%q) = %v, expected %v", tc.name, tc.provider, got, tc.guard)
		}
	}
}

func TestTokenExchangeContext(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	testCases := []struct {
		guard        bool
		allowPrivate bool
		expected     *http.Client
	}{
		// Guarded exchanges use the guarded client...
		{true, false, netguard.Client},
		// ...unless the instance allows private network targets.
		{true, true, http.DefaultClient},
		// Unguarded (built-in) exchanges leave the context untouched, so
		// oauth2 uses its default client.
		{false, false, nil},
	}

	for _, tc := range testCases {
		env.Config.AllowPrivateNetworkTargets = tc.allowPrivate

		ctx := tokenExchangeContext(context.Background(), tc.guard)
		client, _ := ctx.Value(oauth2.HTTPClient).(*http.Client)
		if client != tc.expected {
			t.Errorf("tokenExchangeContext(guard=%v) with AllowPrivateNetworkTargets=%v returned unexpected client", tc.guard, tc.allowPrivate)
		}
	}
}

func TestNewProfileRequest(t *testing.T) {
	RegisterT(t)

	for _, guard := range []bool{true, false} {
		req := newProfileRequest(guard, "https://example.org/me", "TOKEN")
		Expect(req.URL).Equals("https://example.org/me")
		Expect(req.Method).Equals("GET")
		Expect(req.Headers["Authorization"]).Equals("Bearer TOKEN")
		Expect(req.BlockPrivateNetworkTargets).Equals(guard)
	}
}
