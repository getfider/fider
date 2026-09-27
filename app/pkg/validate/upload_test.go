package validate_test

import (
	"bytes"
	"context"
	"image"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestValidateImageUpload(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		fileName string
		count    int
	}{
		{"/app/pkg/web/testdata/logo1.png", 0},
		{"/app/pkg/web/testdata/logo2.jpg", 2},
		{"/app/pkg/web/testdata/logo3.gif", 1},
		{"/app/pkg/web/testdata/logo4.png", 1},
		{"/app/pkg/web/testdata/logo5.png", 0},
		{"/README.md", 1},
		{"/app/pkg/web/testdata/favicon.ico", 1},
	}

	for _, testCase := range testCases {
		img, _ := os.ReadFile(env.Path(testCase.fileName))

		upload := &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		}
		messages, err := validate.ImageUpload(context.Background(), upload, validate.ImageUploadOpts{
			MinHeight:    200,
			MinWidth:     200,
			MaxKilobytes: 100,
			ExactRatio:   true,
		})
		Expect(messages).HasLen(testCase.count)
		Expect(err).IsNil()
	}
}

func TestValidateImageUpload_ExactRatio(t *testing.T) {
	RegisterT(t)

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo3-200w.gif"))
	opts := validate.ImageUploadOpts{
		IsRequired:   false,
		MaxKilobytes: 200,
	}

	upload := &dto.ImageUpload{
		Upload: &dto.ImageUploadData{
			Content: img,
		},
	}
	opts.ExactRatio = true
	messages, err := validate.ImageUpload(context.Background(), upload, opts)
	Expect(messages).HasLen(1)
	Expect(err).IsNil()

	opts.ExactRatio = false
	messages, err = validate.ImageUpload(context.Background(), upload, opts)
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}

func TestValidateImageUpload_Nil(t *testing.T) {
	RegisterT(t)

	messages, err := validate.ImageUpload(context.Background(), nil, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    200,
		MinWidth:     200,
		MaxKilobytes: 50,
		ExactRatio:   true,
	})
	Expect(messages).HasLen(0)
	Expect(err).IsNil()

	messages, err = validate.ImageUpload(context.Background(), &dto.ImageUpload{}, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    200,
		MinWidth:     200,
		MaxKilobytes: 50,
		ExactRatio:   true,
	})
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}

func TestValidateImageUpload_Required(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		upload *dto.ImageUpload
		count  int
	}{
		{nil, 1},
		{&dto.ImageUpload{}, 1},
		{&dto.ImageUpload{
			BlobKey: "some-file.png",
			Remove:  true,
		}, 1},
	}

	for _, testCase := range testCases {
		messages, err := validate.ImageUpload(context.Background(), testCase.upload, validate.ImageUploadOpts{
			IsRequired:   true,
			MinHeight:    200,
			MinWidth:     200,
			MaxKilobytes: 50,
			ExactRatio:   true,
		})
		Expect(messages).HasLen(testCase.count)
		Expect(err).IsNil()
	}
}

func TestValidateMultiImageUpload(t *testing.T) {
	RegisterT(t)

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo3-200w.gif"))

	uploads := []*dto.ImageUpload{
		{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		},
		{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		},
		{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		},
	}

	messages, err := validate.MultiImageUpload(context.Background(), nil, uploads, validate.MultiImageUploadOpts{
		MaxUploads:   2,
		MaxKilobytes: 500,
	})
	Expect(messages).HasLen(1)
	Expect(err).IsNil()
}

func TestValidateMultiImageUpload_Existing(t *testing.T) {
	RegisterT(t)

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo3-200w.gif"))

	uploads := []*dto.ImageUpload{
		{
			BlobKey: "attachments/file1.png",
			Remove:  true,
		},
		{
			BlobKey: "attachments/file2.png",
			Remove:  true,
		},
		{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		},
		{
			Upload: &dto.ImageUploadData{
				Content: img,
			},
		},
	}

	currentAttachments := []string{"attachments/file1.png", "attachments/file2.png"}
	messages, err := validate.MultiImageUpload(context.Background(), currentAttachments, uploads, validate.MultiImageUploadOpts{
		MaxUploads:   2,
		MaxKilobytes: 500,
	})
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}

func imageDimensions(t *testing.T, content []byte) (int, int) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(content))
	Expect(err).IsNil()
	return cfg.Width, cfg.Height
}

