package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

// BillingPage is the billing settings page
func BillingPage() web.HandlerFunc {
	return func(c web.Context) error {
		err := ensureStripeCustomerID(c)
		if err != nil {
			return c.Failure(err)
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

func ensureStripeCustomerID(c web.Context) error {
	billing := c.Tenant().Billing
	if billing.StripeCustomerID == "" {
		_, err := c.Services().Billing.CreateCustomer("")
		if err != nil {
			return err
		}

		err = c.Services().Tenants.UpdateBillingSettings(billing)
		if err != nil {
			return err
		}
	}

	return nil
}
