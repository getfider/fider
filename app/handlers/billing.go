package handlers

import (
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// BillingPage is the billing settings page
func BillingPage() web.HandlerFunc {
	return func(c web.Context) error {
		if !env.IsBillingEnabled() || c.Tenant().Billing == nil {
			return c.Redirect(c.BaseURL())
		}

		return c.Page(web.Props{
			Title:     "Billing · Site Settings",
			ChunkName: "Billing.page",
		})
	}
}
