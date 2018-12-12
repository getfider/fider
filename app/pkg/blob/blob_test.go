package blob_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/errors"

	"github.com/getfider/fider/app/pkg/rand"

	"github.com/getfider/fider/app/pkg/dbx"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/pkg/blob/fs"
	"github.com/getfider/fider/app/pkg/blob/s3"
	"github.com/getfider/fider/app/pkg/blob/sql"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/env"
)

var tenant1 = &models.Tenant{ID: 1}
var tenant2 = &models.Tenant{ID: 2}

func setupS3(t *testing.T) *s3.Storage {
	RegisterT(t)

	err := os.RemoveAll(env.Path("data/s3test/test-bucket"))
	Expect(err).IsNil()

	err = os.MkdirAll(env.Path("data/s3test/test-bucket"), 0777)
	Expect(err).IsNil()

	endpointURL := env.GetEnvOrDefault("S3_ENDPOINT_URL", "")
	region := env.GetEnvOrDefault("S3_REGION", "")
	accessKeyID := env.GetEnvOrDefault("S3_ACCESS_KEY_ID", "")
	secretAccessKey := env.GetEnvOrDefault("S3_SECRET_ACCESS_KEY", "")
	bucket := env.GetEnvOrDefault("S3_BUCKET", "")

	client, err := s3.NewStorage(endpointURL, region, accessKeyID, secretAccessKey, bucket)
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
	{"KeyFormats", KeyFormats},
}

func TestBlobStorage(t *testing.T) {
	for _, tt := range tests {
		t.Run("Test_FS_"+tt.name, func(t *testing.T) {
			client := setupFS(t)
			tt.test(client, t)
		})

		t.Run("Test_S3_"+tt.name, func(t *testing.T) {
			client := setupS3(t)
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

func SameKey_DifferentTenant_Delete(client blob.Storage, t *testing.T) {
	key := "path/to/super-file.txt"
	bytes1, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file.txt"))
	bytes2, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file3.txt"))

	sess1 := client.NewSession(tenant1)
	sess2 := client.NewSession(tenant2)

	err := sess1.Store(&blob.Blob{
		Key:         key,
		Object:      bytes1,
		ContentType: "text/plain; charset=utf-8",
		Size:        int64(len(bytes1)),
	})
	Expect(err).IsNil()

	err = sess2.Store(&blob.Blob{
		Key:         key,
		Object:      bytes2,
		ContentType: "text/plain; charset=utf-8",
		Size:        int64(len(bytes2)),
	})
	Expect(err).IsNil()

	b, err := sess1.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes1))

	b, err = sess2.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes2))

	err = sess1.Delete(key)
	Expect(err).IsNil()

	b, err = sess2.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes2))
	Expect(err).IsNil()

	sess3 := client.NewSession(nil)
	b, err = sess3.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)
}

func KeyFormats(client blob.Storage, t *testing.T) {
	RegisterT(t)
	sess := client.NewSession(nil)

	testCases := []struct {
		key   string
		valid bool
	}{
		{
			key:   "Jon",
			valid: true,
		},
		{
			key:   "ASDHASJDKHAJSDJ.png",
			valid: true,
		},
		{
			key:   "",
			valid: false,
		},
		{
			key:   "/path/to/Jon",
			valid: false,
		},
		{
			key:   "path/to/Jon/",
			valid: false,
		},
		{
			key:   " file.txt",
			valid: false,
		},
		{
			key:   "file with space.txt",
			valid: false,
		},
		{
			key:   rand.String(513),
			valid: false,
		},
	}

	for _, testCase := range testCases {
		err := sess.Store(&blob.Blob{
			Key:         testCase.key,
			Object:      make([]byte, 0),
			ContentType: "text/plain; charset=utf-8",
			Size:        0,
		})
		if testCase.valid {
			Expect(err).IsNil()
		} else {
			Expect(errors.Cause(err)).Equals(blob.ErrInvalidKeyFormat)
		}
	}
}
