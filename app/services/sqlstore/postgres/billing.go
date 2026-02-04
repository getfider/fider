package postgres

import (
	"context"
	"database/sql"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func activateBillingSubscription(ctx context.Context, c *cmd.ActivateBillingSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants
			SET is_pro = true, status = $2
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
			UPDATE tenants
			SET is_pro = false
			WHERE id = $1
		`, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to set tenant to free plan")
		}

		return nil
	})
}

func getStripeBillingState(ctx context.Context, q *query.GetStripeBillingState) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		state := dbEntities.StripeBillingState{}
		err := trx.Get(&state,
			`SELECT stripe_customer_id, stripe_subscription_id, license_key, paddle_subscription_id
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
			CustomerID:           state.StripeCustomerID.String,
			SubscriptionID:       state.StripeSubscriptionID.String,
			LicenseKey:           state.LicenseKey.String,
			PaddleSubscriptionID: state.PaddleSubscriptionID.String,
		}
		return nil
	})
}

func activateStripeSubscription(ctx context.Context, c *cmd.ActivateStripeSubscription) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		// Check for existing Paddle subscription before migrating to Stripe
		var paddleSubscriptionID sql.NullString
		err := trx.Scalar(&paddleSubscriptionID, `
			SELECT paddle_subscription_id FROM tenants_billing WHERE tenant_id = $1
		`, c.TenantID)
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return errors.Wrap(err, "failed to check existing paddle subscription")
		}
		if paddleSubscriptionID.Valid && paddleSubscriptionID.String != "" {
			log.Warnf(ctx, "PADDLE_TO_STRIPE_MIGRATION: Tenant @{TenantID} switching from Paddle (subscription: @{PaddleSubscriptionID}) to Stripe (subscription: @{StripeSubscriptionID})", dto.Props{
				"TenantID":             c.TenantID,
				"PaddleSubscriptionID": paddleSubscriptionID.String,
				"StripeSubscriptionID": c.SubscriptionID,
			})
		}

		_, err = trx.Execute(`
			INSERT INTO tenants_billing (tenant_id, stripe_customer_id, stripe_subscription_id, license_key)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (tenant_id) DO UPDATE
			SET stripe_customer_id = $2, stripe_subscription_id = $3, license_key = $4, paddle_subscription_id = NULL
		`, c.TenantID, c.CustomerID, c.SubscriptionID, c.LicenseKey)
		if err != nil {
			return errors.Wrap(err, "failed to activate stripe subscription")
		}

		_, err = trx.Execute(`
			UPDATE tenants
			SET is_pro = true, status = $2
			WHERE id = $1
		`, c.TenantID, enum.TenantActive)
		if err != nil {
			return errors.Wrap(err, "failed to set tenant to pro plan")
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

		_, err = trx.Execute(`
			UPDATE tenants
			SET is_pro = false
			WHERE id = $1
		`, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to set tenant to free plan")
		}

		return nil
	})
}
