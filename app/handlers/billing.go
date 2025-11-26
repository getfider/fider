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

// CreateStripeCheckoutSession creates a Stripe checkout session for new subscriptions
func CreateStripeCheckoutSession() web.HandlerFunc {
	return func(c *web.Context) error {
		stripe.Key = env.Config.Stripe.SecretKey

		returnURL := c.BaseURL() + "/admin/billing"
		tenantID := c.Tenant().ID

		params := &stripe.CheckoutSessionParams{
			Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(env.Config.Stripe.PriceID),
					Quantity: stripe.Int64(1),
				},
			},
			SuccessURL: stripe.String(returnURL + "?checkout=success"),
			CancelURL:  stripe.String(returnURL + "?checkout=cancelled"),
			Metadata: map[string]string{
				"tenant_id": fmt.Sprintf("%d", tenantID),
			},
			SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
				Metadata: map[string]string{
					"tenant_id": fmt.Sprintf("%d", tenantID),
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
}
