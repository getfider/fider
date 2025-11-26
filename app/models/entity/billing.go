package entity

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

type BillingState struct {
	Status             enum.BillingStatus `json:"status"`
	SubscriptionEndsAt *time.Time         `json:"subscriptionEndsAt"`
}

type BillingSubscription struct {
	UpdateURL          string                    `json:"updateURL"`
	CancelURL          string                    `json:"cancelURL"`
	PaymentInformation BillingPaymentInformation `json:"paymentInformation"`
	LastPayment        BillingPayment            `json:"lastPayment"`
	NextPayment        BillingPayment            `json:"nextPayment"`
}

type BillingPaymentInformation struct {
	PaymentMethod  string `json:"paymentMethod"`
	CardType       string `json:"cardType"`
	LastFourDigits string `json:"lastFourDigits"`
	ExpiryDate     string `json:"expiryDate"`
}

type BillingPayment struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

type StripeBillingState struct {
	CustomerID     string `json:"customerID"`
	SubscriptionID string `json:"subscriptionID"`
}
