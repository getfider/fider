package middlewares

import (
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

// BlockFreemiumBillingAccess blocks access to billing routes for free users when freemium is enabled
func BlockFreemiumBillingAccess() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			// Only check if freemium is enabled
			if env.IsFreemium() {
				// Get the tenant's billing state
				billingState := &query.GetBillingState{}
				if err := bus.Dispatch(c, billingState); err != nil {
					return c.Failure(err)
				}

				// If the tenant is on the free_forever plan, block access
				if billingState.Result.Status == enum.BillingFreeForever {
					return c.Forbidden()
				}
			}

			return next(c)
		}
	}
}
