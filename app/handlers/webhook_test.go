package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

// The test webhook route is a POST (it fires an outbound request), so it must work
// when called the way the frontend's http.post does it: POST, JSON Content-Type,
// no body, through the CSRF middleware.
func TestTestWebhookHandler_POST(t *testing.T) {
	RegisterT(t)

	var triggeredID int
	bus.AddHandler(func(ctx context.Context, c *cmd.TestWebhook) error {
		triggeredID = c.ID
		c.Result = &dto.WebhookTriggerResult{Success: true, StatusCode: http.StatusOK}
		return nil
	})

	server := mock.NewServer()
	server.Use(middlewares.CSRF())

	code, query := server.
		WithURL("http://demo.test.fider.io/_api/admin/webhook/test/42").
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", 42).
		ExecutePostAsJSON(handlers.TestWebhook(), "")

	Expect(code).Equals(http.StatusOK)
	Expect(triggeredID).Equals(42)
	Expect(string(query.Raw("success"))).Equals("true")
	Expect(query.Int32("status_code")).Equals(http.StatusOK)
}
