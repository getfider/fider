package webhooks

import (
	"encoding/json"
	"strconv"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/stripe/stripe-go/v83"
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

	// Get tenant ID from metadata
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

	return c.Ok(web.Map{})
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

	return c.Ok(web.Map{})
}
