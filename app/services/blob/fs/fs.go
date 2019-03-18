package fs

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/pkg/bus"
)

var perm os.FileMode = 0744

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Category() string {
	return "blobstorage"
}

func (s Service) Enabled() bool {
	return env.Config.BlobStorage.Type == "fs"
}

func (s Service) Init() {
	bus.AddHandler(s, retrieveBlob)
	bus.AddHandler(s, storeBlob)
	bus.AddHandler(s, deleteBlob)
}

func retrieveBlob(ctx context.Context, cmd *blob.RetrieveBlob) error {
	fullPath := keyFullPath(ctx, cmd.Key)
	stats, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return blob.ErrNotFound
		}
		return errors.Wrap(err, "failed to get stats '%s' from FileSystem", cmd.Key)
	}

	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return errors.Wrap(err, "failed to get read '%s' from FileSystem", cmd.Key)
	}

	cmd.Blob = &blob.Blob{
		Content:     file,
		ContentType: http.DetectContentType(file),
		Size:        stats.Size(),
	}
	return nil
}

func storeBlob(ctx context.Context, cmd *blob.StoreBlob) error {
	if err := blob.ValidateKey(cmd.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", cmd.Key)
	}

	fullPath := keyFullPath(ctx, cmd.Key)
	err := os.MkdirAll(filepath.Dir(fullPath), perm)

	if err != nil {
		return errors.Wrap(err, "failed to create folder '%s' on FileSystem", fullPath)
	}

	err = ioutil.WriteFile(fullPath, cmd.Content, perm)
	if err != nil {
		return errors.Wrap(err, "failed to create file '%s' on FileSystem", fullPath)
	}

	return nil
}

func deleteBlob(ctx context.Context, cmd *blob.DeleteBlob) error {
	fullPath := keyFullPath(ctx, cmd.Key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete file '%s' from FileSystem", cmd.Key)
	}
	return nil
}

func keyFullPath(ctx context.Context, key string) string {
	basePath := env.Config.BlobStorage.FS.Path
	tenant, ok := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	if ok {
		return path.Join(basePath, "tenants", strconv.Itoa(tenant.ID), key)
	}
	return path.Join(basePath, key)
}
