package stripe

import "encoding/json"

// InvoiceLineType is the list of allowed values for the invoice line's type.
type InvoiceLineType string

// List of values that InvoiceLineType can take.
const (
	InvoiceLineTypeInvoiceItem  InvoiceLineType = "invoiceitem"
	InvoiceLineTypeSubscription InvoiceLineType = "subscription"
)

// InvoiceBilling is the type of billing method for this invoice.
type InvoiceBilling string

// List of values that InvoiceBilling can take.
const (
	InvoiceBillingChargeAutomatically InvoiceBilling = "charge_automatically"
	InvoiceBillingSendInvoice         InvoiceBilling = "send_invoice"
)

// InvoiceBillingReason is the reason why a given invoice was created
type InvoiceBillingReason string

// List of values that InvoiceBillingReason can take.
const (
	InvoiceBillingReasonManual             InvoiceBillingReason = "manual"
	InvoiceBillingReasonSubscription       InvoiceBillingReason = "subscription"
	InvoiceBillingReasonSubscriptionCreate InvoiceBillingReason = "subscription_create"
	InvoiceBillingReasonSubscriptionCycle  InvoiceBillingReason = "subscription_cycle"
	InvoiceBillingReasonSubscriptionUpdate InvoiceBillingReason = "subscription_update"
	InvoiceBillingReasonUpcoming           InvoiceBillingReason = "upcoming"
)

// InvoiceBillingStatus is the reason why a given invoice was created
type InvoiceBillingStatus string

// List of values that InvoiceBillingStatus can take.
const (
	InvoiceBillingStatusDraft         InvoiceBillingStatus = "draft"
	InvoiceBillingStatusOpen          InvoiceBillingStatus = "open"
	InvoiceBillingStatusPaid          InvoiceBillingStatus = "paid"
	InvoiceBillingStatusUncollectible InvoiceBillingStatus = "uncollectible"
	InvoiceBillingStatusVoid          InvoiceBillingStatus = "void"
)

// InvoiceUpcomingInvoiceItemParams is the set of parameters that can be used when adding or modifying
// invoice items on an upcoming invoice.
// For more details see https://stripe.com/docs/api#upcoming_invoice-invoice_items.
type InvoiceUpcomingInvoiceItemParams struct {
	Amount       *int64  `form:"amount"`
	Currency     *string `form:"currency"`
	Description  *string `form:"description"`
	Discountable *bool   `form:"discountable"`
	InvoiceItem  *string `form:"invoiceitem"`
}

// InvoiceParams is the set of parameters that can be used when creating or updating an invoice.
// For more details see https://stripe.com/docs/api#create_invoice, https://stripe.com/docs/api#update_invoice.
type InvoiceParams struct {
	Params              `form:"*"`
	AutoAdvance         *bool    `form:"auto_advance"`
	ApplicationFee      *int64   `form:"application_fee"`
	Billing             *string  `form:"billing"`
	Customer            *string  `form:"customer"`
	DaysUntilDue        *int64   `form:"days_until_due"`
	DefaultSource       *string  `form:"default_source"`
	Description         *string  `form:"description"`
	DueDate             *int64   `form:"due_date"`
	Paid                *bool    `form:"paid"`
	StatementDescriptor *string  `form:"statement_descriptor"`
	Subscription        *string  `form:"subscription"`
	TaxPercent          *float64 `form:"tax_percent"`

	// These are all for exclusive use by GetNext.

	Coupon                         *string                           `form:"coupon"`
	InvoiceItems                   *InvoiceUpcomingInvoiceItemParams `form:"invoice_items"`
	SubscriptionBillingCycleAnchor *int64                            `form:"subscription_billing_cycle_anchor"`
	SubscriptionCancelAtPeriodEnd  *bool                             `form:"subscription_cancel_at_period_end"`
	SubscriptionItems              []*SubscriptionItemsParams        `form:"subscription_items"`
	SubscriptionPlan               *string                           `form:"subscription_plan"`
	SubscriptionProrate            *bool                             `form:"subscription_prorate"`
	SubscriptionProrationDate      *int64                            `form:"subscription_proration_date"`
	SubscriptionQuantity           *int64                            `form:"subscription_quantity"`
	SubscriptionTaxPercent         *float64                          `form:"subscription_tax_percent"`
	SubscriptionTrialEnd           *int64                            `form:"subscription_trial_end"`
	SubscriptionTrialFromPlan      *bool                             `form:"subscription_trial_from_plan"`
}

// InvoiceListParams is the set of parameters that can be used when listing invoices.
// For more details see https://stripe.com/docs/api#list_customer_invoices.
type InvoiceListParams struct {
	ListParams   `form:"*"`
	Billing      *string           `form:"billing"`
	Customer     *string           `form:"customer"`
	Date         *int64            `form:"date"`
	DateRange    *RangeQueryParams `form:"date"`
	DueDate      *int64            `form:"due_date"`
	Subscription *string           `form:"subscription"`
}

// InvoiceLineListParams is the set of parameters that can be used when listing invoice line items.
// For more details see https://stripe.com/docs/api#invoice_lines.
type InvoiceLineListParams struct {
	ListParams `form:"*"`

	Customer *string `form:"customer"`

	// ID is the invoice ID to list invoice lines for.
	ID *string `form:"-"` // Goes in the URL

	Subscription *string `form:"subscription"`
}

