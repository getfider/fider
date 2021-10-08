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

func getSystemSettings(ctx context.Context, q *query.GetSystemSettings) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {

		var value string
		err := trx.Scalar(&value, `SELECT value FROM system_settings WHERE key = $1`, q.Key)
		if err != nil && err != app.ErrNotFound {
			return errors.Wrap(err, "failed to query system_settings with key '%s'", q.Key)
		}

		q.Value = value
		return nil
	})
}

func setSystemSettings(ctx context.Context, c *cmd.SetSystemSettings) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {

		_, err := trx.Execute(`
			INSERT INTO system_settings (key, value) VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;
		`, c.Key, c.Value)
		if err != nil {
			return errors.Wrap(err, "failed to update system_settings with key '%s'", c.Key)
		}

		return nil
	})
}
