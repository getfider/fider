package handlers

import (
	"fmt"
	"net/http"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/stripe/stripe-go/v83"
	portalsession "github.com/stripe/stripe-go/v83/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
)

// ManageBilling is the page used by administrators for Stripe billing settings
func ManageBilling() web.HandlerFunc {
	return func(c *web.Context) error {
		billingState := &query.GetStripeBillingState{}
		if err := bus.Dispatch(c, billingState); err != nil {
			return c.Failure(err)
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageBilling.page",
			Title: "Manage Billing · Site Settings",
			Data: web.Map{
				"stripeCustomerID":     billingState.Result.CustomerID,
				"stripeSubscriptionID": billingState.Result.SubscriptionID,
				"paddleSubscriptionID": billingState.Result.PaddleSubscriptionID,
				"isPro":                c.Tenant().IsPro,
			},
		})
	}
}

// CreateStripePortalSession creates a Stripe customer portal session
func CreateStripePortalSession() web.HandlerFunc {
	return func(c *web.Context) error {
		billingState := &query.GetStripeBillingState{}
		if err := bus.Dispatch(c, billingState); err != nil {
			return c.Failure(err)
		}

		if billingState.Result.CustomerID == "" {
			return c.BadRequest(web.Map{"message": "No Stripe customer found"})
		}

		stripe.Key = env.Config.Stripe.SecretKey

		returnURL := c.BaseURL() + "/admin/billing"

		params := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(billingState.Result.CustomerID),
			ReturnURL: stripe.String(returnURL),
		}

		s, err := portalsession.New(params)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"url": s.URL,
		})
	}
}

// createCheckoutSession creates a Stripe checkout session for the given price ID.
// Uses setup mode so we can collect a billing address before the first invoice is
// generated — the subscription itself is created in the webhook handler, where we
// can attach the UK VAT tax rate if the customer's address country is GB.
func createCheckoutSession(c *web.Context, priceID string) error {
	stripe.Key = env.Config.Stripe.SecretKey

	returnURL := c.BaseURL() + "/admin/billing"
	tenantIDStr := fmt.Sprintf("%d", c.Tenant().ID)

	metadata := map[string]string{
		"tenant_id": tenantIDStr,
		"price_id":  priceID,
	}

	params := &stripe.CheckoutSessionParams{
		Mode:                     stripe.String(string(stripe.CheckoutSessionModeSetup)),
		PaymentMethodTypes:       stripe.StringSlice([]string{"card"}),
		BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionRequired)),
		CustomerCreation:         stripe.String(string(stripe.CheckoutSessionCustomerCreationAlways)),
		SuccessURL:               stripe.String(returnURL + "?checkout=success"),
		CancelURL:                stripe.String(returnURL + "?checkout=cancelled"),
		Metadata:                 metadata,
		SetupIntentData: &stripe.CheckoutSessionSetupIntentDataParams{
			Metadata: metadata,
		},
		CustomText: &stripe.CheckoutSessionCustomTextParams{
			Submit: &stripe.CheckoutSessionCustomTextSubmitParams{
				Message: stripe.String("By submitting, you'll be subscribed to Pro at the price shown on the previous page. Your subscription starts immediately."),
			},
		},
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{
		"url": s.URL,
	})
}

// CreateStripeCheckoutSession creates a Stripe checkout session for new subscriptions
func CreateStripeCheckoutSession() web.HandlerFunc {
	return func(c *web.Context) error {
		return createCheckoutSession(c, env.Config.Stripe.PriceID)
	}
}

// CreateStripeAnnualCheckoutSession creates a Stripe checkout session for annual subscriptions
func CreateStripeAnnualCheckoutSession() web.HandlerFunc {
	return func(c *web.Context) error {
		return createCheckoutSession(c, env.Config.Stripe.AnnualPriceID)
	}
}
