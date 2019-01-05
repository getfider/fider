package billing

import (
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var sc *client.API

func init() {
	stripe.LogLevel = 0
	sc = &client.API{}
	sc.Init(env.Config.Stripe.SecretKey, nil)
}

// Client is a billing client wrapper for Stripe
type Client struct {
	sc     *client.API
	tenant *models.Tenant
	user   *models.User
}

// NewClient creates a new billing client
func NewClient() *Client {
	return &Client{
		sc: sc,
	}
}

// SetCurrentTenant to current context
func (c *Client) SetCurrentTenant(tenant *models.Tenant) {
	c.tenant = tenant
}

// SetCurrentUser to current context
func (c *Client) SetCurrentUser(user *models.User) {
	c.user = user
}

// CreateCustomer on stripe
func (c *Client) CreateCustomer(email string) (string, error) {
	if c.tenant.Billing == nil {
		return "", errors.New("Tenant doesn't have a billing record")
	}

	if c.tenant.Billing.StripeCustomerID == "" {
		customer, err := c.sc.Customers.New(&stripe.CustomerParams{
			Email:       stripe.String(email),
			Description: stripe.String(fmt.Sprintf("%s [%s]", c.tenant.Name, c.tenant.Subdomain)),
		})
		if err != nil {
			return "", errors.Wrap(err, "failed to create Stripe customer")
		}
		return customer.ID, nil
	}

	return c.tenant.Billing.StripeCustomerID, nil
}

// GetPaymentInfo from a stripe card
func (c *Client) GetPaymentInfo() (*models.PaymentInfo, error) {
	customerID := c.tenant.Billing.StripeCustomerID
	if customerID == "" {
		return nil, nil
	}

	customer, err := c.sc.Customers.Get(customerID, &stripe.CustomerParams{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer")
	}

	info := &models.PaymentInfo{
		Email: customer.Email,
	}

	if customer.DefaultSource != nil {
		card, err := c.sc.Cards.Get(customer.DefaultSource.ID, &stripe.CardParams{
			Customer: stripe.String(customerID),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get customer's card")
		}

		info.StripeCardID = card.ID
		info.CardCountry = card.Country
		info.CardBrand = string(card.Brand)
		info.CardLast4 = card.Last4
		info.CardExpMonth = card.ExpMonth
		info.CardExpYear = card.ExpYear
		info.AddressCity = card.AddressCity
		info.AddressCountry = card.AddressCountry
		info.Name = card.Name
		info.AddressLine1 = card.AddressLine1
		info.AddressLine2 = card.AddressLine2
		info.AddressState = card.AddressState
		info.AddressPostalCode = card.AddressZip
	}

	return info, nil
}

// UpdatePaymentInfo creates or updates customer payment info on stripe
func (c *Client) UpdatePaymentInfo(input *models.CreateEditBillingPaymentInfo) error {
	customerID := c.tenant.Billing.StripeCustomerID
	current, err := c.GetPaymentInfo()
	if err != nil {
		return err
	}

	if current.Email != input.Email {
		_, err = c.sc.Customers.Update(customerID, &stripe.CustomerParams{
			Email:       stripe.String(input.Email),
			Description: stripe.String(fmt.Sprintf("%s [%s]", c.tenant.Name, c.tenant.Subdomain)),
		})
		if err != nil {
			return errors.Wrap(err, "failed to update customer billing email")
		}
	}

	if current.StripeCardID == "" {
		_, err = c.sc.Cards.New(&stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create stripe card")
		}
	} else {
		_, err = c.sc.Cards.Update(current.StripeCardID, &stripe.CardParams{
			Customer:       stripe.String(customerID),
			Name:           stripe.String(input.Name),
			AddressCity:    stripe.String(input.AddressCity),
			AddressCountry: stripe.String(input.AddressCountry),
			AddressLine1:   stripe.String(input.AddressLine1),
			AddressLine2:   stripe.String(input.AddressLine2),
			AddressState:   stripe.String(input.AddressState),
			AddressZip:     stripe.String(input.AddressPostalCode),
		})
		if err != nil {
			return errors.Wrap(err, "failed to update stripe card")
		}
	}

	return nil
}
