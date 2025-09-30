package dbEntities

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/dbx"
)

type billingState struct {
	Status             int          `db:"status"`
	PlanID             string       `db:"paddle_plan_id"`
	SubscriptionID     string       `db:"paddle_subscription_id"`
	TrialEndsAt        dbx.NullTime `db:"trial_ends_at"`
	SubscriptionEndsAt dbx.NullTime `db:"subscription_ends_at"`
}

func (s *billingState) toModel(ctx context.Context) *entity.BillingState {
	model := &entity.BillingState{
		Status:         enum.BillingStatus(s.Status),
		PlanID:         s.PlanID,
		SubscriptionID: s.SubscriptionID,
	}

	if s.TrialEndsAt.Valid {
		model.TrialEndsAt = &s.TrialEndsAt.Time
	}

	if s.SubscriptionEndsAt.Valid {
		model.SubscriptionEndsAt = &s.SubscriptionEndsAt.Time
	}

	return model
}
