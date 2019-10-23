package billing

import (
	"context"
	"fmt"
	"strconv"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/goenning/vat"
	"github.com/stripe/stripe-go"
)

func createCustomer(ctx context.Context, c *cmd.CreateBillingCustomer) error {
	return using(ctx, func(tenant *models.Tenant) error {
		if tenant.Billing == nil {
			return errors.New("Tenant doesn't have a billing record")
		}

		if tenant.Billing.StripeCustomerID == "" {
			params := &stripe.CustomerParams{
				Description: stripe.String(tenant.Name),
			}
			params.AddMetadata("tenant_id", strconv.Itoa(tenant.ID))
			params.AddMetadata("tenant_subdomain", tenant.Subdomain)
			customer, err := stripeClient.Customers.New(params)
			if err != nil {
				return errors.Wrap(err, "failed to create Stripe customer")
			}

			tenant.Billing.StripeCustomerID = customer.ID
			return nil
		}

		return nil
	})
}

func deleteCustomer(ctx context.Context, c *cmd.DeleteBillingCustomer) error {
	return using(ctx, func(tenant *models.Tenant) error {
		if !env.IsTest() {
			return errors.New("Stripe customer can only be deleted on test mode")
		}

		_, err := stripeClient.Customers.Del(tenant.Billing.StripeCustomerID, &stripe.CustomerParams{})
		if err != nil {
			return errors.Wrap(err, "failed to delete Stripe customer")
		}
		return nil
	})
}

func getPaymentInfo(ctx context.Context, q *query.GetPaymentInfo) error {
	return using(ctx, func(tenant *models.Tenant) error {
		if tenant.Billing == nil || tenant.Billing.StripeCustomerID == "" {
			return nil
		}

		customerID := tenant.Billing.StripeCustomerID

		customer, err := stripeClient.Customers.Get(customerID, &stripe.CustomerParams{})
		if err != nil {
			return errors.Wrap(err, "failed to get customer")
		}

		if customer.Metadata["tenant_id"] != strconv.Itoa(tenant.ID) {
			panic(fmt.Sprintf("Stripe TenantID (%s) doesn't match current Tenant ID (%s). Aborting.", customer.Metadata["tenant_id"], strconv.Itoa(tenant.ID)))
		}

		if customer.DefaultSource == nil {
			return nil
		}

		card, err := stripeClient.Cards.Get(customer.DefaultSource.ID, &stripe.CardParams{
			Customer: stripe.String(customerID),
		})
		if err != nil {
			return errors.Wrap(err, "failed to get customer's card")
		}

		q.Result = &dto.PaymentInfo{
			Email:             customer.Email,
			Name:              card.Name,
			StripeCardID:      card.ID,
			CardCountry:       card.Country,
			CardBrand:         string(card.Brand),
			CardLast4:         card.Last4,
			CardExpMonth:      card.ExpMonth,
			CardExpYear:       card.ExpYear,
			AddressCity:       card.AddressCity,
			AddressCountry:    card.AddressCountry,
			AddressLine1:      card.AddressLine1,
			AddressLine2:      card.AddressLine2,
			AddressState:      card.AddressState,
			AddressPostalCode: card.AddressZip,
		}

		if customer.TaxInfo != nil {
			q.Result.VATNumber = customer.TaxInfo.TaxID
		}

		return nil
	})
}

func clearPaymentInfo(ctx context.Context, c *cmd.ClearPaymentInfo) error {
	return using(ctx, func(tenant *models.Tenant) error {
		currentInfo := &query.GetPaymentInfo{}
		if err := getPaymentInfo(ctx, currentInfo); err != nil {
			return err
		}

		if currentInfo.Result != nil {
			customerID := tenant.Billing.StripeCustomerID
			_, err := stripeClient.Customers.Update(customerID, &stripe.CustomerParams{
				Description: stripe.String(tenant.Name),
				Email:       stripe.String(""),
				TaxInfo: &stripe.CustomerTaxInfoParams{
					Type:  stripe.String(string(stripe.CustomerTaxInfoTypeVAT)),
					TaxID: stripe.String(""),
				},
			})
			if err != nil {
				return errors.Wrap(err, "failed to delete customer billing email")
			}
			if currentInfo.Result.StripeCardID != "" {
				_, err = stripeClient.Cards.Del(currentInfo.Result.StripeCardID, &stripe.CardParams{
					Customer: stripe.String(customerID),
				})
				if err != nil {
					return errors.Wrap(err, "failed to delete customer card")
				}
			}
		}

		return nil
	})
}

