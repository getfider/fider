package handlers

import (
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
