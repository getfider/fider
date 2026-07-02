package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

// statusKinds enumerates the legal values of statuses.kind. Keep in sync with the
// CHECK constraint in migrations/202606231200_add_statuses_table.up.sql.
var statusKinds = map[string]bool{
	"open":             true,
	"active":           true,
	"closed-completed": true,
	"closed-declined":  true,
	"duplicate":        true,
}

// builtInStatusSeeds defines the six rows seeded for every new tenant. The
// "deleted" status doesn't appear here — it's an internal-only tombstone the
// app writes via the DeletePost handler.
var builtInStatusSeeds = []dbEntities.Status{
	{Slug: "open", Label: "Open", Kind: "open", Color: "blue", Icon: "lightbulb", ShowOnHome: false, ShowOnRoadmap: false, Filterable: true, SortOrder: 10, IsSystem: true, IsActive: true},
	{Slug: "planned", Label: "Planned", Kind: "active", Color: "blue", Icon: "thumbsup", ShowOnHome: true, ShowOnRoadmap: true, Filterable: true, SortOrder: 20, IsSystem: true, IsActive: true},
	{Slug: "started", Label: "Started", Kind: "active", Color: "blue", Icon: "sparkles-outline", ShowOnHome: true, ShowOnRoadmap: true, Filterable: true, SortOrder: 30, IsSystem: true, IsActive: true},
	{Slug: "completed", Label: "Completed", Kind: "closed-completed", Color: "green", Icon: "check-circle", ShowOnHome: true, ShowOnRoadmap: true, Filterable: true, SortOrder: 40, IsSystem: true, IsActive: true},
	{Slug: "declined", Label: "Declined", Kind: "closed-declined", Color: "red", Icon: "thumbsdown", ShowOnHome: true, ShowOnRoadmap: false, Filterable: true, SortOrder: 50, IsSystem: true, IsActive: true},
	{Slug: "duplicate", Label: "Duplicate", Kind: "duplicate", Color: "yellow", Icon: "duplicate", ShowOnHome: true, ShowOnRoadmap: false, Filterable: true, SortOrder: 60, IsSystem: true, IsActive: true},
}

const statusSelectCols = `id, tenant_id, slug, label, kind, color, icon, show_on_home, show_on_roadmap, filterable,
		sort_order, is_system, is_active, created_at, updated_at`

func listActiveStatusesForTenant(ctx context.Context, q *query.ListActiveStatusesForTenant) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		rows := []*dbEntities.Status{}
		err := trx.Select(&rows, `
			SELECT `+statusSelectCols+`
			FROM statuses
			WHERE tenant_id = $1 AND is_active = TRUE
			ORDER BY sort_order, id
		`, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to list active statuses for tenant %d", tenant.ID)
		}
		q.Result = make([]*entity.Status, len(rows))
		for i, r := range rows {
			q.Result[i] = r.ToModel()
		}
		return nil
	})
}

func getStatusByID(ctx context.Context, q *query.GetStatusByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		row := dbEntities.Status{}
		err := trx.Get(&row, `
			SELECT `+statusSelectCols+`
			FROM statuses
			WHERE tenant_id = $1 AND id = $2
		`, tenant.ID, q.ID)
		if err == sql.ErrNoRows {
			return app.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "failed to get status by id %d", q.ID)
		}
		q.Result = row.ToModel()
		return nil
	})
}

func getStatusBySlug(ctx context.Context, q *query.GetStatusBySlug) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		row := dbEntities.Status{}
		err := trx.Get(&row, `
			SELECT `+statusSelectCols+`
			FROM statuses
			WHERE tenant_id = $1 AND slug = $2
		`, tenant.ID, q.Slug)
		if err == sql.ErrNoRows {
			return app.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "failed to get status by slug %q", q.Slug)
		}
		q.Result = row.ToModel()
		return nil
	})
}

