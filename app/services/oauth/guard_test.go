package oauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"
	"github.com/getfider/fider/app/services/oauth"
)

// End-to-end check that getOAuthRawProfile sends the profile request with the
// SSRF guard flag set whenever the resolved config is tenant-controlled.
// AllowPrivateNetworkTargets is enabled only so the local httptest token
// endpoint passes the URL preflight; the flag on the request is independent of
// it (netguard.ClientFor applies the escape hatch when the request runs).
func TestGetOAuthRawProfile_TenantConfigProfileRequestIsGuarded(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	env.Config.AllowPrivateNetworkTargets = true
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	// Built-in GitHub disabled at instance level: getConfig falls through to
	// GetCustomOAuthConfigByProvider, i.e. a tenant row with provider='github'.
	oauth.SetSystemProviderStatus(t, app.GitHubProvider, enum.OAuthConfigDisabled)

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"ACCESS_TOKEN","token_type":"bearer"}`))
	}))
	t.Cleanup(tokenServer.Close)

	for _, provider := range []string{"_custom", app.GitHubProvider} {
		httpclientmock.Reset()
		bus.Init(&oauth.Service{}, httpclientmock.Service{})
		bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
			q.Result = &entity.OAuthConfig{
				Provider:     q.Provider,
				Status:       enum.OAuthConfigEnabled,
				ClientID:     "CU_CL_ID",
				ClientSecret: "CU_CL_SECRET",
				AuthorizeURL: "https://example.org/oauth/authorize",
				TokenURL:     tokenServer.URL + "/token",
				ProfileURL:   "https://203.0.113.10/me",
			}
			return nil
		})

		ctx := newGetContext("http://login.test.fider.io:3000")
		rawProfile := &query.GetOAuthRawProfile{Provider: provider, Code: "CODE"}
		err := bus.Dispatch(ctx, rawProfile)
		Expect(err).IsNil()

		Expect(httpclientmock.CommandsHistory).HasLen(1)
		req := httpclientmock.CommandsHistory[0]
		Expect(req.URL).Equals("https://203.0.113.10/me")
		Expect(req.Headers["Authorization"]).Equals("Bearer ACCESS_TOKEN")
		if !req.BlockPrivateNetworkTargets {
			t.Errorf("profile request for %q is not guarded", provider)
		}
	}
}
