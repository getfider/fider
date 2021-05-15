package sql

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/pkg/bus"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "SQL"
}

func (s Service) Category() string {
	return "blobstorage"
}

func (s Service) Enabled() bool {
	return env.Config.BlobStorage.Type == "sql"
}

func (s Service) Init() {
	bus.AddHandler(listBlobs)
	bus.AddHandler(getBlobByKey)
	bus.AddHandler(storeBlob)
	bus.AddHandler(deleteBlob)
}

type dbBlob struct {
	Key         string `db:"key"`
	ContentType string `db:"content_type"`
	Size        int64  `db:"size"`
	Content     []byte `db:"file"`
}

func listBlobs(ctx context.Context, q *query.ListBlobs) error {
	return using(ctx, func(tenant *entity.Tenant) error {
		trx, err := dbx.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open transaction")
		}
		defer trx.MustCommit()

		blobs := []*dbBlob{}
		err = trx.Select(&blobs, "SELECT key FROM blobs WHERE tenant_id = $1", tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed list blobs")
		}

		files := make([]string, len(blobs))
		for i, b := range blobs {
			files[i] = b.Key
		}

		sort.Strings(files)
		q.Result = files
		return nil
	})
}

func getBlobByKey(ctx context.Context, q *query.GetBlobByKey) error {
	return using(ctx, func(tenant *entity.Tenant) error {
		var tenantID sql.NullInt64
		if tenant != nil {
			err := tenantID.Scan(tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed scan tenant id")
			}
		}

		trx, err := dbx.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open transaction")
		}
		defer trx.MustCommit()

		b := dbBlob{}
		err = trx.Get(&b, "SELECT file, content_type, size FROM blobs WHERE key = $1 AND (tenant_id = $2 OR ($2 IS NULL AND tenant_id IS NULL))", q.Key, tenantID)
		if err != nil {
			if err == app.ErrNotFound {
				return blob.ErrNotFound
			}
			return errors.Wrap(err, "failed to get blob with key '%s'", q.Key)
		}

		q.Result = &dto.Blob{
			Size:        b.Size,
			ContentType: b.ContentType,
			Content:     b.Content,
		}
		return nil
	})
}

func storeBlob(ctx context.Context, c *cmd.StoreBlob) error {
	if err := blob.ValidateKey(c.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", c.Key)
	}

	return using(ctx, func(tenant *entity.Tenant) error {
		var tenantID sql.NullInt64
		if tenant != nil {
			err := tenantID.Scan(tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed scan tenant id")
			}
		}

		trx, err := dbx.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open transaction")
		}
		defer trx.MustCommit()

		now := time.Now()
		_, err = trx.Execute(`
		INSERT INTO blobs (tenant_id, key, size, content_type, file, created_at, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (tenant_id, key)
		DO UPDATE SET size = $3, content_type = $4, file = $5, modified_at = $7
		`, tenantID, c.Key, int64(len(c.Content)), c.ContentType, c.Content, now, now)
		if err != nil {
			return errors.Wrap(err, "failed to store blob with key '%s'", c.Key)
		}

		if err != nil {
			return errors.Wrap(err, "failed to commit store of blob with key '%s'", c.Key)
		}

		return nil
	})
}

func deleteBlob(ctx context.Context, c *cmd.DeleteBlob) error {
	return using(ctx, func(tenant *entity.Tenant) error {
		var tenantID sql.NullInt64
		if tenant != nil {
			err := tenantID.Scan(tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed scan tenant id")
			}
		}

		trx, err := dbx.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open transaction")
		}
		defer trx.MustCommit()

		_, err = trx.Execute("DELETE FROM blobs WHERE key = $1 AND (tenant_id = $2 OR ($2 IS NULL AND tenant_id IS NULL))", c.Key, tenantID)
		if err != nil {
			return errors.Wrap(err, "failed to delete blob with key '%s'", c.Key)
		}

		if err != nil {
			return errors.Wrap(err, "failed to commit deletion of blob with key '%s'", c.Key)
		}

		return nil
	})
}

func using(ctx context.Context, handler func(tenant *entity.Tenant) error) error {
	tenant, _ := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	return handler(tenant)
}
