package web

import (
	"context"
	"strings"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/services/blob"

	"golang.org/x/crypto/acme/autocert"
)

// Making sure that we're adhering to the autocert.Cache interface.
var _ autocert.Cache = (*Cache)(nil)

// Cache provides a Blob backend to the autocert cache.
type Cache struct {
	db *dbx.Database
}

// NewAutoCert returns a new AutoCert cache using Blob Storage
func NewAutoCert(db *dbx.Database) *Cache {
	return &Cache{db}
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	ctx = c.withDatabase(ctx)

	cmd := &blob.RetrieveBlob{
		Key: c.formatKey(key),
	}
	err := bus.Dispatch(ctx, cmd)
	if err == blob.ErrNotFound {
		return nil, autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return cmd.Blob.Content, nil
}

// Put stores the data in the cache under the specified key.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	ctx = c.withDatabase(ctx)

	return bus.Dispatch(ctx, &blob.StoreBlob{
		Key: c.formatKey(key),
		Blob: blob.Blob{
			Content:     data,
			ContentType: "application/x-pem-file",
			Size:        int64(len(data)),
		},
	})
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (c *Cache) Delete(ctx context.Context, key string) error {
	ctx = c.withDatabase(ctx)
	return bus.Dispatch(ctx, &blob.DeleteBlob{
		Key: c.formatKey(key),
	})
}

func (c *Cache) withDatabase(ctx context.Context) context.Context {
	return context.WithValue(ctx, app.DatabaseCtxKey, c.db)
}

func (c *Cache) formatKey(key string) string {
	return "autocert/" + strings.ToLower(key)
}
