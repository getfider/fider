package blob_test

import (
	"context"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"testing"

	az "github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/getfider/fider/app/pkg/blob/azblob"
	"github.com/getfider/fider/app/pkg/blob/fs"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/env"
)

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
	container.Create(context.Background(), az.Metadata{}, az.PublicAccessNone)

	client, err := azblob.NewStorage(endpointURL, accountName, accountKey, containerName)
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
	}
}

func AllOperations(client blob.Storage, t *testing.T) {
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
		err := client.Store(&blob.Blob{
			Key:         testCase.key,
			Object:      bytes,
			ContentType: testCase.contentType,
			Size:        int64(len(bytes)),
		})
		Expect(err).IsNil()

		b, err := client.Get(testCase.key)
		Expect(err).IsNil()
		Expect(b.Key).Equals(testCase.key)
		Expect(b.Object).Equals(bytes)
		Expect(b.Size).Equals(int64(len(bytes)))
		Expect(b.ContentType).Equals(testCase.contentType)

		err = client.Delete(testCase.key)
		Expect(err).IsNil()

		b, err = client.Get(testCase.key)
		Expect(b).IsNil()
		Expect(err).Equals(blob.ErrNotFound)
	}
}

func DeleteUnkownFile(client blob.Storage, t *testing.T) {
	err := client.Delete("path/somefile.txt")
	Expect(err).IsNil()
}
