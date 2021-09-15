package handlers

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ManageBilling is the page used by administrators for billing settings
func ManageBilling() web.HandlerFunc {
	return func(c *web.Context) error {
		getState := &query.GetBillingState{}
		if err := bus.Dispatch(c, getState); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Manage Billing · Site Settings",
			ChunkName: "ManageBilling.page",
			Data: web.Map{
				"status":             getState.Result.Status,
				"trialEndsAt":        getState.Result.TrialEndsAt,
				"subscriptionEndsAt": getState.Result.SubscriptionEndsAt,
			},
		})
	}
}

// GenerateCheckoutLink generates a Paddle-hosted checkout link for the service subscription
func GenerateCheckoutLink() web.HandlerFunc {
	return func(c *web.Context) error {
		generateLink := &cmd.GenerateCheckoutLink{
			Email:     c.User().Email,
			ReturnURL: c.BaseURL() + "/admin/billing",
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
