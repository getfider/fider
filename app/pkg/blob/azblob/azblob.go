package azblob

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	az "github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/errors"
)

var _ blob.Storage = (*Storage)(nil)

// Storage stores blobs on an Azure Blob Container
type Storage struct {
	container az.ContainerURL
}

// Session is a per-request object to interact with the storage
type Session struct {
	storage *Storage
	tenant  *models.Tenant
}

func isNotFound(err error) bool {
	resp, ok := err.(az.StorageError)
	return ok && (resp.ServiceCode() == "BlobNotFound" || resp.Response().StatusCode == http.StatusNotFound)
}

// NewStorage creates a new Azure Blob Container storage
func NewStorage(endpointURL, accountName, accountKey, containerName string) (*Storage, error) {
	credential, err := az.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	p := az.NewPipeline(credential, az.PipelineOptions{})
	u, _ := url.Parse(endpointURL)
	service := az.NewServiceURL(*u, p)
	container := service.NewContainerURL(containerName)

	return &Storage{
		container,
	}, nil
}

// NewSession creates a new session
func (s *Storage) NewSession(tenant *models.Tenant) blob.Session {
	return &Session{
		storage: s,
		tenant:  tenant,
	}
}

func (s *Session) keyFullPathURL(key string) string {
	if s.tenant != nil {
		return path.Join("tenants", strconv.Itoa(s.tenant.ID), key)
	}
	return key
}

// Get returns a blob with given key
func (s *Session) Get(key string) (*blob.Blob, error) {
	blobURL := s.storage.container.NewBlockBlobURL(s.keyFullPathURL(key))
	get, err := blobURL.Download(context.Background(), 0, 0, az.BlobAccessConditions{}, false)
	if err != nil {
		if isNotFound(err) {
			return nil, blob.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get blob '%s' from Azure Blob Container", key)
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(az.RetryReaderOptions{})
	downloadedData.ReadFrom(reader)
	defer reader.Close()

	return &blob.Blob{
		Key:         key,
		Size:        get.ContentLength(),
		ContentType: get.ContentType(),
		Object:      downloadedData.Bytes(),
	}, nil
}

// Delete a blob with given key
func (s *Session) Delete(key string) error {
	blobURL := s.storage.container.NewBlockBlobURL(s.keyFullPathURL(key))
	_, err := blobURL.Delete(context.Background(), az.DeleteSnapshotsOptionNone, az.BlobAccessConditions{})
	if err != nil {
		if isNotFound(err) {
			return nil
		}
		return errors.Wrap(err, "failed to remove blob '%s' from Azure Blob Container", key)
	}
	return nil
}

// Store a blob with given key and content. Blobs with same key are replaced.
func (s *Session) Store(b *blob.Blob) error {
	blobURL := s.storage.container.NewBlockBlobURL(s.keyFullPathURL(b.Key))
	reader := bytes.NewReader(b.Object)
	_, err := blobURL.Upload(context.Background(), reader, az.BlobHTTPHeaders{ContentType: b.ContentType}, az.Metadata{}, az.BlobAccessConditions{})
	if err != nil {
		return errors.Wrap(err, "failed to upload blob '%s' to Azure Blob Container", b.Key)
	}
	return nil
}