func countPostsByStatusID(ctx context.Context, q *query.CountPostsByStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		var slug string
		err := trx.Scalar(&slug, `SELECT slug FROM statuses WHERE tenant_id = $1 AND id = $2`, tenant.ID, q.StatusID)
		if err == sql.ErrNoRows {
			return app.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "failed to look up slug for status %d", q.StatusID)
		}
		// trx.Count counts the number of rows RETURNED by a query, not the
		// value of a COUNT(*) aggregate. Use Scalar to read the aggregate.
		var count int
		err = trx.Scalar(&count, `SELECT COUNT(*) FROM posts WHERE tenant_id = $1 AND status_slug = $2`, tenant.ID, slug)
		if err != nil {
			return errors.Wrap(err, "failed to count posts for status %d", q.StatusID)
		}
		q.Result = count
		return nil
	})
}

func createStatus(ctx context.Context, c *cmd.CreateStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		slug := strings.TrimSpace(strings.ToLower(c.Slug))
		if slug == "" {
			return errors.New("status slug is required")
		}
		if !statusKinds[c.Kind] {
			return errors.New("status kind %q is not allowed", c.Kind)
		}
		row := dbEntities.Status{}
		err := trx.Get(&row, `
			INSERT INTO statuses (tenant_id, slug, label, kind, color, icon, show_on_home, show_on_roadmap, filterable, sort_order, is_system, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE, TRUE)
			RETURNING `+statusSelectCols+`
		`, tenant.ID, slug, c.Label, c.Kind, c.Color, c.Icon, c.ShowOnHome, c.ShowOnRoadmap, c.Filterable, c.SortOrder)
		if err != nil {
			return errors.Wrap(err, "failed to create status")
		}
		c.Result = row.ToModel()
		return nil
	})
}

func updateStatus(ctx context.Context, c *cmd.UpdateStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		_, err := trx.Execute(`
			UPDATE statuses SET
				label = $1, color = $2, icon = $3,
				show_on_home = $4, show_on_roadmap = $5, filterable = $6, sort_order = $7, is_active = $8,
				updated_at = NOW()
			WHERE tenant_id = $9 AND id = $10
		`, c.Label, c.Color, c.Icon, c.ShowOnHome, c.ShowOnRoadmap, c.Filterable, c.SortOrder, c.IsActive, tenant.ID, c.ID)
		if err != nil {
			return errors.Wrap(err, "failed to update status %d", c.ID)
		}
		return nil
	})
}

func deleteStatus(ctx context.Context, c *cmd.DeleteStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, _ *entity.User) error {
		var isSystem bool
		err := trx.Scalar(&isSystem, `SELECT is_system FROM statuses WHERE tenant_id = $1 AND id = $2`, tenant.ID, c.ID)
		if err == sql.ErrNoRows {
			return app.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "failed to look up status %d", c.ID)
		}
		if isSystem {
			return errors.New("cannot delete a system status; deactivate it instead")
		}
		_, err = trx.Execute(`DELETE FROM statuses WHERE tenant_id = $1 AND id = $2`, tenant.ID, c.ID)
		if err != nil {
			return errors.Wrap(err, "failed to delete status %d", c.ID)
		}
		return nil
	})
}

func seedTenantStatuses(ctx context.Context, c *cmd.SeedTenantStatuses) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		for _, seed := range builtInStatusSeeds {
			_, err := trx.Execute(`
				INSERT INTO statuses (tenant_id, slug, label, kind, color, icon, show_on_home, show_on_roadmap, filterable, sort_order, is_system, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				ON CONFLICT (tenant_id, slug) DO NOTHING
			`, c.TenantID, seed.Slug, seed.Label, seed.Kind, seed.Color, seed.Icon, seed.ShowOnHome, seed.ShowOnRoadmap, seed.Filterable, seed.SortOrder, seed.IsSystem, seed.IsActive)
			if err != nil {
				return errors.Wrap(err, "failed to seed status %q for tenant %d", seed.Slug, c.TenantID)
			}
		}
		return nil
	})
}
