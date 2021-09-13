package postgres

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
)

type dbBillingState struct {
	Status             int          `db:"status"`
	PlanID             string       `db:"paddle_plan_id"`
	SubscriptionID     string       `db:"paddle_subscription_id"`
	TrialEndsAt        dbx.NullTime `db:"trial_ends_at"`
	SubscriptionEndsAt dbx.NullTime `db:"subscription_ends_at"`
}

func (s *dbBillingState) toModel(ctx context.Context) *entity.BillingState {
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

func getBillingState(ctx context.Context, q *query.GetBillingState) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil

		state := dbBillingState{}
		err := trx.Get(&state,
			`SELECT
				trial_ends_at,
				subscription_ends_at,
				paddle_subscription_id,
				paddle_plan_id,
				status
			FROM tenants_billing
			WHERE tenant_id = $1`, tenant.ID)

		if err != nil {
			return err
		}

		q.Result = state.toModel(ctx)
		return nil
	})
}
