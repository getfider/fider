package blob

import (
	"errors"
	"strings"

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

// ErrInvalidKeyFormat is returned when blob key is in invalid format
var ErrInvalidKeyFormat = errors.New("Blob key is in invalid format")

// ValidateKey checks if key is is valid format
func ValidateKey(key string) error {
	if len(key) == 0 || len(key) > 512 || strings.Contains(key, " ") {
		return ErrInvalidKeyFormat
	}
	return nil
}

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
