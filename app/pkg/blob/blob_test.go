package blob_test

import (
	"context"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"

	"github.com/getfider/fider/app/models"

	az "github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/getfider/fider/app/pkg/blob/azblob"
	"github.com/getfider/fider/app/pkg/blob/fs"
	"github.com/getfider/fider/app/pkg/blob/sql"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/env"
)

var tenant1 = &models.Tenant{ID: 1}
var tenant2 = &models.Tenant{ID: 2}

func setupAZBLOB(t *testing.T) *azblob.Storage {
	RegisterT(t)

	endpointURL := env.GetEnvOrDefault("AZBLOB_ENDPOINT_URL", "")
	accountName := env.GetEnvOrDefault("AZBLOB_ACCOUNT_NAME", "")
	accountKey := env.GetEnvOrDefault("AZBLOB_ACCOUNT_KEY", "")
	containerName := env.GetEnvOrDefault("AZBLOB_CONTAINER", "")

	credential, err := az.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	p := az.NewPipeline(credential, az.PipelineOptions{})
	u, _ := url.Parse(endpointURL)
	service := az.NewServiceURL(*u, p)
	container := service.NewContainerURL(containerName)
	container.Delete(context.Background(), az.ContainerAccessConditions{})
	container.Create(context.Background(), az.Metadata{}, az.PublicAccessNone)

	client, err := azblob.NewStorage(endpointURL, accountName, accountKey, containerName)
	Expect(err).IsNil()

	return client
}

func setupSQL(t *testing.T) *sql.Storage {
	RegisterT(t)

	db := dbx.New()
	db.Seed()
	client, err := sql.NewStorage(db)
	Expect(err).IsNil()

	return client
}

func setupFS(t *testing.T) *fs.Storage {
	RegisterT(t)

	rootFolder := env.Path("tmp/fs_test")

	err := os.RemoveAll(rootFolder)
	Expect(err).IsNil()

	client, err := fs.NewStorage(rootFolder)
	Expect(err).IsNil()

	return client
}

type blobTestCase func(client blob.Storage, t *testing.T)
type setupFunction func(t *testing.T) blob.Storage

var tests = []struct {
	name string
	test blobTestCase
}{
	{"AllOperations", AllOperations},
	{"DeleteUnkownFile", DeleteUnkownFile},
	{"SameKey_DifferentTenant", SameKey_DifferentTenant},
}

func TestBlobStorage(t *testing.T) {
	for _, tt := range tests {
		t.Run("Test_FS_"+tt.name, func(t *testing.T) {
			client := setupFS(t)
			tt.test(client, t)
		})

		t.Run("Test_AZBLOB_"+tt.name, func(t *testing.T) {
			client := setupAZBLOB(t)
			tt.test(client, t)
		})

		t.Run("Test_SQL_"+tt.name, func(t *testing.T) {
			client := setupSQL(t)
			tt.test(client, t)
		})
	}
}

func AllOperations(client blob.Storage, t *testing.T) {
	sess := client.NewSession(nil)
	var testCases = []struct {
		localPath   string
		key         string
		contentType string
	}{
		{"/app/pkg/blob/testdata/file.txt", "some/path/to/file.txt", "text/plain; charset=utf-8"},
		{"/app/pkg/blob/testdata/file2.png", "file2.png", "image/png"},
	}

	for _, testCase := range testCases {
		bytes, _ := ioutil.ReadFile(env.Path(testCase.localPath))
		err := sess.Store(&blob.Blob{
			Key:         testCase.key,
			Object:      bytes,
			ContentType: testCase.contentType,
			Size:        int64(len(bytes)),
		})
		Expect(err).IsNil()

		b, err := sess.Get(testCase.key)
		Expect(err).IsNil()
		Expect(b.Key).Equals(testCase.key)
		Expect(b.Object).Equals(bytes)
		Expect(b.Size).Equals(int64(len(bytes)))
		Expect(b.ContentType).Equals(testCase.contentType)

		err = sess.Delete(testCase.key)
		Expect(err).IsNil()

		b, err = sess.Get(testCase.key)
		Expect(b).IsNil()
		Expect(err).Equals(blob.ErrNotFound)
	}
}

func DeleteUnkownFile(client blob.Storage, t *testing.T) {
	sess := client.NewSession(nil)
	err := sess.Delete("path/somefile.txt")
	Expect(err).IsNil()
}

func SameKey_DifferentTenant(client blob.Storage, t *testing.T) {
	sess := client.NewSession(tenant1)
	key := "path/to/file3.txt"
	bytes, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file3.txt"))

	err := sess.Store(&blob.Blob{
		Key:         key,
		Object:      bytes,
		ContentType: "text/plain; charset=utf-8",
		Size:        int64(len(bytes)),
	})
	Expect(err).IsNil()

	b, err := sess.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(bytes)

	sess2 := client.NewSession(tenant2)

	b, err = sess2.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)

	sess3 := client.NewSession(nil)

	b, err = sess3.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)
}
