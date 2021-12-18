package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/handlers"
)

func TestManageBillingHandler_RedirectWhenUsingCNAME(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("https://feedback.demo.com/admin/billing").
		Execute(handlers.ManageBilling())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("https://demo.test.fider.io/admin/billing")
}

func TestManageBillingHandler_ReturnsCorrectBillingInformation(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetBillingState) error {
		trialEndsAt := time.Date(2021, time.February, 2, 4, 2, 2, 0, time.UTC)
		q.Result = &entity.BillingState{
			Status:         enum.BillingActive,
			PlanID:         "PLAN-123",
			SubscriptionID: "SUB-123",
			TrialEndsAt:    &trialEndsAt,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetBillingSubscription) error {
		Expect(q.SubscriptionID).Equals("SUB-123")

		q.Result = &entity.BillingSubscription{
			UpdateURL: "https://sandbox-subscription-management.paddle.com/subscription/SUB-123/hash/1111/update",
			CancelURL: "https://sandbox-subscription-management.paddle.com/subscription/SUB-123/hash/1111/cancel",
			PaymentInformation: entity.BillingPaymentInformation{
				PaymentMethod:  "card",
				CardType:       "visa",
				LastFourDigits: "1111",
				ExpiryDate:     "10/2031",
			},
			LastPayment: entity.BillingLastPayment{
				Amount:   float64(30),
				Currency: "USD",
				Date:     "2021-11-09",
			},
		}
		return nil
	})

	server := mock.NewServer()
	code, page := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("https://demo.test.fider.io/admin/billing").
		ExecuteAsPage(handlers.ManageBilling())

	Expect(code).Equals(http.StatusOK)
	Expect(page.Data).ContainsProps(dto.Props{
		"status":             float64(2),
		"trialEndsAt":        "2021-02-02T04:02:02Z",
		"subscriptionEndsAt": nil,
	})

	Expect(page.Data["subscription"]).ContainsProps(dto.Props{
		"updateURL": "https://sandbox-subscription-management.paddle.com/subscription/SUB-123/hash/1111/update",
		"cancelURL": "https://sandbox-subscription-management.paddle.com/subscription/SUB-123/hash/1111/cancel",
	})

	ExpectHandler(&query.GetBillingState{}).CalledOnce()
	ExpectHandler(&query.GetBillingSubscription{}).CalledOnce()
}

func TestGenerateCheckoutLinkHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.GenerateCheckoutLink) error {
		c.URL = "https://paddle.com/fake-checkout-url"
		return nil
	})

	server := mock.NewServer()
	code, json := server.
		WithURL("http://demo.test.fider.io/_api/billing/checkout-link").
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(handlers.GenerateCheckoutLink())

	Expect(code).Equals(http.StatusOK)
	Expect(json.String("url")).Equals("https://paddle.com/fake-checkout-url")
	ExpectHandler(&cmd.GenerateCheckoutLink{}).CalledOnce()
}
