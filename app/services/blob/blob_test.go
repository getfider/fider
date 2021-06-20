package blob_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/rand"

	"github.com/getfider/fider/app/services/blob/fs"
	"github.com/getfider/fider/app/services/blob/s3"
	"github.com/getfider/fider/app/services/blob/sql"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/blob"
)

var tenant1 = &entity.Tenant{ID: 1}
var tenant2 = &entity.Tenant{ID: 2}

func setupS3(t *testing.T) {
	RegisterT(t)

	bus.Init(s3.Service{})

	bucket := aws.String(env.Config.BlobStorage.S3.BucketName)
	_, _ = s3.DefaultClient.CreateBucket(&awss3.CreateBucketInput{
		Bucket: bucket,
	})
}

func setupSQL(t *testing.T) {
	RegisterT(t)

	dbx.Seed()

	bus.Init(sql.Service{})
}

func setupFS(t *testing.T) {
	RegisterT(t)

	bus.Init(fs.Service{})

	err := os.RemoveAll(env.Config.BlobStorage.FS.Path)
	Expect(err).IsNil()
}

type blobTestCase func(ctx context.Context)

var tests = []struct {
	name string
	test blobTestCase
}{
	{"AllOperations", AllOperations},
	{"DeleteUnkownFile", DeleteUnkownFile},
	{"KeyFormats", KeyFormats},
	{"SameKey_DifferentTenant", SameKey_DifferentTenant},
	{"SameKey_DifferentTenant_Delete", SameKey_DifferentTenant_Delete},
	{"ListBlobsFromTenant", ListBlobsFromTenant},
	{"ListBlobsOutsideTenant", ListBlobsOutsideTenant},
	{"ListUnauthorizedBlobs", ListUnauthorizedBlobs},
}

func TestBlobStorage(t *testing.T) {
	for _, tt := range tests {
		t.Run("Test_FS_"+tt.name, func(t *testing.T) {
			setupFS(t)
			tt.test(context.Background())
		})

		t.Run("Test_S3_"+tt.name, func(t *testing.T) {
			setupS3(t)
			tt.test(context.Background())
		})

		t.Run("Test_SQL_"+tt.name, func(t *testing.T) {
			setupSQL(t)
			tt.test(context.Background())
		})
	}
}

func AllOperations(ctx context.Context) {
	var testCases = []struct {
		localPath   string
		key         string
		contentType string
	}{
		{"/app/services/blob/testdata/file.txt", "some/path/to/file.txt", "text/plain; charset=utf-8"},
		{"/app/services/blob/testdata/file2.png", "file2.png", "image/png"},
	}

	for _, testCase := range testCases {
		bytes, _ := ioutil.ReadFile(env.Path(testCase.localPath))
		err := bus.Dispatch(ctx, &cmd.StoreBlob{
			Key:         testCase.key,
			Content:     bytes,
			ContentType: testCase.contentType,
		})
		Expect(err).IsNil()

		q := &query.GetBlobByKey{
			Key: testCase.key,
		}
		err = bus.Dispatch(ctx, q)
		Expect(err).IsNil()
		Expect(q.Key).Equals(testCase.key)
		Expect(q.Result.Size).Equals(int64(len(bytes)))
		Expect(q.Result.ContentType).Equals(testCase.contentType)
		Expect(q.Result.Content).Equals(bytes)

		err = bus.Dispatch(ctx, &cmd.DeleteBlob{
			Key: testCase.key,
		})
		Expect(err).IsNil()

		q = &query.GetBlobByKey{
			Key: testCase.key,
		}
		err = bus.Dispatch(ctx, q)
		Expect(q.Result).IsNil()
		Expect(err).Equals(blob.ErrNotFound)
	}
}

func DeleteUnkownFile(ctx context.Context) {
	err := bus.Dispatch(ctx, &cmd.DeleteBlob{
		Key: "path/somefile.txt",
	})
	Expect(err).IsNil()
}

func SameKey_DifferentTenant(ctx context.Context) {
	ctxWithTenant1 := context.WithValue(ctx, app.TenantCtxKey, tenant1)
	ctxWithTenant2 := context.WithValue(ctx, app.TenantCtxKey, tenant2)

	key := "path/to/file3.txt"
	bytes, _ := ioutil.ReadFile(env.Path("/app/services/blob/testdata/file3.txt"))

	err := bus.Dispatch(ctxWithTenant1, &cmd.StoreBlob{
		Key:         key,
		Content:     bytes,
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	q := &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctxWithTenant1, q)
	Expect(err).IsNil()
	Expect(q.Result.Content).Equals(bytes)

	q = &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctxWithTenant2, q)
	Expect(err).Equals(blob.ErrNotFound)
	Expect(q.Result).IsNil()

	q = &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctx, q)
	Expect(err).Equals(blob.ErrNotFound)
	Expect(q.Result).IsNil()

	err = bus.Dispatch(ctxWithTenant1, &cmd.DeleteBlob{Key: key})
	Expect(err).IsNil()
}

