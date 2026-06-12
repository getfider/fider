package postgres

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func scheduleTenantDeletion(ctx context.Context, c *cmd.ScheduleTenantDeletion) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants
			SET scheduled_deletion_at = $2, deletion_requested_by = $3, deletion_cancel_key = $4
			WHERE id = $1
		`, c.TenantID, c.ScheduledAt, c.RequestedByUserID, c.CancelKey)
		if err != nil {
			return errors.Wrap(err, "failed to schedule tenant deletion")
		}
		return nil
	})
}

func cancelTenantDeletion(ctx context.Context, c *cmd.CancelTenantDeletion) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE tenants
			SET scheduled_deletion_at = NULL, deletion_requested_by = NULL, deletion_cancel_key = NULL
			WHERE id = $1
		`, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to cancel tenant deletion")
		}
		return nil
	})
}

func getTenantsPendingDeletion(ctx context.Context, q *query.GetTenantsPendingDeletion) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		q.Result = []*entity.Tenant{}
		var tenants []*dbEntities.Tenant
		err := trx.Select(&tenants, `
			SELECT id, name, subdomain, scheduled_deletion_at
			FROM tenants
			WHERE scheduled_deletion_at IS NOT NULL AND scheduled_deletion_at <= now()
			ORDER BY scheduled_deletion_at ASC
		`)
		if err != nil {
			return errors.Wrap(err, "failed to get tenants pending deletion")
		}
		for _, t := range tenants {
			q.Result = append(q.Result, t.ToModel())
		}
		return nil
	})
}

func getTenantByCancelKey(ctx context.Context, q *query.GetTenantByCancelKey) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		tenant := dbEntities.Tenant{}
		err := trx.Get(&tenant, `
			SELECT id, name, subdomain, scheduled_deletion_at
			FROM tenants
			WHERE deletion_cancel_key = $1
		`, q.Key)
		if err != nil {
			return errors.Wrap(err, "failed to get tenant by cancel key")
		}
		q.Result = tenant.ToModel()
		return nil
	})
}

func getTenantOwner(ctx context.Context, q *query.GetTenantOwner) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		// Scan directly into a minimal row rather than going through queryUser/ToModel — the
		// latter builds an avatar URL via web.AssetsURL, which panics in background-job context
		// because there is no HTTP request to derive the host/scheme from.
		row := struct {
			ID     int    `db:"id"`
			Name   string `db:"name"`
			Email  string `db:"email"`
			Role   int    `db:"role"`
			Status int    `db:"status"`
		}{}
		err := trx.Get(&row, `
			SELECT id, name, email, role, status
			FROM users
			WHERE tenant_id = $1 AND role = $2 AND status = $3
			ORDER BY id ASC
			LIMIT 1
		`, q.TenantID, enum.RoleAdministrator, enum.UserActive)
		if err != nil {
			return errors.Wrap(err, "failed to get owner of tenant '%d'", q.TenantID)
		}
		q.Result = &entity.User{
			ID:     row.ID,
			Name:   row.Name,
			Email:  row.Email,
			Role:   enum.Role(row.Role),
			Status: enum.UserStatus(row.Status),
		}
		return nil
	})
}

// tenantScopedTables lists every table holding tenant-owned rows, ordered children → parents
// so deletes respect the composite (id, tenant_id) foreign keys (all NO ACTION, no CASCADE).
// `reactions` is handled separately first because it has no tenant_id column.
var tenantScopedTables = []string{
	"attachments",
	"mention_notifications",
	"notifications",
	"post_subscribers",
	"post_votes",
	"post_tags",
	"comments",
	"posts",
	"tags",
	"email_verifications",
	"user_providers",
	"user_settings",
	"webhooks",
	"events",
	"blobs",
	"oauth_providers",
	"tenant_providers",
	"users",
	"tenants_billing",
}

func deleteTenant(ctx context.Context, c *cmd.DeleteTenant) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		// reactions has no tenant_id; scope via the tenant's comments before comments/users go.
		if _, err := trx.Execute(
			"DELETE FROM reactions WHERE comment_id IN (SELECT id FROM comments WHERE tenant_id = $1)",
			c.TenantID,
		); err != nil {
			return errors.Wrap(err, "failed to delete reactions for tenant '%d'", c.TenantID)
		}

		for _, table := range tenantScopedTables {
			if _, err := trx.Execute(
				fmt.Sprintf("DELETE FROM %s WHERE tenant_id = $1", table),
				c.TenantID,
			); err != nil {
				return errors.Wrap(err, "failed to delete %s for tenant '%d'", table, c.TenantID)
			}
		}

		if _, err := trx.Execute("DELETE FROM tenants WHERE id = $1", c.TenantID); err != nil {
			return errors.Wrap(err, "failed to delete tenant '%d'", c.TenantID)
		}
		return nil
	})
}
