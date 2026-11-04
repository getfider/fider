package webhooks_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers/webhooks"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/webhook"
)

const whSecret = "whsec_test"

// fakeStripe stands in for the Stripe API: subscriptions by ID, which customer owns them,
// and a record of subscriptions created.
type fakeStripe struct {
	subs    map[string]map[string]any
	created int
}

func (f *fakeStripe) sub(id, status, customer string) {
	f.subs[id] = map[string]any{
		"id": id, "object": "subscription", "status": status, "customer": customer,
		"metadata": map[string]string{"tenant_id": "42"},
	}
}

func (f *fakeStripe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet && path == "/v1/subscriptions":
		data := []any{}
		for _, s := range f.subs {
			if s["customer"] == r.URL.Query().Get("customer") {
				data = append(data, s)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"object": "list", "data": data, "has_more": false, "url": "/v1/subscriptions"})
	case r.Method == http.MethodPost && path == "/v1/subscriptions":
		f.created++
		_ = r.ParseForm()
		id := fmt.Sprintf("sub_new%d", f.created)
		f.sub(id, "active", r.PostForm.Get("customer"))
		_ = json.NewEncoder(w).Encode(f.subs[id])
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/subscriptions/"):
		s, ok := f.subs[strings.TrimPrefix(path, "/v1/subscriptions/")]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"type":"invalid_request_error","code":"resource_missing"}}`))
			return
		}
		_ = json.NewEncoder(w).Encode(s)
	case strings.HasPrefix(path, "/v1/customers/"):
		_, _ = w.Write([]byte(`{"id":"cus_new","object":"customer","address":{"country":"US"}}`))
	case strings.HasPrefix(path, "/v1/setup_intents/"):
		_, _ = w.Write([]byte(`{"id":"seti_1","object":"setup_intent","payment_method":"pm_1"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// setup wires a fake Stripe API, a stored subscription for tenant 42, and records activations/cancellations.
func setup(t *testing.T, storedSubID string) (*fakeStripe, *[]*cmd.ActivateStripeSubscription) {
	RegisterT(t)
	env.Config.Stripe.WebhookSecret = whSecret
	env.Config.Stripe.SecretKey = "sk_test_fake"
	env.Config.UserList.Enabled = false

	fake := &fakeStripe{subs: map[string]map[string]any{}}
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	stripe.SetBackend(stripe.APIBackend, stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
		URL:               stripe.String(srv.URL),
		MaxNetworkRetries: stripe.Int64(0),
		LeveledLogger:     &stripe.LeveledLogger{Level: stripe.LevelNull},
	}))
	t.Cleanup(func() { stripe.SetBackend(stripe.APIBackend, nil) })

	bus.AddHandler(func(ctx context.Context, q *query.GetStripeBillingState) error {
		tenant, _ := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
		Expect(tenant.ID).Equals(42)
		q.Result = &entity.StripeBillingState{SubscriptionID: storedSubID}
		return nil
	})
	activated := &[]*cmd.ActivateStripeSubscription{}
	bus.AddHandler(func(ctx context.Context, c *cmd.ActivateStripeSubscription) error {
		*activated = append(*activated, c)
		return nil
	})
	return fake, activated
}

func post(eventType string, object map[string]any, previous ...map[string]any) int {
	data := map[string]any{"object": object}
	if len(previous) > 0 {
		data["previous_attributes"] = previous[0]
	}
	body, _ := json.Marshal(map[string]any{
		"id": "evt_1", "object": "event", "api_version": stripe.APIVersion, "type": eventType, "data": data,
	})
	signed := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{Payload: body, Secret: whSecret})
	code, _ := mock.NewServer().
		AddHeader("Stripe-Signature", signed.Header).
		ExecutePost(webhooks.IncomingStripeWebhook(), string(body))
	return code
}

