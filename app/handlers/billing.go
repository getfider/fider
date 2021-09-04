package handlers

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// ManageBilling is the page used by administrators for billing settings
func ManageBilling() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:     "Manage Billing · Site Settings",
			ChunkName: "ManageBilling.page",
		})
	}
}

// GenerateCheckoutLink generates a Paddle-hosted checkout link for the service subscription
func GenerateCheckoutLink() web.HandlerFunc {
	return func(c *web.Context) error {
		generateLink := &cmd.GenerateCheckoutLink{}
		if err := bus.Dispatch(c, generateLink); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"url": generateLink.URL,
		})
	}
}