// InvoiceFinalizeParams is the set of parameters that can be used when finalizing invoices.
type InvoiceFinalizeParams struct {
	Params      `form:"*"`
	AutoAdvance *bool `form:"auto_advance"`
}

// InvoiceMarkUncollectibleParams is the set of parameters that can be used when marking
// invoices as uncollectible.
type InvoiceMarkUncollectibleParams struct {
	Params `form:"*"`
}

// InvoicePayParams is the set of parameters that can be used when
// paying invoices. For more details, see:
// https://stripe.com/docs/api#pay_invoice.
type InvoicePayParams struct {
	Params        `form:"*"`
	Forgive       *bool   `form:"forgive"`
	PaidOutOfBand *bool   `form:"paid_out_of_band"`
	Source        *string `form:"source"`
}

// InvoiceSendParams is the set of parameters that can be used when sending invoices.
type InvoiceSendParams struct {
	Params `form:"*"`
}

// InvoiceVoidParams is the set of parameters that can be used when voiding invoices.
type InvoiceVoidParams struct {
	Params `form:"*"`
}

// Invoice is the resource representing a Stripe invoice.
// For more details see https://stripe.com/docs/api#invoice_object.
type Invoice struct {
	AmountDue                 int64                `json:"amount_due"`
	AmountPaid                int64                `json:"amount_paid"`
	AmountRemaining           int64                `json:"amount_remaining"`
	ApplicationFee            int64                `json:"application_fee"`
	AttemptCount              int64                `json:"attempt_count"`
	Attempted                 bool                 `json:"attempted"`
	AutoAdvance               bool                 `json:"auto_advance"`
	Billing                   InvoiceBilling       `json:"billing"`
	BillingReason             InvoiceBillingReason `json:"billing_reason"`
	Charge                    *Charge              `json:"charge"`
	Currency                  Currency             `json:"currency"`
	Customer                  *Customer            `json:"customer"`
	Date                      int64                `json:"date"`
	DefaultSource             *PaymentSource       `json:"default_source"`
	Description               string               `json:"description"`
	Discount                  *Discount            `json:"discount"`
	DueDate                   int64                `json:"due_date"`
	EndingBalance             int64                `json:"ending_balance"`
	FinalizedAt               int64                `json:"finalized_at"`
	HostedInvoiceURL          string               `json:"hosted_invoice_url"`
	ID                        string               `json:"id"`
	InvoicePDF                string               `json:"invoice_pdf"`
	Lines                     *InvoiceLineList     `json:"lines"`
	Livemode                  bool                 `json:"livemode"`
	Metadata                  map[string]string    `json:"metadata"`
	NextPaymentAttempt        int64                `json:"next_payment_attempt"`
	Number                    string               `json:"number"`
	Paid                      bool                 `json:"paid"`
	PeriodEnd                 int64                `json:"period_end"`
	PeriodStart               int64                `json:"period_start"`
	ReceiptNumber             string               `json:"receipt_number"`
	StartingBalance           int64                `json:"starting_balance"`
	StatementDescriptor       string               `json:"statement_descriptor"`
	Status                    InvoiceBillingStatus `json:"status"`
	Subscription              string               `json:"subscription"`
	SubscriptionProrationDate int64                `json:"subscription_proration_date"`
	Subtotal                  int64                `json:"subtotal"`
	Tax                       int64                `json:"tax"`
	TaxPercent                float64              `json:"tax_percent"`
	Total                     int64                `json:"total"`
	WebhooksDeliveredAt       int64                `json:"webhooks_delivered_at"`
}

// InvoiceList is a list of invoices as retrieved from a list endpoint.
type InvoiceList struct {
	ListMeta
	Data []*Invoice `json:"data"`
}

// InvoiceLine is the resource representing a Stripe invoice line item.
// For more details see https://stripe.com/docs/api#invoice_line_item_object.
type InvoiceLine struct {
	Amount           int64             `json:"amount"`
	Currency         Currency          `json:"currency"`
	Description      string            `json:"description"`
	Discountable     bool              `json:"discountable"`
	ID               string            `json:"id"`
	Livemode         bool              `json:"live_mode"`
	Metadata         map[string]string `json:"metadata"`
	Period           *Period           `json:"period"`
	Plan             *Plan             `json:"plan"`
	Proration        bool              `json:"proration"`
	Quantity         int64             `json:"quantity"`
	Subscription     string            `json:"subscription"`
	SubscriptionItem string            `json:"subscription_item"`
	Type             InvoiceLineType   `json:"type"`
}

// Period is a structure representing a start and end dates.
type Period struct {
	End   int64 `json:"end"`
	Start int64 `json:"start"`
}

// InvoiceLineList is a list object for invoice line items.
type InvoiceLineList struct {
	ListMeta
	Data []*InvoiceLine `json:"data"`
}

// UnmarshalJSON handles deserialization of an Invoice.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *Invoice) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type invoice Invoice
	var v invoice
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Invoice(v)
	return nil
}