func TestValidateImageUpload_DecompressionBomb(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		name    string
		content []byte
	}{
		// ~12000x12000 = 144MP, decodes to ~576MB of NRGBA but is only a few KB compressed
		{"png 12000x12000", mock.UniformPNG(12000, 12000)},
		// 40000x40000 = 1.6 gigapixels
		{"gif 40000x40000", mock.GIFHeader(40000, 40000)},
		// Only one dimension is huge, but the pixel count is still way over the budget
		{"png 60000x1000", mock.UniformPNG(60000, 1000)},
	}

	for _, testCase := range testCases {
		Expect(len(testCase.content) < 100*1024).IsTrue()
		original := append([]byte{}, testCase.content...)

		upload := &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content: testCase.content,
			},
		}

		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		messages, err := validate.ImageUpload(context.Background(), upload, validate.ImageUploadOpts{
			MaxKilobytes: 5120,
		})
		runtime.ReadMemStats(&after)

		Expect(err).IsNil()
		Expect(messages).HasLen(1)
		Expect(messages[0]).Equals("The image dimensions are too large. The maximum is 40 megapixels.")
		Expect(upload.Upload.Content).Equals(original)

		// Must not have decoded the image (which would allocate hundreds of MB)
		allocated := after.TotalAlloc - before.TotalAlloc
		Expect(allocated < 10*1024*1024).IsTrue()
	}
}

func TestValidateImageUpload_TooManyBytes_DoesNotResize(t *testing.T) {
	RegisterT(t)

	content := mock.UniformPNG(2000, 1000)
	upload := &dto.ImageUpload{
		Upload: &dto.ImageUploadData{
			Content: content,
		},
	}

	messages, err := validate.ImageUpload(context.Background(), upload, validate.ImageUploadOpts{
		MaxKilobytes: 0,
	})
	Expect(err).IsNil()
	Expect(messages).HasLen(1)

	width, height := imageDimensions(t, upload.Upload.Content)
	Expect(width).Equals(2000)
	Expect(height).Equals(1000)
}

func TestValidateImageUpload_Resize(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{2000, 1000, 1500, 750},
		{1000, 2000, 750, 1500},
		// Only one dimension is over the limit
		{3000, 500, 1500, 250},
		{500, 3000, 250, 1500},
		// Within limits, not resized (and never upscaled)
		{1500, 1500, 1500, 1500},
		{800, 600, 800, 600},
	}

	for _, testCase := range testCases {
		upload := &dto.ImageUpload{
			Upload: &dto.ImageUploadData{
				Content: mock.UniformPNG(testCase.width, testCase.height),
			},
		}

		messages, err := validate.ImageUpload(context.Background(), upload, validate.ImageUploadOpts{
			MaxKilobytes: 5120,
		})
		Expect(err).IsNil()
		Expect(messages).HasLen(0)

		width, height := imageDimensions(t, upload.Upload.Content)
		Expect(width).Equals(testCase.expectedWidth)
		Expect(height).Equals(testCase.expectedHeight)
	}
}

func TestIsWithinPixelBudget(t *testing.T) {
	RegisterT(t)

	Expect(validate.IsWithinPixelBudget(1, 1)).IsTrue()
	Expect(validate.IsWithinPixelBudget(8000, 5000)).IsTrue()
	Expect(validate.IsWithinPixelBudget(40_000_000, 1)).IsTrue()
	Expect(validate.IsWithinPixelBudget(8000, 5001)).IsFalse()
	Expect(validate.IsWithinPixelBudget(12000, 12000)).IsFalse()
	Expect(validate.IsWithinPixelBudget(0, 100)).IsFalse()
	Expect(validate.IsWithinPixelBudget(100, 0)).IsFalse()
	Expect(validate.IsWithinPixelBudget(-100, -100)).IsFalse()
	Expect(validate.IsWithinPixelBudget(math.MaxInt32, math.MaxInt32)).IsFalse()
	Expect(validate.IsWithinPixelBudget(math.MaxInt, math.MaxInt)).IsFalse()
}

func TestImageWithinPixelBudget(t *testing.T) {
	RegisterT(t)

	ok, err := validate.ImageWithinPixelBudget(mock.UniformPNG(2000, 1000))
	Expect(err).IsNil()
	Expect(ok).IsTrue()

	ok, err = validate.ImageWithinPixelBudget(mock.UniformPNG(12000, 12000))
	Expect(err).IsNil()
	Expect(ok).IsFalse()

	ok, err = validate.ImageWithinPixelBudget(mock.GIFHeader(40000, 40000))
	Expect(err).IsNil()
	Expect(ok).IsFalse()

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo2.jpg"))
	ok, err = validate.ImageWithinPixelBudget(img)
	Expect(err).IsNil()
	Expect(ok).IsTrue()

	_, err = validate.ImageWithinPixelBudget([]byte("not an image"))
	Expect(err).IsNotNil()
}
