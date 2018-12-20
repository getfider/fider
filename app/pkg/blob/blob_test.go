package blob_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
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

	bucket := "test-bucket"
	s3.DefaultClient.DeleteBucket(&awss3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	s3.DefaultClient.CreateBucket(&awss3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	return s3.NewStorage(bucket)
}

func setupSQL(t *testing.T) *sql.Storage {
	RegisterT(t)

	db := dbx.New()
	db.Seed()
	return sql.NewStorage(db)
}

func setupFS(t *testing.T) *fs.Storage {
	RegisterT(t)

	err := os.RemoveAll("./tmp/fs_test")
	Expect(err).IsNil()

	return fs.NewStorage("./tmp/fs_test")
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
	sqlClient := setupSQL(t)
	s3Client := setupS3(t)
	fsClient := setupFS(t)

	for _, tt := range tests {
		t.Run("Test_FS_"+tt.name, func(t *testing.T) {
			tt.test(fsClient, t)
		})

		t.Run("Test_S3_"+tt.name, func(t *testing.T) {
			tt.test(s3Client, t)
		})

		t.Run("Test_SQL_"+tt.name, func(t *testing.T) {
			tt.test(sqlClient, t)
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
		err := client.Put(testCase.key, bytes, testCase.contentType)
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

func SameKey_DifferentTenant(client blob.Storage, t *testing.T) {
	client.SetCurrentTenant(tenant1)
	key := "path/to/file3.txt"
	bytes, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file3.txt"))

	err := client.Put(key, bytes, "text/plain; charset=utf-8")
	Expect(err).IsNil()

	b, err := client.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(bytes)

	client.SetCurrentTenant(tenant2)

	b, err = client.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)

	client.SetCurrentTenant(nil)

	b, err = client.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)
}

func SameKey_DifferentTenant_Delete(client blob.Storage, t *testing.T) {
	key := "path/to/super-file.txt"
	bytes1, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file.txt"))
	bytes2, _ := ioutil.ReadFile(env.Path("/app/pkg/blob/testdata/file3.txt"))

	client.SetCurrentTenant(tenant1)
	err := client.Put(key, bytes1, "text/plain; charset=utf-8")
	Expect(err).IsNil()

	client.SetCurrentTenant(tenant2)
	err = client.Put(key, bytes2, "text/plain; charset=utf-8")
	Expect(err).IsNil()

	client.SetCurrentTenant(tenant1)
	b, err := client.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes1))

	client.SetCurrentTenant(tenant2)
	b, err = client.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes2))

	client.SetCurrentTenant(tenant1)
	err = client.Delete(key)
	Expect(err).IsNil()

	client.SetCurrentTenant(tenant2)
	b, err = client.Get(key)
	Expect(err).IsNil()
	Expect(b.Object).Equals(len(bytes2))
	Expect(err).IsNil()

	client.SetCurrentTenant(nil)
	b, err = client.Get(key)
	Expect(b).IsNil()
	Expect(err).Equals(blob.ErrNotFound)
}

func KeyFormats(client blob.Storage, t *testing.T) {
	RegisterT(t)

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
		err := client.Put(testCase.key, make([]byte, 0), "text/plain; charset=utf-8")
		if testCase.valid {
			Expect(err).IsNil()
		} else {
			Expect(errors.Cause(err)).Equals(blob.ErrInvalidKeyFormat)
		}
	}
}

func TestSanitizeFileName(t *testing.T) {
	RegisterT(t)

	Expect(blob.SanitizeFileName("João.txt")).Equals("joao.txt")
	Expect(blob.SanitizeFileName(" Jon")).Equals("jon")
	Expect(blob.SanitizeFileName("Jon.png")).Equals("jon.png")
	Expect(blob.SanitizeFileName("Jon Snow.png ")).Equals("jon-snow.png")
	Expect(blob.SanitizeFileName(" ヒキワリ.png")).Equals("hikiwari.png")
	Expect(blob.SanitizeFileName("люди рождаются свободными.png")).Equals("liudi-rozhdaiutsia-svobodnymi.png")
}
