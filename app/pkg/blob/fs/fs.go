package fs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/errors"
)

var perm os.FileMode = 0744

// Storage stores blobs on FileSystem
type Storage struct {
	path string
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

// Get returns a blob with given key
func (s *Storage) Get(key string) (*blob.Blob, error) {
	fullPath := path.Join(s.path, key)
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
	fullPath := path.Join(s.path, key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete file '%s' from FileSystem", key)
	}
	return nil
}

// Store a blob with given key and content. Blobs with same key are replaced.
func (s *Storage) Store(b *blob.Blob) error {
	fullPath := path.Join(s.path, b.Key)
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
