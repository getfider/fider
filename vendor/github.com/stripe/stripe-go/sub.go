package stripe

import (
	"encoding/json"

	"github.com/stripe/stripe-go/form"
)

// SubscriptionStatus is the list of allowed values for the subscription's status.
type SubscriptionStatus string

// List of values that SubscriptionStatus can take.
const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusAll      SubscriptionStatus = "all"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	SubscriptionStatusUnpaid   SubscriptionStatus = "unpaid"
)

// SubscriptionBilling is the type of billing method for this subscription's invoices.
type SubscriptionBilling string

// List of values that SubscriptionBilling can take.
const (
	SubscriptionBillingChargeAutomatically SubscriptionBilling = "charge_automatically"
	SubscriptionBillingSendInvoice         SubscriptionBilling = "send_invoice"
)

// SubscriptionParams is the set of parameters that can be used when creating or updating a subscription.
// For more details see https://stripe.com/docs/api#create_subscription and https://stripe.com/docs/api#update_subscription.
type SubscriptionParams struct {
	Params                      `form:"*"`
	ApplicationFeePercent       *float64                   `form:"application_fee_percent"`
	Billing                     *string                    `form:"billing"`
	BillingCycleAnchor          *int64                     `form:"billing_cycle_anchor"`
	BillingCycleAnchorNow       *bool                      `form:"-"` // See custom AppendTo
	BillingCycleAnchorUnchanged *bool                      `form:"-"` // See custom AppendTo
	CancelAtPeriodEnd           *bool                      `form:"cancel_at_period_end"`
	Card                        *CardParams                `form:"card"`
	Coupon                      *string                    `form:"coupon"`
	Customer                    *string                    `form:"customer"`
	DaysUntilDue                *int64                     `form:"days_until_due"`
	DefaultSource               *string                    `form:"default_source"`
	Items                       []*SubscriptionItemsParams `form:"items"`
	OnBehalfOf                  *string                    `form:"on_behalf_of"`
	Plan                        *string                    `form:"plan"`
	Prorate                     *bool                      `form:"prorate"`
	ProrationDate               *int64                     `form:"proration_date"`
	Quantity                    *int64                     `form:"quantity"`
	TaxPercent                  *float64                   `form:"tax_percent"`
	TrialEnd                    *int64                     `form:"trial_end"`
	TrialEndNow                 *bool                      `form:"-"` // See custom AppendTo
	TrialFromPlan               *bool                      `form:"trial_from_plan"`
	TrialPeriodDays             *int64                     `form:"trial_period_days"`
}

// SubscriptionCancelParams is the set of parameters that can be used when canceling a subscription.
// For more details see https://stripe.com/docs/api#cancel_subscription
type SubscriptionCancelParams struct {
	Params     `form:"*"`
	InvoiceNow *bool `form:"invoice_now"`
	Prorate    *bool `form:"prorate"`
}

// AppendTo implements custom encoding logic for SubscriptionParams so that the special
// "now" value for billing_cycle_anchor and trial_end can be implemented
// (they're otherwise timestamps rather than strings).
func (p *SubscriptionParams) AppendTo(body *form.Values, keyParts []string) {
	if BoolValue(p.BillingCycleAnchorNow) {
		body.Add(form.FormatKey(append(keyParts, "billing_cycle_anchor")), "now")
	}

	if BoolValue(p.BillingCycleAnchorUnchanged) {
		body.Add(form.FormatKey(append(keyParts, "billing_cycle_anchor")), "unchanged")
	}

	if BoolValue(p.TrialEndNow) {
		body.Add(form.FormatKey(append(keyParts, "trial_end")), "now")
	}
}

// SubscriptionItemsParams is the set of parameters that can be used when creating or updating a subscription item on a subscription
// For more details see https://stripe.com/docs/api#create_subscription and https://stripe.com/docs/api#update_subscription.
type SubscriptionItemsParams struct {
	Params     `form:"*"`
	ClearUsage *bool   `form:"clear_usage"`
	Deleted    *bool   `form:"deleted"`
	ID         *string `form:"id"`
	Plan       *string `form:"plan"`
	Quantity   *int64  `form:"quantity"`
}

// SubscriptionListParams is the set of parameters that can be used when listing active subscriptions.
// For more details see https://stripe.com/docs/api#list_subscriptions.
type SubscriptionListParams struct {
	ListParams   `form:"*"`
	Billing      string            `form:"billing"`
	Created      int64             `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
	Customer     string            `form:"customer"`
	Plan         string            `form:"plan"`
	Status       string            `form:"status"`
}

// Subscription is the resource representing a Stripe subscription.
// For more details see https://stripe.com/docs/api#subscriptions.
type Subscription struct {
	ApplicationFeePercent float64               `json:"application_fee_percent"`
	Billing               SubscriptionBilling   `json:"billing"`
	BillingCycleAnchor    int64                 `json:"billing_cycle_anchor"`
	CanceledAt            int64                 `json:"canceled_at"`
	Created               int64                 `json:"created"`
	CurrentPeriodEnd      int64                 `json:"current_period_end"`
	CurrentPeriodStart    int64                 `json:"current_period_start"`
	Customer              *Customer             `json:"customer"`
	DaysUntilDue          int64                 `json:"days_until_due"`
	DefaultSource         *PaymentSource        `json:"default_source"`
	Discount              *Discount             `json:"discount"`
	CancelAtPeriodEnd     bool                  `json:"cancel_at_period_end"`
	EndedAt               int64                 `json:"ended_at"`
	ID                    string                `json:"id"`
	Items                 *SubscriptionItemList `json:"items"`
	Metadata              map[string]string     `json:"metadata"`
	OnBehalfOf            *Account              `json:"on_behalf_of"`
	Plan                  *Plan                 `json:"plan"`
	Quantity              int64                 `json:"quantity"`
	Start                 int64                 `json:"start"`
	Status                SubscriptionStatus    `json:"status"`
	TaxPercent            float64               `json:"tax_percent"`
	TrialEnd              int64                 `json:"trial_end"`
	TrialStart            int64                 `json:"trial_start"`
}

// SubscriptionList is a list object for subscriptions.
type SubscriptionList struct {
	ListMeta
	Data []*Subscription `json:"data"`
}

// UnmarshalJSON handles deserialization of a Subscription.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (s *Subscription) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		s.ID = id
		return nil
	}

	type subscription Subscription
	var v subscription
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Subscription(v)
	return nil
}
