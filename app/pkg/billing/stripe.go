package billing

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var sc *client.API
var mu sync.RWMutex
var plans []*models.BillingPlan

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
		params := &stripe.CustomerParams{
			Email:       stripe.String(email),
			Description: stripe.String(customerDesc(c.tenant)),
		}
		params.AddMetadata("tenant_id", strconv.Itoa(c.tenant.ID))
		customer, err := c.sc.Customers.New(params)
		if err != nil {
			return "", errors.Wrap(err, "failed to create Stripe customer")
		}
		c.tenant.Billing.StripeCustomerID = customer.ID
		return customer.ID, nil
	}

	return c.tenant.Billing.StripeCustomerID, nil
}

// DeleteCustomer on stripe
func (c *Client) DeleteCustomer() error {
	if !env.IsTest() {
		return errors.New("Stripe customer can only be deleted on test mode")
	}

	_, err := c.sc.Customers.Del(c.tenant.Billing.StripeCustomerID, &stripe.CustomerParams{})
	if err != nil {
		return errors.Wrap(err, "failed to delete Stripe customer")
	}
	return nil
}

// GetPaymentInfo from a stripe card
func (c *Client) GetPaymentInfo() (*models.PaymentInfo, error) {
	if c.tenant.Billing == nil || c.tenant.Billing.StripeCustomerID == "" {
		return nil, nil
	}

	customerID := c.tenant.Billing.StripeCustomerID

	customer, err := c.sc.Customers.Get(customerID, &stripe.CustomerParams{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer")
	}

	if customer.Metadata["tenant_id"] != strconv.Itoa(c.tenant.ID) {
		panic(fmt.Sprintf("Stripe TenantID (%s) doesn't match current Tenant ID (%s). Aborting.", customer.Metadata["tenant_id"], strconv.Itoa(c.tenant.ID)))
	}

	if customer.DefaultSource == nil {
		return nil, nil
	}

	card, err := c.sc.Cards.Get(customer.DefaultSource.ID, &stripe.CardParams{
		Customer: stripe.String(customerID),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer's card")
	}

	info := &models.PaymentInfo{
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

	return info, nil
}

// ClearPaymentInfo removes all payment information from stripe
func (c *Client) ClearPaymentInfo() error {
	current, err := c.GetPaymentInfo()
	if err != nil {
		return err
	}

	if current != nil {
		customerID := c.tenant.Billing.StripeCustomerID
		_, err = c.sc.Customers.Update(customerID, &stripe.CustomerParams{
			Email: stripe.String(""),
		})
		if err != nil {
			return errors.Wrap(err, "failed to delete customer billing email")
		}
		if current.StripeCardID != "" {
			_, err = c.sc.Cards.Del(current.StripeCardID, &stripe.CardParams{
				Customer: stripe.String(customerID),
			})
			if err != nil {
				return errors.Wrap(err, "failed to delete customer card")
			}
		}
	}

	return nil
}

// UpdatePaymentInfo creates or updates customer payment info on stripe
func (c *Client) UpdatePaymentInfo(input *models.CreateEditBillingPaymentInfo) error {
	customerID := c.tenant.Billing.StripeCustomerID
	current, err := c.GetPaymentInfo()
	if err != nil {
		return err
	}

	// email is different, update it
	if current == nil || current.Email != input.Email {
		_, err = c.sc.Customers.Update(customerID, &stripe.CustomerParams{
			Email:       stripe.String(input.Email),
			Description: stripe.String(customerDesc(c.tenant)),
		})
		if err != nil {
			return errors.Wrap(err, "failed to update customer billing email")
		}
	}

	// new card, just create it
	if current == nil || current.StripeCardID == "" {
		_, err = c.sc.Cards.New(&stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create stripe card")
		}
		return nil
	}

	// replacing card, create new and delete old
	if input.Card != nil && input.Card.Token != "" {
		_, err = c.sc.Cards.New(&stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create new stripe card")
		}

		_, err = c.sc.Cards.Del(current.StripeCardID, &stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to delete old stripe card")
		}
		return nil
	}

	// updating card, just update current card
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
	return nil
}

func customerDesc(tenant *models.Tenant) string {
	return fmt.Sprintf("%s [%s]", tenant.Name, tenant.Subdomain)
}

// ListPlans on stripe
func (c *Client) ListPlans() ([]*models.BillingPlan, error) {
	if plans != nil {
		return plans, nil
	}

	mu.Lock()
	defer mu.Unlock()

	if plans == nil {
		plans = make([]*models.BillingPlan, 0)
		it := c.sc.Plans.List(&stripe.PlanListParams{
			Active: stripe.Bool(true),
		})
		for it.Next() {
			plan := it.Plan()
			tier, _ := strconv.Atoi(plan.Metadata["tier"])
			plans = append(plans, &models.BillingPlan{
				ID:          plan.ID,
				Name:        plan.Nickname,
				Description: plan.Metadata["description"],
				Tier:        tier,
				Price:       plan.Amount,
				Interval:    string(plan.Interval),
			})
		}
		if err := it.Err(); err != nil {
			return nil, err
		}
		sort.Slice(plans, func(i, j int) bool {
			return plans[i].Price < plans[j].Price
		})
	}

	return plans, nil
}
