package fs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/errors"
)

var perm os.FileMode = 0744

var _ blob.Storage = (*Storage)(nil)

// Storage stores blobs on FileSystem
type Storage struct {
	path string
}

// Session is a per-request object to interact with the storage
type Session struct {
	storage *Storage
	tenant  *models.Tenant
}

// NewStorage creates a new FileSystem Storage
func NewStorage(path string) (*Storage, error) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create path '%s'", path)
	}

	return &Storage{
		path,
	}, nil
}

// NewSession creates a new session
func (s *Storage) NewSession(tenant *models.Tenant) blob.Session {
	return &Session{
		storage: s,
		tenant:  tenant,
	}
}

func (s *Session) keyFullPath(key string) string {
	fullPath := path.Join(s.storage.path, key)
	if s.tenant != nil {
		fullPath = path.Join("tenants", strconv.Itoa(s.tenant.ID), fullPath)
	}
	return fullPath
}

// Get returns a blob with given key
func (s *Session) Get(key string) (*blob.Blob, error) {
	fullPath := s.keyFullPath(key)
	stats, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, blob.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get stats '%s' from FileSystem", key)
	}

	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get read '%s' from FileSystem", key)
	}

	return &blob.Blob{
		Key:         key,
		Size:        stats.Size(),
		ContentType: http.DetectContentType(file),
		Object:      file,
	}, nil
}

// Delete a blob with given key
func (s *Session) Delete(key string) error {
	fullPath := s.keyFullPath(key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete file '%s' from FileSystem", key)
	}
	return nil
}

// Store a blob with given key and content. Blobs with same key are replaced.
func (s *Session) Store(b *blob.Blob) error {
	if err := blob.ValidateKey(b.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", b.Key)
	}

	fullPath := s.keyFullPath(b.Key)
	err := os.MkdirAll(filepath.Dir(fullPath), perm)

	if err != nil {
		return errors.Wrap(err, "failed to create folder '%s' on FileSystem", fullPath)
	}

	err = ioutil.WriteFile(fullPath, b.Object, perm)
	if err != nil {
		return errors.Wrap(err, "failed to create file '%s' on FileSystem", fullPath)
	}

	return nil
}
