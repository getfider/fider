package cmd

type ActivateBillingSubscription struct {
	TenantID int
}

type CancelBillingSubscription struct {
	TenantID int
}

type ActivateStripeSubscription struct {
	TenantID       int
	CustomerID     string
	SubscriptionID string
}

type CancelStripeSubscription struct {
	TenantID int
}
