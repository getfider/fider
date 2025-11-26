package postgres

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func getBillingState(ctx context.Context, q *query.GetBillingState) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil

		state := dbEntities.BillingState{}
		err := trx.Get(&state,
			`SELECT
				subscription_ends_at,
				status
			FROM tenants_billing
			WHERE tenant_id = $1`, tenant.ID)

		if err != nil {
			return err
		}

		q.Result = state.ToModel(ctx)
		return nil
	})
}

func activateBillingSubscription(ctx context.Context, c *cmd.ActivateBillingSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants_billing
			SET subscription_ends_at = null, status = $2
			WHERE tenant_id = $1
		`, c.TenantID, enum.BillingActive)
		if err != nil {
			return errors.Wrap(err, "failed activate billing subscription")
		}

		_, err = trx.Execute(`
			UPDATE tenants
			SET status = $2
			WHERE id = $1
		`, c.TenantID, enum.TenantActive)
		if err != nil {
			return errors.Wrap(err, "failed activate tenant")
		}

		return nil
	})
}

func cancelBillingSubscription(ctx context.Context, c *cmd.CancelBillingSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants_billing
			SET subscription_ends_at = $2, status = $3
			WHERE tenant_id = $1
		`, c.TenantID, c.SubscriptionEndsAt, enum.BillingCancelled)
		if err != nil {
			return errors.Wrap(err, "failed cancel billing subscription")
		}
		return nil
	})
}

func getStripeBillingState(ctx context.Context, q *query.GetStripeBillingState) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		state := dbEntities.StripeBillingState{}
		err := trx.Get(&state,
			`SELECT stripe_customer_id, stripe_subscription_id
			FROM tenants_billing
			WHERE tenant_id = $1`, tenant.ID)

		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				// No billing record for this tenant, return empty state
				q.Result = &entity.StripeBillingState{}
				return nil
			}
			return err
		}

		q.Result = &entity.StripeBillingState{
			CustomerID:     state.StripeCustomerID.String,
			SubscriptionID: state.StripeSubscriptionID.String,
		}
		return nil
	})
}

func activateStripeSubscription(ctx context.Context, c *cmd.ActivateStripeSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			INSERT INTO tenants_billing (tenant_id, status, stripe_customer_id, stripe_subscription_id)
			VALUES ($1, 0, $2, $3)
			ON CONFLICT (tenant_id) DO UPDATE
			SET stripe_customer_id = $2, stripe_subscription_id = $3
		`, c.TenantID, c.CustomerID, c.SubscriptionID)
		if err != nil {
			return errors.Wrap(err, "failed to activate stripe subscription")
		}
		return nil
	})
}

func cancelStripeSubscription(ctx context.Context, c *cmd.CancelStripeSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants_billing
			SET stripe_subscription_id = NULL
			WHERE tenant_id = $1
		`, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to cancel stripe subscription")
		}
		return nil
	})
}
