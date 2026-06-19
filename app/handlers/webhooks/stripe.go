package webhooks

import (
	"encoding/json"
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/setupintent"
	stripesub "github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/webhook"
)

// IncomingStripeWebhook handles all incoming requests from Stripe Webhooks
func IncomingStripeWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		payload := []byte(c.Request.Body)
		sigHeader := c.Request.GetHeader("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, sigHeader, env.Config.Stripe.WebhookSecret)
		if err != nil {
			return c.Failure(errors.Wrap(err, "failed to verify stripe webhook signature"))
		}

		switch event.Type {
		case "checkout.session.completed":
			return handleCheckoutSessionCompleted(c, event)
		case "customer.subscription.deleted":
			return handleSubscriptionDeleted(c, event)
		default:
			log.Debugf(c, "Ignoring Stripe webhook event: '@{EventType}'", dto.Props{
				"EventType": event.Type,
			})
			return c.Ok(web.Map{})
		}
	}
}

func handleCheckoutSessionCompleted(c *web.Context, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse checkout session"))
	}

	if session.Mode == stripe.CheckoutSessionModeSetup {
		return handleSetupCheckoutCompleted(c, &session)
	}
	return handleSubscriptionCheckoutCompleted(c, &session)
}

// handleSubscriptionCheckoutCompleted handles the legacy mode=subscription flow.
// Kept to drain any in-flight sessions created before the switch to setup mode.
func handleSubscriptionCheckoutCompleted(c *web.Context, session *stripe.CheckoutSession) error {
	tenantIDStr, ok := session.Metadata["tenant_id"]
	if !ok {
		return c.Failure(errors.New("tenant_id not found in session metadata"))
	}

	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse tenant_id"))
	}

	activate := &cmd.ActivateStripeSubscription{
		TenantID:       tenantID,
		CustomerID:     session.Customer.ID,
		SubscriptionID: session.Subscription.ID,
	}

	if err := bus.Dispatch(c, activate); err != nil {
		return c.Failure(err)
	}

	log.Infof(c, "Stripe subscription activated for tenant @{TenantID}", dto.Props{
		"TenantID": tenantID,
	})

	updateUserListPlan(c, tenantID, enum.PlanPro)

	return c.Ok(web.Map{})
}

// handleSetupCheckoutCompleted creates the Stripe Subscription server-side once the
// customer has supplied a payment method and billing address via Checkout. The address
// determines whether the UK VAT tax rate is attached, so VAT lands on the first invoice.
func handleSetupCheckoutCompleted(c *web.Context, session *stripe.CheckoutSession) error {
	if session.Customer == nil {
		return c.Failure(errors.New("setup-mode checkout session has no customer"))
	}
	if session.SetupIntent == nil {
		return c.Failure(errors.New("setup-mode checkout session has no setup intent"))
	}

	tenantIDStr, ok := session.Metadata["tenant_id"]
	if !ok {
		return c.Failure(errors.New("tenant_id not found in session metadata"))
	}
	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse tenant_id"))
	}

	priceID, ok := session.Metadata["price_id"]
	if !ok {
		return c.Failure(errors.New("price_id not found in session metadata"))
	}

	stripe.Key = env.Config.Stripe.SecretKey

	// In setup mode the address lives on the Customer, not on session.CustomerDetails.
	cus, err := customer.Get(session.Customer.ID, nil)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to fetch customer"))
	}
	country := ""
	if cus.Address != nil {
		country = cus.Address.Country
	}

	si, err := setupintent.Get(session.SetupIntent.ID, nil)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to fetch setup intent"))
	}
	if si.PaymentMethod == nil {
		return c.Failure(errors.New("setup intent has no payment method"))
	}

	subParams := &stripe.SubscriptionParams{
		Customer:             stripe.String(session.Customer.ID),
		DefaultPaymentMethod: stripe.String(si.PaymentMethod.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(priceID)},
		},
		Metadata: map[string]string{
			"tenant_id": tenantIDStr,
		},
	}
	if country == "GB" {
		if env.Config.Stripe.UKVATTaxRateID == "" {
			return c.Failure(errors.New("UK customer detected but STRIPE_UK_VAT_TAX_RATE_ID is not configured"))
		}
		subParams.DefaultTaxRates = stripe.StringSlice([]string{env.Config.Stripe.UKVATTaxRateID})
	}
	// session.ID is stable per Checkout session, so a retried webhook returns the same subscription.
	subParams.SetIdempotencyKey(session.ID)

	newSub, err := stripesub.New(subParams)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to create subscription"))
	}

	// Only activate the tenant if Stripe successfully collected payment on the first invoice.
	// allow_incomplete (the default) leaves the sub in `incomplete` if the charge failed (e.g. SCA
	// challenge on off-session payment) — in that case Stripe will keep retrying and eventually
	// fire subscription.deleted, which the existing handler already handles.
	if newSub.Status != stripe.SubscriptionStatusActive && newSub.Status != stripe.SubscriptionStatusTrialing {
		log.Warnf(c, "Stripe subscription created in non-active state '@{Status}' for tenant @{TenantID}", dto.Props{
			"Status":         string(newSub.Status),
			"TenantID":       tenantID,
			"SubscriptionID": newSub.ID,
		})
		return c.Ok(web.Map{})
	}

	activate := &cmd.ActivateStripeSubscription{
		TenantID:       tenantID,
		CustomerID:     session.Customer.ID,
		SubscriptionID: newSub.ID,
	}
	if err := bus.Dispatch(c, activate); err != nil {
		return c.Failure(err)
	}

	log.Infof(c, "Stripe subscription activated for tenant @{TenantID}", dto.Props{
		"TenantID": tenantID,
	})

	updateUserListPlan(c, tenantID, enum.PlanPro)

	return c.Ok(web.Map{})
}

func updateUserListPlan(c *web.Context, tenantID int, plan enum.Plan) {
	if !env.Config.UserList.Enabled {
		return
	}
	getTenant := &query.GetTenantByDomain{Domain: strconv.Itoa(tenantID)}
	if err := bus.Dispatch(c, getTenant); err == nil && getTenant.Result != nil {
		c.Enqueue(tasks.UserListUpdateCompany(&dto.UserListUpdateCompany{
			TenantID: tenantID,
			Name:     getTenant.Result.Name,
			Plan:     plan,
		}))
	}
}

func handleSubscriptionDeleted(c *web.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse subscription"))
	}

	// Get tenant ID from metadata
	tenantIDStr, ok := subscription.Metadata["tenant_id"]
	if !ok {
		log.Warnf(c, "tenant_id not found in subscription metadata for subscription @{SubscriptionID}", dto.Props{
			"SubscriptionID": subscription.ID,
		})
		return c.Ok(web.Map{})
	}

	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse tenant_id"))
	}

	cancel := &cmd.CancelStripeSubscription{
		TenantID: tenantID,
	}

	if err := bus.Dispatch(c, cancel); err != nil {
		return c.Failure(err)
	}

	log.Infof(c, "Stripe subscription cancelled for tenant @{TenantID}", dto.Props{
		"TenantID": tenantID,
	})

	updateUserListPlan(c, tenantID, enum.PlanFree)

	return c.Ok(web.Map{})
}
