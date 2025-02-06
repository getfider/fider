package cmd

import (
	"time"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
)

type GenerateCheckoutLink struct {
	PlanID      string
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
	TenantsLocked      []int
}
