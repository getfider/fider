package handlers

import (
	"fmt"
	"net/http"

	"github.com/Spicy-Bush/fider-tarkov-community/app/actions"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/env"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
)

// ManageBilling is the page used by administrators for billing settings
func ManageBilling() web.HandlerFunc {
	return func(c *web.Context) error {

		// It's not possible to use custom domains on billing page, so redirect to Fider url
		if c.Request.IsCustomDomain() {
			url := fmt.Sprintf("https://%s.%s/admin/billing", c.Tenant().Subdomain, env.Config.HostDomain)
			return c.Redirect(url)
		}

		billingState := &query.GetBillingState{}
		if err := bus.Dispatch(c, billingState); err != nil {
			return c.Failure(err)
		}

		billingSubscription := &query.GetBillingSubscription{
			SubscriptionID: billingState.Result.SubscriptionID,
		}
		if billingState.Result.Status == enum.BillingActive {
			if err := bus.Dispatch(c, billingSubscription); err != nil {
				return c.Failure(err)
			}
		}

		return c.Page(http.StatusOK, web.Props{
			Page:  "Administration/pages/ManageBilling.page",
			Title: "Manage Billing · Site Settings",
			Data: web.Map{
				"paddle": web.Map{
					"isSandbox":     env.Config.Paddle.IsSandbox,
					"vendorId":      env.Config.Paddle.VendorID,
					"monthlyPlanId": env.Config.Paddle.MonthlyPlanID,
					"yearlyPlanId":  env.Config.Paddle.YearlyPlanID,
				},
				"status":             billingState.Result.Status,
				"trialEndsAt":        billingState.Result.TrialEndsAt,
				"subscriptionEndsAt": billingState.Result.SubscriptionEndsAt,
				"subscription":       billingSubscription.Result,
			},
		})
	}
}

// GenerateCheckoutLink generates a Paddle-hosted checkout link for the service subscription
func GenerateCheckoutLink() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.GenerateCheckoutLink)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		generateLink := &cmd.GenerateCheckoutLink{
			PlanID: action.PlanID,
			Passthrough: dto.PaddlePassthrough{
				TenantID: c.Tenant().ID,
			},
		}

		if err := bus.Dispatch(c, generateLink); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"url": generateLink.URL,
		})
	}
}
