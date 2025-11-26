package cmd

import (
	"time"
)

type ActivateBillingSubscription struct {
	TenantID int
}

type CancelBillingSubscription struct {
	TenantID           int
	SubscriptionEndsAt time.Time
}

type ActivateStripeSubscription struct {
	TenantID       int
	CustomerID     string
	SubscriptionID string
}

type CancelStripeSubscription struct {
	TenantID int
}
