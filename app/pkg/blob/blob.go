package blob

import (
	"errors"

	"github.com/getfider/fider/app/models"
)

// Blob is a file persisted somewhere
type Blob struct {
	Key         string
	Object      []byte
	Size        int64
	ContentType string
}

// ErrNotFound is returned when given blob is not found
var ErrNotFound = errors.New("Blob not found")

// Storage is how Fider persists blobs
type Storage interface {
	NewSession(tenant *models.Tenant) Session
}

// Session is a per-request object to interact with the storage
type Session interface {
	Get(key string) (*Blob, error)
	Delete(key string) error
	Store(b *Blob) error
}
