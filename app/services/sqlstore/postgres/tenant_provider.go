package postgres

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbTenantProvider struct {
	ID        int  `db:"id"`
	TenantID  int  `db:"tenant_id"`
	Provider  string `db:"provider"`
	IsEnabled bool `db:"is_enabled"`
}

func (m *dbTenantProvider) toModel() *entity.TenantProvider {
	return &entity.TenantProvider{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Provider:  m.Provider,
		IsEnabled: m.IsEnabled,
	}
}

func getTenantProviderStatus(ctx context.Context, q *query.GetTenantProviderStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant == nil {
			return app.ErrNotFound
		}

		tp := &dbTenantProvider{}
		err := trx.Get(tp, `
			SELECT id, tenant_id, provider, is_enabled
			FROM tenant_providers
			WHERE tenant_id = $1 AND provider = $2
		`, tenant.ID, q.Provider)
		
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				// If no record exists, assume it's enabled (default behavior)
				q.Result = &entity.TenantProvider{
					TenantID:  tenant.ID,
					Provider:  q.Provider,
					IsEnabled: true,
				}
				return nil
			}
			return errors.Wrap(err, "failed to get tenant provider status")
		}

		q.Result = tp.toModel()
		return nil
	})
}

func setTenantProviderStatus(ctx context.Context, c *cmd.SetTenantProviderStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant == nil {
			return app.ErrNotFound
		}

		// Insert or update the tenant provider status
		query := `
			INSERT INTO tenant_providers (tenant_id, provider, is_enabled, created_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (tenant_id, provider)
			DO UPDATE SET is_enabled = $3
		`
		
		_, err := trx.Execute(query, tenant.ID, c.Provider, c.IsEnabled)
		if err != nil {
			return errors.Wrap(err, "failed to set tenant provider status")
		}

		return nil
	})
}
