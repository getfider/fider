package query

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
)

type GetBillingState struct {
	// Output
	Result *entity.BillingState
}

type GetBillingSubscription struct {
	SubscriptionID string

	// Output
	Result *entity.BillingSubscription
}
