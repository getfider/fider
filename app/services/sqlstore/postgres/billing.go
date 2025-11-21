package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
	"github.com/lib/pq"
)

func getBillingState(ctx context.Context, q *query.GetBillingState) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil

		state := dbEntities.BillingState{}
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

		q.Result = state.ToModel(ctx)
		return nil
	})
}

func activateBillingSubscription(ctx context.Context, c *cmd.ActivateBillingSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants_billing
			SET subscription_ends_at = null, paddle_subscription_id = $2, paddle_plan_id = $3, status = $4
			WHERE tenant_id = $1
		`, c.TenantID, c.SubscriptionID, c.PlanID, enum.BillingActive)
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

func lockExpiredTenants(ctx context.Context, c *cmd.LockExpiredTenants) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		now := time.Now()

		type tenant struct {
			Id int `db:"id"`
		}

		tenants := []*tenant{}
		err := trx.Select(&tenants, `
			SELECT id
			FROM tenants t
			INNER JOIN tenants_billing tb
			ON t.id = tb.tenant_id
			WHERE t.status <> $1 AND t.status <> $2
			AND (
				(tb.status = $3 AND trial_ends_at <= $5)
				OR (tb.status = $4 AND subscription_ends_at <= $5)
			)`, enum.TenantLocked, enum.TenantDisabled, enum.BillingTrial, enum.BillingCancelled, now)
		if err != nil {
			return errors.Wrap(err, "failed to get expired trial/cancelled tenants")
		}

		if len(tenants) > 0 {
			ids := make([]int, 0)
			for _, tenant := range tenants {
				ids = append(ids, tenant.Id)
			}

			count, err := trx.Execute(`
				UPDATE tenants
				SET status = $1
				WHERE id = ANY($2)
			`, enum.TenantLocked, pq.Array(ids))
			if err != nil {
				return errors.Wrap(err, "failed to lock trial/cancelled tenants")
			}

			c.NumOfTenantsLocked = count
			c.TenantsLocked = ids
		}

		return nil
	})
}

func getTrialingTenantContacts(ctx context.Context, q *query.GetTrialingTenantContacts) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		var users []*dbEntities.User
		err := trx.Select(&users, `
			SELECT
				u.name,
				u.email,
				u.role,
				u.status,
				t.subdomain as tenant_subdomain
			FROM tenants_billing tb
			INNER JOIN tenants t
			ON t.id = tb.tenant_id
			INNER JOIN users u
			ON u.tenant_id = tb.tenant_id
			AND u.role = $1
			WHERE date(trial_ends_at) = date($2)
			AND tb.status = $3`, enum.RoleAdministrator, q.TrialExpiresOn, enum.BillingTrial)
		if err != nil {
			return errors.Wrap(err, "failed to get trialing tenant contacts")
		}

		q.Contacts = make([]*entity.User, len(users))
		for i, user := range users {
			q.Contacts[i] = user.ToModel(ctx)
		}
		return nil
	})
}
