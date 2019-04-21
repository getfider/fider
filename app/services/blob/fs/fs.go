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
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
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

func (s Service) Name() string {
	return "FileSystem"
}

func (s Service) Category() string {
	return "blobstorage"
}

func (s Service) Enabled() bool {
	return env.Config.BlobStorage.Type == "fs"
}

func (s Service) Init() {
	bus.AddHandler(getBlobByKey)
	bus.AddHandler(storeBlob)
	bus.AddHandler(deleteBlob)
}

func getBlobByKey(ctx context.Context, q *query.GetBlobByKey) error {
	fullPath := keyFullPath(ctx, q.Key)
	stats, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return blob.ErrNotFound
		}
		return errors.Wrap(err, "failed to get stats '%s' from FileSystem", q.Key)
	}

	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return errors.Wrap(err, "failed to get read '%s' from FileSystem", q.Key)
	}

	q.Result = &dto.Blob{
		Content:     file,
		ContentType: http.DetectContentType(file),
		Size:        stats.Size(),
	}
	return nil
}

func storeBlob(ctx context.Context, c *cmd.StoreBlob) error {
	if err := blob.ValidateKey(c.Key); err != nil {
		return errors.Wrap(err, "failed to validate blob key '%s'", c.Key)
	}

	fullPath := keyFullPath(ctx, c.Key)
	err := os.MkdirAll(filepath.Dir(fullPath), perm)

	if err != nil {
		return errors.Wrap(err, "failed to create folder '%s' on FileSystem", fullPath)
	}

	err = ioutil.WriteFile(fullPath, c.Content, perm)
	if err != nil {
		return errors.Wrap(err, "failed to create file '%s' on FileSystem", fullPath)
	}

	return nil
}

func deleteBlob(ctx context.Context, c *cmd.DeleteBlob) error {
	fullPath := keyFullPath(ctx, c.Key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to delete file '%s' from FileSystem", c.Key)
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