func SameKey_DifferentTenant_Delete(ctx context.Context) {
	ctxWithTenant1 := context.WithValue(ctx, app.TenantCtxKey, tenant1)
	ctxWithTenant2 := context.WithValue(ctx, app.TenantCtxKey, tenant2)

	key := "path/to/super-file.txt"
	bytes1, _ := ioutil.ReadFile(env.Path("/app/services/blob/testdata/file.txt"))
	bytes2, _ := ioutil.ReadFile(env.Path("/app/services/blob/testdata/file3.txt"))

	err := bus.Dispatch(ctxWithTenant1, &cmd.StoreBlob{
		Key:         key,
		Content:     bytes1,
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	err = bus.Dispatch(ctxWithTenant2, &cmd.StoreBlob{
		Key:         key,
		Content:     bytes2,
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	q := &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctxWithTenant1, q)
	Expect(err).IsNil()
	Expect(q.Result.Content).Equals(bytes1)

	q = &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctxWithTenant2, q)
	Expect(err).IsNil()
	Expect(q.Result.Content).Equals(bytes2)

	err = bus.Dispatch(ctxWithTenant1, &cmd.DeleteBlob{Key: key})
	Expect(err).IsNil()

	q = &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctxWithTenant2, q)
	Expect(err).IsNil()
	Expect(q.Result.Content).Equals(bytes2)

	q = &query.GetBlobByKey{Key: key}
	err = bus.Dispatch(ctx, q)
	Expect(err).Equals(blob.ErrNotFound)
	Expect(q.Result).IsNil()

	err = bus.Dispatch(ctxWithTenant2, &cmd.DeleteBlob{Key: key})
	Expect(err).IsNil()
}

func ListBlobsFromTenant(ctx context.Context) {
	ctxWithTenant1 := context.WithValue(ctx, app.TenantCtxKey, tenant1)
	ctxWithTenant2 := context.WithValue(ctx, app.TenantCtxKey, tenant2)

	err := bus.Dispatch(ctxWithTenant1, &cmd.StoreBlob{
		Key:         "texts/hello.txt",
		Content:     make([]byte, 0),
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	err = bus.Dispatch(ctxWithTenant1, &cmd.StoreBlob{
		Key:         "memos/hello1.txt",
		Content:     make([]byte, 0),
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	err = bus.Dispatch(ctxWithTenant2, &cmd.StoreBlob{
		Key:         "texts/hello.txt",
		Content:     make([]byte, 0),
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	tenant1Files := &query.ListBlobs{}
	err = bus.Dispatch(ctxWithTenant1, tenant1Files)
	Expect(err).IsNil()
	Expect(tenant1Files.Result).HasLen(2)
	Expect(tenant1Files.Result).Equals([]string{"memos/hello1.txt", "texts/hello.txt"})

	tenant2Files := &query.ListBlobs{}
	err = bus.Dispatch(ctxWithTenant2, tenant2Files)
	Expect(err).IsNil()
	Expect(tenant2Files.Result).HasLen(1)
	Expect(tenant2Files.Result).Equals([]string{"texts/hello.txt"})
}

func ListBlobsOutsideTenant(ctx context.Context) {
	err := bus.Dispatch(ctx, &cmd.StoreBlob{
		Key:         "texts/hello.txt",
		Content:     make([]byte, 0),
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	err = bus.Dispatch(ctx, &cmd.StoreBlob{
		Key:         "texts/world.txt",
		Content:     make([]byte, 0),
		ContentType: "text/plain; charset=utf-8",
	})
	Expect(err).IsNil()

	textFiles := &query.ListBlobs{Prefix: "texts/"}
	err = bus.Dispatch(ctx, textFiles)
	Expect(err).IsNil()
	Expect(textFiles.Result).HasLen(2)
	Expect(textFiles.Result).Equals([]string{"texts/hello.txt", "texts/world.txt"})

	imageFiles := &query.ListBlobs{Prefix: "images/"}
	err = bus.Dispatch(ctx, imageFiles)
	Expect(err).IsNil()
	Expect(imageFiles.Result).HasLen(0)
	Expect(imageFiles.Result).Equals([]string{})
}

func ListUnauthorizedBlobs(ctx context.Context) {
	Expect(func() {
		_ = bus.Dispatch(ctx, &query.ListBlobs{Prefix: "tenants/"})
	}).Panics()

	Expect(func() {
		_ = bus.Dispatch(ctx, &query.ListBlobs{Prefix: "tenants"})
	}).Panics()

	ctxWithTenant := context.WithValue(ctx, app.TenantCtxKey, tenant1)
	err := bus.Dispatch(ctxWithTenant, &query.ListBlobs{Prefix: "tenants"})
	Expect(err).IsNil()
}

func KeyFormats(ctx context.Context) {
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
		err := bus.Dispatch(ctx, &cmd.StoreBlob{
			Key:         testCase.key,
			Content:     []byte{},
			ContentType: "text/plain; charset=utf-8",
		})
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
