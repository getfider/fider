package dto

import "time"

// Country is a valid country within Fider
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
	IsEU bool   `json:"isEU"`
}

// BillingPlan is the model for billing plan from Stripe
type BillingPlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Currency    string `json:"currency"`
	MaxUsers    int    `json:"maxUsers"`
	Price       int64  `json:"price"`
	Interval    string `json:"interval"`
}

// UpcomingInvoice is the model for upcoming invoice from Stripe
type UpcomingInvoice struct {
	Currency  string    `json:"currency"`
	DueDate   time.Time `json:"dueDate"`
	AmountDue int64     `json:"amountDue"`
}

// PaymentInfo is the model for billing payment info
type PaymentInfo struct {
	StripeCardID      string `json:"-"`
	CardCountry       string `json:"cardCountry"`
	CardBrand         string `json:"cardBrand"`
	CardLast4         string `json:"cardLast4"`
	CardExpMonth      uint8  `json:"cardExpMonth"`
	CardExpYear       uint16 `json:"cardExpYear"`
	AddressCity       string `json:"addressCity"`
	AddressCountry    string `json:"addressCountry"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	AddressState      string `json:"addressState"`
	AddressPostalCode string `json:"addressPostalCode"`
	VATNumber         string `json:"vatNumber"`
}

// CreateEditBillingPaymentInfo is the input model to create or edit billing payment info
type CreateEditBillingPaymentInfo struct {
	Name              string                            `json:"name"`
	Email             string                            `json:"email"`
	Card              *CreateEditBillingPaymentInfoCard `json:"card"`
	AddressLine1      string                            `json:"addressLine1"`
	AddressLine2      string                            `json:"addressLine2"`
	AddressCity       string                            `json:"addressCity"`
	AddressState      string                            `json:"addressState"`
	AddressPostalCode string                            `json:"addressPostalCode"`
	AddressCountry    string                            `json:"addressCountry" format:"upper"`
	VATNumber         string                            `json:"vatNumber"`
}

// CreateEditBillingPaymentInfoCard is the input model for a card during billing payment info update
type CreateEditBillingPaymentInfoCard struct {
	Type    string `json:"type"`
	Token   string `json:"token"`
	Country string `json:"country"`
}