func setupCheckoutCompleted() (string, map[string]any) {
	return "checkout.session.completed", map[string]any{
		"id": "cs_1", "object": "checkout.session", "mode": "setup",
		"customer": "cus_new", "setup_intent": "seti_1",
		"metadata": map[string]string{"tenant_id": "42", "price_id": "price_1"},
	}
}

func TestStripe_IncompleteSubscriptionPaidLater_ActivatesTenant(t *testing.T) {
	fake, activated := setup(t, "")
	fake.sub("sub_1", "active", "cus_1")

	code := post("customer.subscription.updated", fake.subs["sub_1"], map[string]any{"status": "incomplete"})

	Expect(code).Equals(http.StatusOK)
	Expect(*activated).HasLen(1)
	Expect((*activated)[0].TenantID).Equals(42)
	Expect((*activated)[0].CustomerID).Equals("cus_1")
	Expect((*activated)[0].SubscriptionID).Equals("sub_1")
}

func TestStripe_IncompleteSubscriptionPaidLater_ActsOnCurrentState(t *testing.T) {
	fake, activated := setup(t, "")
	fake.sub("sub_1", "canceled", "cus_1") // a late event: the subscription has since ended

	code := post("customer.subscription.updated", map[string]any{"id": "sub_1", "status": "active"}, map[string]any{"status": "incomplete"})

	Expect(code).Equals(http.StatusOK)
	Expect(*activated).HasLen(0)
}

func TestStripe_OtherSubscriptionUpdates_Ignored(t *testing.T) {
	fake, activated := setup(t, "")
	fake.sub("sub_1", "active", "cus_1")

	code := post("customer.subscription.updated", fake.subs["sub_1"], map[string]any{"cancel_at_period_end": true})

	Expect(code).Equals(http.StatusOK)
	Expect(*activated).HasLen(0)
}

func TestStripe_SecondLiveSubscription_NotActivated(t *testing.T) {
	fake, activated := setup(t, "sub_current")
	fake.sub("sub_current", "active", "cus_0")
	fake.sub("sub_1", "active", "cus_1")

	code := post("customer.subscription.updated", fake.subs["sub_1"], map[string]any{"status": "incomplete"})

	Expect(code).Equals(http.StatusOK)
	Expect(*activated).HasLen(0)
}

func TestStripe_EndedStoredSubscription_Replaced(t *testing.T) {
	fake, activated := setup(t, "sub_old")
	fake.sub("sub_old", "canceled", "cus_0")
	fake.sub("sub_1", "active", "cus_1")

	code := post("customer.subscription.updated", fake.subs["sub_1"], map[string]any{"status": "incomplete"})

	Expect(code).Equals(http.StatusOK)
	Expect(*activated).HasLen(1)
	Expect((*activated)[0].SubscriptionID).Equals("sub_1")
}

func TestStripe_SetupCheckout_CreatesAndActivates(t *testing.T) {
	fake, activated := setup(t, "")

	code := post(setupCheckoutCompleted())

	Expect(code).Equals(http.StatusOK)
	Expect(fake.created).Equals(1)
	Expect(*activated).HasLen(1)
	Expect((*activated)[0].SubscriptionID).Equals("sub_new1")
}

func TestStripe_SetupCheckoutRedelivered_ReusesSubscription(t *testing.T) {
	// e.g. a Dashboard resend after the 24h idempotency key has expired
	fake, activated := setup(t, "")
	fake.sub("sub_existing", "active", "cus_new")

	code := post(setupCheckoutCompleted())

	Expect(code).Equals(http.StatusOK)
	Expect(fake.created).Equals(0)
	Expect(*activated).HasLen(1)
	Expect((*activated)[0].SubscriptionID).Equals("sub_existing")
}

func TestStripe_SetupCheckout_TenantAlreadySubscribed_NoNewSubscription(t *testing.T) {
	fake, activated := setup(t, "sub_current")
	fake.sub("sub_current", "active", "cus_0")

	code := post(setupCheckoutCompleted())

	Expect(code).Equals(http.StatusOK)
	Expect(fake.created).Equals(0)
	Expect(*activated).HasLen(0)
}
