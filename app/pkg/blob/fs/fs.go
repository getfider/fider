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
	path   string
	tenant *models.Tenant
}

// NewStorage creates a new FileSystem Storage
func NewStorage(path string) *Storage {
	err := os.MkdirAll(path, perm)
	if err != nil {
		panic(errors.Wrap(err, "failed to create path '%s'", path))
	}

	return &Storage{
		path: path,
	}
}

// SetCurrentTenant to current context
func (s *Storage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

func (s *Storage) keyFullPath(key string) string {
	if s.tenant != nil {
		return path.Join(s.path, "tenants", strconv.Itoa(s.tenant.ID), key)
	}
	return path.Join(s.path, key)
}

// Get returns a blob with given key
func (s *Storage) Get(key string) (*blob.Blob, error) {
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
func (s *Storage) Delete(key string) error {
	fullPath := s.keyFullPath(key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete file '%s' from FileSystem", key)
	}
	return nil
}

// Put a blob with given key and content. Blobs with same key are replaced.
func (s *Storage) Put(key string, content []byte, contentType string) error {
	if err := blob.ValidateKey(key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", key)
	}

	fullPath := s.keyFullPath(key)
	err := os.MkdirAll(filepath.Dir(fullPath), perm)

	if err != nil {
		return errors.Wrap(err, "failed to create folder '%s' on FileSystem", fullPath)
	}

	err = ioutil.WriteFile(fullPath, content, perm)
	if err != nil {
		return errors.Wrap(err, "failed to create file '%s' on FileSystem", fullPath)
	}

	return nil
}
