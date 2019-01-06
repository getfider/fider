package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

// BillingPage is the billing settings page
func BillingPage() web.HandlerFunc {
	return func(c web.Context) error {
		_, err := c.Services().Billing.CreateCustomer("")
		if err != nil {
			return err
		}

		err = c.Services().Tenants.UpdateBillingSettings(c.Tenant().Billing)
		if err != nil {
			return err
		}

		paymentInfo, err := c.Services().Billing.GetPaymentInfo()
		if err != nil {
			return c.Failure(err)
		}

		plans, err := c.Services().Billing.ListPlans()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Billing · Site Settings",
			ChunkName: "Billing.page",
			Data: web.Map{
				"plans":       plans,
				"paymentInfo": paymentInfo,
				"countries":   models.GetAllCountries(),
			},
		})
	}
}

// UpdatePaymentInfo on stripe based on given input
func UpdatePaymentInfo() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateEditBillingPaymentInfo)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Billing.UpdatePaymentInfo(input.Model)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// BillingSubscribe subscribes current tenant to given plan on stripe
func BillingSubscribe() web.HandlerFunc {
	return func(c web.Context) error {
		planID := c.Param("planID")
		err := c.Services().Billing.Subscribe(planID)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.UpdateBillingSettings(c.Tenant().Billing)
		if err != nil {
			return err
		}
		err = c.Services().Tenants.Activate(c.Tenant().ID)
		if err != nil {
			return err
		}

		return c.Ok(web.Map{})
	}
}

// CancelBillingSubscription cancels current subscription from current tenant
func CancelBillingSubscription() web.HandlerFunc {
	return func(c web.Context) error {
		err := c.Services().Billing.CancelSubscription()
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.UpdateBillingSettings(c.Tenant().Billing)
		if err != nil {
			return err
		}

		return c.Ok(web.Map{})
	}
}
