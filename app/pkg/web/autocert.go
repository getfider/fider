package web

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/blob"

	"golang.org/x/crypto/acme/autocert"
)

// Making sure that we're adhering to the autocert.Cache interface.
var _ autocert.Cache = (*Cache)(nil)

// Cache provides a Blob backend to the autocert cache.
type Cache struct {
}

// NewAutoCertCache returns a new AutoCert cache using Blob Storage
func NewAutoCertCache() *Cache {
	return &Cache{}
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	q := &query.GetBlobByKey{
		Key: c.formatKey(key),
	}
	err := bus.Dispatch(ctx, q)
	if errors.Cause(err) == blob.ErrNotFound {
		return nil, autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return q.Result.Content, nil
}

// Put stores the data in the cache under the specified key.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	return bus.Dispatch(ctx, &cmd.StoreBlob{
		Key:         c.formatKey(key),
		Content:     data,
		ContentType: "application/x-pem-file",
	})
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return bus.Dispatch(ctx, &cmd.DeleteBlob{
		Key: c.formatKey(key),
	})
}

func (c *Cache) formatKey(key string) string {
	return "autocert/" + strings.ToLower(key)
}
