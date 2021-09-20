package entity

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

type BillingState struct {
	Status             enum.BillingStatus `json:"status"`
	PlanID             string             `json:"planID"`
	SubscriptionID     string             `json:"subscriptionID"`
	TrialEndsAt        *time.Time         `json:"trialEndsAt"`
	SubscriptionEndsAt *time.Time         `json:"subscriptionEndsAt"`
}

type BillingSubscription struct {
	SignupAt           string                    `json:"signupAt"`
	UpdateURL          string                    `json:"updateURL"`
	CancelURL          string                    `json:"cancelURL"`
	PaymentInformation BillingPaymentInformation `json:"paymentInformation"`
	LastPayment        BillingLastPayment        `json:"lastPayment"`
}

type BillingPaymentInformation struct {
	PaymentMethod  string `json:"paymentMethod"`
	CardType       string `json:"cardType"`
	LastFourDigits string `json:"lastFourDigits"`
	ExpiryDate     string `json:"expiryDate"`
}

type BillingLastPayment struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}
