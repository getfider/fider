package models

// CreateEditBillingPaymentInfo is the input model to create or edit billing payment info
type CreateEditBillingPaymentInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Card  *struct {
		Type    string `json:"type"`
		Token   string `json:"token"`
		Country string `json:"country"`
	} `json:"card"`
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	AddressCity       string `json:"addressCity"`
	AddressState      string `json:"addressState"`
	AddressPostalCode string `json:"addressPostalCode"`
	AddressCountry    string `json:"addressCountry" format:"upper"`
}

// PaymentInfo is the model for billing payment info
type PaymentInfo struct {
	StripeCardID      string `json:"-"`
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
}
