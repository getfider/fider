package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

// BillingPage is the billing settings page
func BillingPage() web.HandlerFunc {
	return func(c web.Context) error {
		var err error
		if c.Tenant().Billing.StripeCustomerID == "" {
			_, err = c.Services().Billing.CreateCustomer("")
			if err != nil {
				return err
			}

			err = c.Services().Tenants.UpdateBillingSettings(c.Tenant().Billing)
			if err != nil {
				return err
			}
		}

		var invoiceDue *models.UpcomingInvoice
		if c.Tenant().Billing.StripeSubscriptionID != "" {
			invoiceDue, err = c.Services().Billing.GetUpcomingInvoice()
			if err != nil {
				return c.Failure(err)
			}
		}

		paymentInfo, err := c.Services().Billing.GetPaymentInfo()
		if err != nil {
			return c.Failure(err)
		}

		var plans []*models.BillingPlan
		if paymentInfo != nil {
			plans, err = c.Services().Billing.ListPlans(paymentInfo.AddressCountry)
			if err != nil {
				return c.Failure(err)
			}
		}

		tenantUserCount, err := c.Services().Users.Count()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Billing · Site Settings",
			ChunkName: "Billing.page",
			Data: web.Map{
				"invoiceDue":      invoiceDue,
				"tenantUserCount": tenantUserCount,
				"plans":           plans,
				"paymentInfo":     paymentInfo,
				"countries":       models.GetAllCountries(),
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

// GetBillingPlans returns a list of plans for given country code
func GetBillingPlans() web.HandlerFunc {
	return func(c web.Context) error {
		countryCode := c.Param("countryCode")
		plans, err := c.Services().Billing.ListPlans(countryCode)
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(plans)
	}
}

// BillingSubscribe subscribes current tenant to given plan on stripe
func BillingSubscribe() web.HandlerFunc {
	return func(c web.Context) error {
		planID := c.Param("planID")

		paymentInfo, err := c.Services().Billing.GetPaymentInfo()
		if err != nil {
			return c.Failure(err)
		}

		plan, err := c.Services().Billing.GetPlanByID(paymentInfo.AddressCountry, planID)
		if err != nil {
			return c.Failure(err)
		}

		userCount, err := c.Services().Users.Count()
		if err != nil {
			return c.Failure(err)
		}

		if plan.MaxUsers > 0 && userCount > plan.MaxUsers {
			return c.Unauthorized()
		}

		err = c.Services().Billing.Subscribe(plan.ID)
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
