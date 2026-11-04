package webhooks

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
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
		case "customer.subscription.updated":
			return handleSubscriptionUpdated(c, event)
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

	stripe.Key = env.Config.Stripe.SecretKey
	return activateTenant(c, tenantID, session.Customer.ID, session.Subscription.ID)
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

	// The idempotency key is only kept for 24h, but Stripe retries webhooks for 3 days and the
	// Dashboard can resend them for 15. CustomerCreation=always gives every Checkout session its
	// own customer, so any subscription on it was created by an earlier delivery of this event.
	listParams := &stripe.SubscriptionListParams{Customer: stripe.String(session.Customer.ID), Status: stripe.String("all")}
	listParams.Limit = stripe.Int64(1)
	iter := stripesub.List(listParams)
	var newSub *stripe.Subscription
	if iter.Next() {
		newSub = iter.Subscription()
	}
	if err := iter.Err(); err != nil {
		return c.Failure(errors.Wrap(err, "failed to list customer subscriptions"))
	}

	if newSub == nil {
		// The session only saved a card, so nothing has been charged yet. Don't start a second
		// subscription for a tenant that is already paying (e.g. two Checkout sessions left open).
		liveSubID, err := liveSubscriptionID(c, tenantID)
		if err != nil {
			return c.Failure(err)
		}
		if liveSubID != "" {
			log.Errorf(c, "Not creating Stripe subscription for tenant @{TenantID}: it already has live subscription @{SubscriptionID} (checkout customer @{CustomerID})", dto.Props{
				"TenantID":       tenantID,
				"SubscriptionID": liveSubID,
				"CustomerID":     session.Customer.ID,
			})
			return c.Ok(web.Map{})
		}

		newSub, err = stripesub.New(subParams)
		if err != nil {
			return c.Failure(errors.Wrap(err, "failed to create subscription"))
		}
	}

	// Only activate the tenant if Stripe collected payment on the first invoice. Otherwise
	// (declined, or 3-D Secure needed off-session) the subscription stays `incomplete`: Stripe
	// does not retry that first charge, but if the customer pays the open invoice within 23h it
	// becomes `active` and handleSubscriptionUpdated activates the tenant. If not, it expires.
	if !isLive(newSub.Status) {
		log.Warnf(c, "Stripe subscription @{SubscriptionID} is '@{Status}' for tenant @{TenantID}", dto.Props{
			"Status":         string(newSub.Status),
			"TenantID":       tenantID,
			"SubscriptionID": newSub.ID,
		})
		return c.Ok(web.Map{})
	}

	return activateTenant(c, tenantID, session.Customer.ID, newSub.ID)
}

// handleSubscriptionUpdated activates the tenant when a subscription that was created
// `incomplete` gets its first invoice paid (see handleSetupCheckoutCompleted).
func handleSubscriptionUpdated(c *web.Context, event stripe.Event) error {
	if event.GetPreviousValue("status") != string(stripe.SubscriptionStatusIncomplete) {
		return c.Ok(web.Map{})
	}

	// Events can arrive late or out of order, so act on the subscription's current state.
	stripe.Key = env.Config.Stripe.SecretKey
	sub, err := stripesub.Get(event.GetObjectValue("id"), nil)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to fetch subscription"))
	}
	if !isLive(sub.Status) || sub.Customer == nil {
		return c.Ok(web.Map{})
	}

	tenantIDStr, ok := sub.Metadata["tenant_id"]
	if !ok {
		log.Warnf(c, "tenant_id not found in subscription metadata for subscription @{SubscriptionID}", dto.Props{
			"SubscriptionID": sub.ID,
		})
		return c.Ok(web.Map{})
	}
	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		return c.Failure(errors.Wrap(err, "failed to parse tenant_id"))
	}

	return activateTenant(c, tenantID, sub.Customer.ID, sub.ID)
}

func isLive(status stripe.SubscriptionStatus) bool {
	return status == stripe.SubscriptionStatusActive || status == stripe.SubscriptionStatusTrialing
}

// liveSubscriptionID returns the tenant's stored Stripe subscription if Stripe still considers it
// live (anything but ended), or "" if there is none. Expects stripe.Key to be set.
func liveSubscriptionID(c *web.Context, tenantID int) (string, error) {
	tenantCtx := context.WithValue(c, app.TenantCtxKey, &entity.Tenant{ID: tenantID})
	billing := &query.GetStripeBillingState{}
	if err := bus.Dispatch(tenantCtx, billing); err != nil {
		return "", errors.Wrap(err, "failed to read stripe billing state")
	}
	subID := billing.Result.SubscriptionID
	if subID == "" {
		return "", nil
	}

	sub, err := stripesub.Get(subID, nil)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok && stripeErr.Code == stripe.ErrorCodeResourceMissing {
			return "", nil
		}
		return "", errors.Wrap(err, "failed to fetch stored subscription")
	}
	if sub.Status == stripe.SubscriptionStatusCanceled || sub.Status == stripe.SubscriptionStatusIncompleteExpired {
		return "", nil
	}
	return subID, nil
}

// activateTenant records the subscription against the tenant and upgrades it to Pro, unless the
// tenant already has a different live subscription: overwriting it would orphan that one (still
// billing, but unreachable from "Manage billing"), so that case is logged for manual cleanup.
// Expects stripe.Key to be set.
func activateTenant(c *web.Context, tenantID int, customerID, subscriptionID string) error {
	liveSubID, err := liveSubscriptionID(c, tenantID)
	if err != nil {
		return c.Failure(err)
	}
	if liveSubID != "" && liveSubID != subscriptionID {
		log.Errorf(c, "Duplicate Stripe subscription @{NewSubscriptionID} (customer @{CustomerID}) for tenant @{TenantID}, which already has @{SubscriptionID}: not activated, cancel/refund it manually", dto.Props{
			"TenantID":          tenantID,
			"SubscriptionID":    liveSubID,
			"NewSubscriptionID": subscriptionID,
			"CustomerID":        customerID,
		})
		return c.Ok(web.Map{})
	}

	activate := &cmd.ActivateStripeSubscription{
		TenantID:       tenantID,
		CustomerID:     customerID,
		SubscriptionID: subscriptionID,
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
		TenantID:       tenantID,
		SubscriptionID: subscription.ID,
	}

	if err := bus.Dispatch(c, cancel); err != nil {
		return c.Failure(err)
	}

	if !cancel.Result {
		log.Warnf(c, "Ignoring deletion of Stripe subscription @{SubscriptionID}: not the current subscription of tenant @{TenantID}", dto.Props{
			"SubscriptionID": subscription.ID,
			"TenantID":       tenantID,
		})
		return c.Ok(web.Map{})
	}

	log.Infof(c, "Stripe subscription cancelled for tenant @{TenantID}", dto.Props{
		"TenantID": tenantID,
	})

	updateUserListPlan(c, tenantID, enum.PlanFree)

	return c.Ok(web.Map{})
}
