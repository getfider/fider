package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/blob"
)

func readTestImage(name string) []byte {
	content, err := os.ReadFile(env.Path("/app/pkg/web/testdata/" + name))
	Expect(err).IsNil()
	return content
}

func TestUploadImage(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	bus.AddHandler(func(ctx context.Context, c *cmd.StoreBlob) error {
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetBlobByKey) error {
		return blob.ErrNotFound
	})

	uploadImage := &cmd.UploadImage{
		Image: &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content:     readTestImage("logo1.png"),
				ContentType: "image/png",
			},
		},
		Folder: "avatars",
	}
	err := bus.Dispatch(ctx, uploadImage)
	Expect(err).IsNil()
	Expect(uploadImage.Image.BlobKey).ContainsSubstring("avatars/")
	Expect(uploadImage.Image.BlobKey).HasLen(73)
}

func TestUploadImage_ContentTypeFromImageBytes(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	var storedContentType string
	bus.AddHandler(func(ctx context.Context, c *cmd.StoreBlob) error {
		storedContentType = c.ContentType
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetBlobByKey) error {
		return blob.ErrNotFound
	})

	testCases := []struct {
		file     string
		expected string
	}{
		{"logo1.png", "image/png"},
		{"logo2.jpg", "image/jpeg"},
		{"logo3.gif", "image/gif"},
	}

	for _, testCase := range testCases {
		uploadImage := &cmd.UploadImage{
			Image: &dto.ImageUpload{
				Upload: &dto.ImageUploadData{
					FileName:    "image.html",
					Content:     readTestImage(testCase.file),
					ContentType: "text/html",
				},
			},
			Folder: "attachments",
		}
		err := bus.Dispatch(ctx, uploadImage)
		Expect(err).IsNil()
		Expect(storedContentType).Equals(testCase.expected)
		Expect(uploadImage.Image.Upload.ContentType).Equals(testCase.expected)
	}
}

func TestUploadImage_NotAnImage(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	storeBlobCalled := false
	bus.AddHandler(func(ctx context.Context, c *cmd.StoreBlob) error {
		storeBlobCalled = true
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetBlobByKey) error {
		return blob.ErrNotFound
	})

	uploadImage := &cmd.UploadImage{
		Image: &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content:     []byte("<script>alert(1)</script>"),
				ContentType: "image/png",
			},
		},
		Folder: "attachments",
	}
	err := bus.Dispatch(ctx, uploadImage)
	Expect(err).IsNotNil()
	Expect(storeBlobCalled).IsFalse()
}

func TestUploadImage_NoContent(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	bus.AddHandler(func(ctx context.Context, c *cmd.StoreBlob) error {
		return nil
	})

	uploadImage := &cmd.UploadImage{
		Image: &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content: []byte(""),
			},
		},
		Folder: "avatars",
	}
	err := bus.Dispatch(ctx, uploadImage)
	Expect(err).IsNil()
	Expect(uploadImage.Image.BlobKey).Equals("")
}

func TestUploadMultipleImages(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	bus.AddHandler(func(ctx context.Context, c *cmd.StoreBlob) error {
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetBlobByKey) error {
		return blob.ErrNotFound
	})

	uploadImages := &cmd.UploadImages{
		Images: []*dto.ImageUpload{
			{
				Upload: &dto.ImageUploadData{
					Content:     readTestImage("logo1.png"),
					ContentType: "image/png",
				},
			},
			{
				Upload: &dto.ImageUploadData{
					Content:     readTestImage("logo2.jpg"),
					ContentType: "image/jpeg",
				},
			},
		},
		Folder: "avatars",
	}
	err := bus.Dispatch(ctx, uploadImages)
	Expect(err).IsNil()

	Expect(uploadImages.Images[0].BlobKey).ContainsSubstring("avatars/")
	Expect(uploadImages.Images[0].BlobKey).HasLen(73)

	Expect(uploadImages.Images[1].BlobKey).ContainsSubstring("avatars/")
	Expect(uploadImages.Images[1].BlobKey).HasLen(73)
}