func updatePaymentInfo(ctx context.Context, c *cmd.UpdatePaymentInfo) error {
	return using(ctx, func(tenant *models.Tenant) error {
		customerID := tenant.Billing.StripeCustomerID

		currentInfo := &query.GetPaymentInfo{}
		if err := getPaymentInfo(ctx, currentInfo); err != nil {
			return err
		}

		if !vat.IsEU(c.Input.AddressCountry) {
			c.Input.VATNumber = ""
		}

		// update customer info
		params := &stripe.CustomerParams{
			Email:       stripe.String(c.Input.Email),
			Description: stripe.String(c.Input.Name),
			Shipping: &stripe.CustomerShippingDetailsParams{
				Name: stripe.String(c.Input.Name),
				Address: &stripe.AddressParams{
					City:       stripe.String(c.Input.AddressCity),
					Country:    stripe.String(c.Input.AddressCountry),
					Line1:      stripe.String(c.Input.AddressLine1),
					Line2:      stripe.String(c.Input.AddressLine2),
					PostalCode: stripe.String(c.Input.AddressPostalCode),
					State:      stripe.String(c.Input.AddressState),
				},
			},
			TaxInfo: &stripe.CustomerTaxInfoParams{
				Type:  stripe.String(string(stripe.CustomerTaxInfoTypeVAT)),
				TaxID: stripe.String(c.Input.VATNumber),
			},
		}
		_, err := stripeClient.Customers.Update(customerID, params)
		if err != nil {
			return errors.Wrap(err, "failed to update customer billing email")
		}

		// new card, just create it
		if currentInfo.Result == nil || currentInfo.Result.StripeCardID == "" {
			_, err = stripeClient.Cards.New(&stripe.CardParams{
				Customer: stripe.String(customerID),
				Token:    stripe.String(c.Input.Card.Token),
			})
			if err != nil {
				return errors.Wrap(err, "failed to create stripe card")
			}
			return nil
		}

		// replacing card, create new and delete old
		if c.Input.Card != nil && c.Input.Card.Token != "" {
			_, err = stripeClient.Cards.New(&stripe.CardParams{
				Customer: stripe.String(customerID),
				Token:    stripe.String(c.Input.Card.Token),
			})
			if err != nil {
				return errors.Wrap(err, "failed to create new stripe card")
			}

			_, err = stripeClient.Cards.Del(currentInfo.Result.StripeCardID, &stripe.CardParams{
				Customer: stripe.String(customerID),
				Token:    stripe.String(c.Input.Card.Token),
			})
			if err != nil {
				return errors.Wrap(err, "failed to delete old stripe card")
			}
			return nil
		}

		// updating card, just update current card
		_, err = stripeClient.Cards.Update(currentInfo.Result.StripeCardID, &stripe.CardParams{
			Customer:       stripe.String(customerID),
			Name:           stripe.String(c.Input.Name),
			AddressCity:    stripe.String(c.Input.AddressCity),
			AddressCountry: stripe.String(c.Input.AddressCountry),
			AddressLine1:   stripe.String(c.Input.AddressLine1),
			AddressLine2:   stripe.String(c.Input.AddressLine2),
			AddressState:   stripe.String(c.Input.AddressState),
			AddressZip:     stripe.String(c.Input.AddressPostalCode),
		})
		if err != nil {
			return errors.Wrap(err, "failed to update stripe card")
		}
		return nil
	})
}
