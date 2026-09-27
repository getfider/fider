package webhook_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	pkgwebhook "github.com/getfider/fider/app/pkg/webhook"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"
	"github.com/getfider/fider/app/services/webhook"
)

// Webhook URLs are admin-configurable, so the outgoing request must opt in to
// the dial-time SSRF guard (the WebhookURL preflight alone is bypassable via
// DNS rebinding).
func TestTriggerWebhooks_BlocksPrivateNetworkTargets(t *testing.T) {
	RegisterT(t)
	bus.Init(webhook.Service{}, httpclientmock.Service{})

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveWebhooksByType) error {
		q.Result = []*entity.Webhook{
			{
				ID:         1,
				Name:       "Public webhook",
				Type:       enum.WebhookNewPost,
				Status:     enum.WebhookEnabled,
				Url:        "https://203.0.113.10/hook", // IP literal: no DNS lookup in tests
				Content:    "hello",
				HttpMethod: "POST",
			},
		}
		return nil
	})

	err := bus.Dispatch(context.Background(), &cmd.TriggerWebhooks{
		Type:  enum.WebhookNewPost,
		Props: pkgwebhook.Props{},
	})
	Expect(err).IsNil()

	Expect(httpclientmock.CommandsHistory).HasLen(1)
	Expect(httpclientmock.CommandsHistory[0].URL).Equals("https://203.0.113.10/hook")
	Expect(httpclientmock.CommandsHistory[0].BlockPrivateNetworkTargets).IsTrue()
}
