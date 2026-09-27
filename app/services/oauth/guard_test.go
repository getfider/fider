package oauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"
	"github.com/getfider/fider/app/services/oauth"
)

// End-to-end check that getOAuthRawProfile sends the custom provider's profile
// request with the SSRF guard flag set. AllowPrivateNetworkTargets is enabled
// only so the local httptest token endpoint passes the URL preflight; the flag
// on the request is independent of it (netguard.ClientFor applies the escape
// hatch when the request is executed).
func TestGetOAuthRawProfile_CustomProviderProfileRequestIsGuarded(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	env.Config.AllowPrivateNetworkTargets = true
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"ACCESS_TOKEN","token_type":"bearer"}`))
	}))
	t.Cleanup(tokenServer.Close)

	bus.Init(&oauth.Service{}, httpclientmock.Service{})
	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		q.Result = &entity.OAuthConfig{
			Provider:     q.Provider,
			ClientID:     "CU_CL_ID",
			ClientSecret: "CU_CL_SECRET",
			AuthorizeURL: "https://example.org/oauth/authorize",
			TokenURL:     tokenServer.URL + "/token",
			ProfileURL:   "https://203.0.113.10/me",
		}
		return nil
	})

	ctx := newGetContext("http://login.test.fider.io:3000")
	rawProfile := &query.GetOAuthRawProfile{Provider: "_custom", Code: "CODE"}
	err := bus.Dispatch(ctx, rawProfile)
	Expect(err).IsNil()

	Expect(httpclientmock.CommandsHistory).HasLen(1)
	req := httpclientmock.CommandsHistory[0]
	Expect(req.URL).Equals("https://203.0.113.10/me")
	Expect(req.Headers["Authorization"]).Equals("Bearer ACCESS_TOKEN")
	Expect(req.BlockPrivateNetworkTargets).IsTrue()
}
