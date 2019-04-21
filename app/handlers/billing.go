package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// BillingPage is the billing settings page
func BillingPage() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Tenant().Billing.StripeCustomerID == "" {
			if err := bus.Dispatch(c, &cmd.CreateBillingCustomer{}); err != nil {
				return c.Failure(err)
			}

			if err := bus.Dispatch(c, &cmd.UpdateTenantBillingSettings{
				Settings: c.Tenant().Billing,
			}); err != nil {
				return c.Failure(err)
			}
		}

		getUpcomingInvoiceQuery := &query.GetUpcomingInvoice{}
		if c.Tenant().Billing.StripeSubscriptionID != "" {
			err := bus.Dispatch(c, getUpcomingInvoiceQuery)
			if err != nil {
				return c.Failure(err)
			}
		}

		paymentInfo := query.GetPaymentInfo{}
		err := bus.Dispatch(c, paymentInfo)
		if err != nil {
			return c.Failure(err)
		}

		listPlansQuery := &query.ListBillingPlans{}
		if paymentInfo.Result != nil {
			listPlansQuery.CountryCode = paymentInfo.Result.AddressCountry
			err = bus.Dispatch(c, listPlansQuery)
			if err != nil {
				println(err.Error())
				return c.Failure(err)
			}
		}

		countUsers := &query.CountUsers{}
		allCountries := &query.GetAllCountries{}
		if err := bus.Dispatch(c, countUsers); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Billing · Site Settings",
			ChunkName: "Billing.page",
			Data: web.Map{
				"invoiceDue":      getUpcomingInvoiceQuery.Result,
				"tenantUserCount": countUsers.Result,
				"plans":           listPlansQuery.Result,
				"paymentInfo":     paymentInfo.Result,
				"countries":       allCountries.Result,
			},
		})
	}
}

// UpdatePaymentInfo on stripe based on given input
func UpdatePaymentInfo() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateEditBillingPaymentInfo)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.UpdatePaymentInfo{Input: input.Model}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// GetBillingPlans returns a list of plans for given country code
func GetBillingPlans() web.HandlerFunc {
	return func(c *web.Context) error {
		countryCode := c.Param("countryCode")
		listPlansQuery := &query.ListBillingPlans{CountryCode: countryCode}
		err := bus.Dispatch(c, listPlansQuery)
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(listPlansQuery.Result)
	}
}

// BillingSubscribe subscribes current tenant to given plan on stripe
func BillingSubscribe() web.HandlerFunc {
	return func(c *web.Context) error {
		planID := c.Param("planID")

		paymentInfoQuery := &query.GetPaymentInfo{}
		err := bus.Dispatch(c, paymentInfoQuery)
		if err != nil {
			return c.Failure(err)
		}

		getPlanByIDQuery := &query.GetBillingPlanByID{
			PlanID:      planID,
			CountryCode: paymentInfoQuery.Result.AddressCountry,
		}
		err = bus.Dispatch(c, getPlanByIDQuery)
		if err != nil {
			return c.Failure(err)
		}
		plan := getPlanByIDQuery.Result

		countUsers := &query.CountUsers{}
		err = bus.Dispatch(c, countUsers)
		if err != nil {
			return c.Failure(err)
		}

		if plan.MaxUsers > 0 && countUsers.Result > plan.MaxUsers {
			return c.Unauthorized()
		}

		if err = bus.Dispatch(c, &cmd.CreateBillingSubscription{
			PlanID: plan.ID,
		}); err != nil {
			return c.Failure(err)
		}

		updateBilling := &cmd.UpdateTenantBillingSettings{Settings: c.Tenant().Billing}
		activateTenant := &cmd.ActivateTenant{TenantID: c.Tenant().ID}
		if err := bus.Dispatch(c, updateBilling, activateTenant); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// CancelBillingSubscription cancels current subscription from current tenant
func CancelBillingSubscription() web.HandlerFunc {
	return func(c *web.Context) error {
		err := bus.Dispatch(c, &cmd.CancelBillingSubscription{})
		if err != nil {
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.UpdateTenantBillingSettings{
			Settings: c.Tenant().Billing,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
