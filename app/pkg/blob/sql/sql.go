package sql

import (
	"database/sql"
	"os"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

var perm os.FileMode = 0744

var _ blob.Storage = (*Storage)(nil)

type dbBlob struct {
	ContentType string `db:"content_type"`
	Size        int64  `db:"size"`
	Content     []byte `db:"file"`
}

// Storage stores blobs on FileSystem
type Storage struct {
	db *dbx.Database
}

// Session is a per-request object to interact with the storage
type Session struct {
	storage *Storage
	tenant  *models.Tenant
}

// NewStorage creates a new FileSystem Storage
func NewStorage(db *dbx.Database) (*Storage, error) {
	return &Storage{
		db,
	}, nil
}

// NewSession creates a new session
func (s *Storage) NewSession(tenant *models.Tenant) blob.Session {
	return &Session{
		storage: s,
		tenant:  tenant,
	}
}

// Get returns a blob with given key
func (s *Session) Get(key string) (*blob.Blob, error) {
	trx, err := s.storage.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open transaction")
	}
	defer trx.Commit()

	var tenantID sql.NullInt64
	if s.tenant != nil {
		tenantID.Scan(s.tenant.ID)
	}

	b := dbBlob{}
	err = trx.Get(&b, "SELECT file, content_type, size FROM blobs WHERE key = $1 AND (tenant_id = $2 OR ($2 IS NULL AND tenant_id IS NULL))", key, tenantID)
	if err != nil {
		if err == app.ErrNotFound {
			return nil, blob.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get blob with key '%s'", key)
	}

	return &blob.Blob{
		Key:         key,
		Size:        b.Size,
		ContentType: b.ContentType,
		Object:      b.Content,
	}, nil
}

// Delete a blob with given key
func (s *Session) Delete(key string) error {
	trx, err := s.storage.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to open transaction")
	}

	var tenantID sql.NullInt64
	if s.tenant != nil {
		tenantID.Scan(s.tenant.ID)
	}

	_, err = trx.Execute("DELETE FROM blobs WHERE key = $1 AND (tenant_id = $2 OR ($2 IS NULL AND tenant_id IS NULL))", key, tenantID)
	if err != nil {
		return errors.Wrap(err, "failed to delete blob with key '%s'", key)
	}

	err = trx.Commit()
	if err != nil {
		defer trx.Rollback()
		return errors.Wrap(err, "failed to commit deletion of blob with key '%s'", key)
	}

	return nil
}

// Store a blob with given key and content. Blobs with same key are replaced.
func (s *Session) Store(b *blob.Blob) error {
	if err := blob.ValidateKey(b.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", b.Key)
	}

	trx, err := s.storage.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to open transaction")
	}

	var tenantID sql.NullInt64
	if s.tenant != nil {
		tenantID.Scan(s.tenant.ID)
	}

	now := time.Now()
	_, err = trx.Execute(`
	INSERT INTO blobs (tenant_id, key, size, content_type, file, created_at, modified_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (tenant_id, key)
	DO UPDATE SET size = $3, content_type = $4, file = $5, modified_at = $7
	`, tenantID, b.Key, b.Size, b.ContentType, b.Object, now, now)
	if err != nil {
		return errors.Wrap(err, "failed to store blob with key '%s'", b.Key)
	}

	err = trx.Commit()
	if err != nil {
		defer trx.Rollback()
		return errors.Wrap(err, "failed to commit store of blob with key '%s'", b.Key)
	}

	return nil
}
