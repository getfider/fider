package billing

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goenning/vat"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var stripeClient *client.API
var mu sync.RWMutex
var plans []*models.BillingPlan

func init() {
	stripe.LogLevel = 0
	stripeClient = &client.API{}
	stripeClient.Init(env.Config.Stripe.SecretKey, nil)
}

// Client is a billing client wrapper for Stripe
type Client struct {
	stripe *client.API
	tenant *models.Tenant
	user   *models.User
}

// NewClient creates a new billing client
func NewClient() *Client {
	return &Client{
		stripe: stripeClient,
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

// CreateCustomer on Stripe
func (c *Client) CreateCustomer(email string) (string, error) {
	if c.tenant.Billing == nil {
		return "", errors.New("Tenant doesn't have a billing record")
	}

	if c.tenant.Billing.StripeCustomerID == "" {
		params := &stripe.CustomerParams{
			Email:       stripe.String(email),
			Description: stripe.String(c.tenant.Name),
		}
		params.AddMetadata("tenant_id", strconv.Itoa(c.tenant.ID))
		params.AddMetadata("tenant_subdomain", c.tenant.Subdomain)
		customer, err := c.stripe.Customers.New(params)
		if err != nil {
			return "", errors.Wrap(err, "failed to create Stripe customer")
		}
		c.tenant.Billing.StripeCustomerID = customer.ID
		return customer.ID, nil
	}

	return c.tenant.Billing.StripeCustomerID, nil
}

// DeleteCustomer on Stripe
func (c *Client) DeleteCustomer() error {
	if !env.IsTest() {
		return errors.New("Stripe customer can only be deleted on test mode")
	}

	_, err := c.stripe.Customers.Del(c.tenant.Billing.StripeCustomerID, &stripe.CustomerParams{})
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

	customer, err := c.stripe.Customers.Get(customerID, &stripe.CustomerParams{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer")
	}

	if customer.Metadata["tenant_id"] != strconv.Itoa(c.tenant.ID) {
		panic(fmt.Sprintf("Stripe TenantID (%s) doesn't match current Tenant ID (%s). Aborting.", customer.Metadata["tenant_id"], strconv.Itoa(c.tenant.ID)))
	}

	if customer.DefaultSource == nil {
		return nil, nil
	}

	card, err := c.stripe.Cards.Get(customer.DefaultSource.ID, &stripe.CardParams{
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

	if customer.TaxInfo != nil {
		info.VATNumber = customer.TaxInfo.TaxID
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
		_, err = c.stripe.Customers.Update(customerID, &stripe.CustomerParams{
			Description: stripe.String(c.tenant.Name),
			Email:       stripe.String(""),
			TaxInfo: &stripe.CustomerTaxInfoParams{
				Type:  stripe.String(string(stripe.CustomerTaxInfoTypeVAT)),
				TaxID: stripe.String(""),
			},
		})
		if err != nil {
			return errors.Wrap(err, "failed to delete customer billing email")
		}
		if current.StripeCardID != "" {
			_, err = c.stripe.Cards.Del(current.StripeCardID, &stripe.CardParams{
				Customer: stripe.String(customerID),
			})
			if err != nil {
				return errors.Wrap(err, "failed to delete customer card")
			}
		}
	}

	return nil
}

// UpdatePaymentInfo creates or updates customer payment info on Stripe
func (c *Client) UpdatePaymentInfo(input *models.CreateEditBillingPaymentInfo) error {
	customerID := c.tenant.Billing.StripeCustomerID
	current, err := c.GetPaymentInfo()
	if err != nil {
		return err
	}

	if !vat.IsEU(input.AddressCountry) {
		input.VATNumber = ""
	}

	// update customer info
	params := &stripe.CustomerParams{
		Email:       stripe.String(input.Email),
		Description: stripe.String(input.Name),
		Shipping: &stripe.CustomerShippingDetailsParams{
			Name: stripe.String(input.Name),
			Address: &stripe.AddressParams{
				City:       stripe.String(input.AddressCity),
				Country:    stripe.String(input.AddressCountry),
				Line1:      stripe.String(input.AddressLine1),
				Line2:      stripe.String(input.AddressLine2),
				PostalCode: stripe.String(input.AddressPostalCode),
				State:      stripe.String(input.AddressState),
			},
		},
		TaxInfo: &stripe.CustomerTaxInfoParams{
			Type:  stripe.String(string(stripe.CustomerTaxInfoTypeVAT)),
			TaxID: stripe.String(input.VATNumber),
		},
	}
	_, err = c.stripe.Customers.Update(customerID, params)
	if err != nil {
		return errors.Wrap(err, "failed to update customer billing email")
	}

	// new card, just create it
	if current == nil || current.StripeCardID == "" {
		_, err = c.stripe.Cards.New(&stripe.CardParams{
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
		_, err = c.stripe.Cards.New(&stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create new stripe card")
		}

		_, err = c.stripe.Cards.Del(current.StripeCardID, &stripe.CardParams{
			Customer: stripe.String(customerID),
			Token:    stripe.String(input.Card.Token),
		})
		if err != nil {
			return errors.Wrap(err, "failed to delete old stripe card")
		}
		return nil
	}

	// updating card, just update current card
	_, err = c.stripe.Cards.Update(current.StripeCardID, &stripe.CardParams{
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

// GetPlanByID return a plan by its ID
func (c *Client) GetPlanByID(countryCode, planID string) (*models.BillingPlan, error) {
	plans, err := c.ListPlans(countryCode)
	if err != nil {
		return nil, err
	}

	for _, plan := range plans {
		if plan.ID == planID {
			return plan, nil
		}
	}
	return nil, errors.New("failed to get plan by id: " + planID)
}

// ListPlans on Stripe
func (c *Client) ListPlans(countryCode string) ([]*models.BillingPlan, error) {
	if plans != nil {
		return c.filterByCountryCode(plans, countryCode), nil
	}

	mu.Lock()
	defer mu.Unlock()

	if plans == nil {
		plans = make([]*models.BillingPlan, 0)
		it := c.stripe.Plans.List(&stripe.PlanListParams{
			Active: stripe.Bool(true),
		})
		for it.Next() {
			plan := it.Plan()
			name, ok := plan.Metadata["friendly_name"]
			if !ok {
				name = plan.Nickname
			}
			maxUsers, _ := strconv.Atoi(plan.Metadata["max_users"])
			plans = append(plans, &models.BillingPlan{
				ID:          plan.ID,
				Name:        name,
				Description: plan.Metadata["description"],
				MaxUsers:    maxUsers,
				Price:       plan.Amount,
				Currency:    strings.ToUpper(string(plan.Currency)),
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

	return c.filterByCountryCode(plans, countryCode), nil
}

func (c *Client) filterByCountryCode(plans []*models.BillingPlan, countryCode string) []*models.BillingPlan {
	currency := "USD"
	if vat.IsEU(countryCode) {
		currency = "EUR"
	}

	filteredPlans := make([]*models.BillingPlan, 0)
	for _, p := range plans {
		if p.Currency == currency {
			filteredPlans = append(filteredPlans, p)
		}
	}
	return filteredPlans
}

// Subscribe current tenant to given plan on Stripe
func (c *Client) Subscribe(planID string) error {
	customerID := c.tenant.Billing.StripeCustomerID
	if c.tenant.Billing.StripeSubscriptionID != "" {
		sub, err := c.stripe.Subscriptions.Get(c.tenant.Billing.StripeSubscriptionID, nil)
		if err != nil {
			return errors.Wrap(err, "failed to get stripe subscription")
		}
		sub, err = c.stripe.Subscriptions.Update(c.tenant.Billing.StripeSubscriptionID, &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(false),
			Items: []*stripe.SubscriptionItemsParams{
				{
					ID:   stripe.String(sub.Items.Data[0].ID),
					Plan: stripe.String(planID),
				},
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to update stripe subscription")
		}

		c.tenant.Billing.SubscriptionEndsAt = nil
	} else {
		sub, err := c.stripe.Subscriptions.New(&stripe.SubscriptionParams{
			Customer: stripe.String(customerID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Plan: stripe.String(planID),
				},
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to create stripe subscription")
		}

		c.tenant.Billing.StripeSubscriptionID = sub.ID
	}

	c.tenant.Billing.StripePlanID = planID
	return nil
}

// CancelSubscription from current tenant on Stripe
func (c *Client) CancelSubscription() error {
	sub, err := c.stripe.Subscriptions.Update(c.tenant.Billing.StripeSubscriptionID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	})
	if err != nil {
		return errors.Wrap(err, "failed to cancel stripe subscription")
	}
	endDate := time.Unix(sub.CurrentPeriodEnd, 0)
	c.tenant.Billing.SubscriptionEndsAt = &endDate
	return nil
}

// GetUpcomingInvoice returns the next due invoice for current tenant
func (c *Client) GetUpcomingInvoice() (*models.UpcomingInvoice, error) {
	inv, err := c.stripe.Invoices.GetNext(&stripe.InvoiceParams{
		Customer:     stripe.String(c.tenant.Billing.StripeCustomerID),
		Subscription: stripe.String(c.tenant.Billing.StripeSubscriptionID),
	})
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			if stripeErr.HTTPStatusCode == 404 {
				return nil, nil
			}
		}
		return nil, errors.Wrap(err, "failed to get upcoming invoice")
	}

	dueDate := time.Unix(inv.Date, 0)
	return &models.UpcomingInvoice{
		Currency:  strings.ToUpper(string(inv.Currency)),
		AmountDue: inv.AmountDue,
		DueDate:   dueDate,
	}, nil
}
