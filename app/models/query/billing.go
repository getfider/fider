package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetBillingSubscription struct {
	SubscriptionID string

	// Output
	Result *entity.BillingSubscription
}

type GetStripeBillingState struct {
	// Output
	Result *entity.StripeBillingState
}
