package dbEntities

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/dbx"
)

type BillingState struct {
	Status             int          `db:"status"`
	SubscriptionEndsAt dbx.NullTime `db:"subscription_ends_at"`
}

type StripeBillingState struct {
	StripeCustomerID     dbx.NullString `db:"stripe_customer_id"`
	StripeSubscriptionID dbx.NullString `db:"stripe_subscription_id"`
}

func (s *BillingState) ToModel(ctx context.Context) *entity.BillingState {
	model := &entity.BillingState{
		Status: enum.BillingStatus(s.Status),
	}

	if s.SubscriptionEndsAt.Valid {
		model.SubscriptionEndsAt = &s.SubscriptionEndsAt.Time
	}

	return model
}
