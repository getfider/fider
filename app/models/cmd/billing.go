package cmd

import (
	"time"

	"github.com/getfider/fider/app/models/dto"
)

type GenerateCheckoutLink struct {
	Email       string
	ReturnURL   string
	Passthrough dto.PaddlePassthrough

	// Output
	URL string
}

type ActivateBillingSubscription struct {
	TenantID       int
	SubscriptionID string
	PlanID         string
}

type CancelBillingSubscription struct {
	TenantID           int
	SubscriptionEndsAt time.Time
}

type LockExpiredTenants struct {
	//Output
	NumOfTenantsLocked int64
}
