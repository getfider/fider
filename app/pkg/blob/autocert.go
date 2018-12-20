package blob

import (
	"context"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

// Making sure that we're adhering to the autocert.Cache interface.
var _ autocert.Cache = (*Cache)(nil)

// Cache provides a Blob backend to the autocert cache.
type Cache struct {
	storage Storage
}

// NewAutoCert returns a new AutoCert cache using Blob Storage
func NewAutoCert(storage Storage) *Cache {
	return &Cache{storage}
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := c.storage.Get(c.formatKey(key))
	if err == ErrNotFound {
		return nil, autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return b.Object, nil
}

// Put stores the data in the cache under the specified key.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	err := c.storage.Put(c.formatKey(key), data, "application/x-pem-file")
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (c *Cache) Delete(ctx context.Context, key string) error {
	err := c.storage.Delete(c.formatKey(key))
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) formatKey(key string) string {
	return "autocert/" + strings.ToLower(key)
}
