package validate_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/goenning/imagic"
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

// jpegHeaderWithDimensions returns a JPEG whose header (SOF) claims the given dimensions.
// Only the header is valid, which is all image.DecodeConfig reads.
func jpegHeaderWithDimensions(t *testing.T, width, height int) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, image.NewYCbCr(image.Rect(0, 0, 16, 16), image.YCbCrSubsampleRatio444), nil)
	Expect(err).IsNil()
	content := buf.Bytes()
	sof := bytes.Index(content, []byte{0xFF, 0xC0})
	Expect(sof > 0).IsTrue()
	// FF C0, length (2), precision (1), height (2), width (2)
	binary.BigEndian.PutUint16(content[sof+5:], uint16(height))
	binary.BigEndian.PutUint16(content[sof+7:], uint16(width))
	return content
}

func TestValidateImageUpload_DecompressionBomb(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		name    string
		content []byte
	}{
		// 12000x12000 = 144MP, decodes to hundreds of MB but is only a few KB compressed
		{"png 12000x12000", mock.UniformPNG(12000, 12000)},
		// 40000x40000 = 1.6 gigapixels
		{"gif 40000x40000", mock.GIFHeader(40000, 40000)},
		// Small pixel count, but extreme aspect ratio
		{"png 100000x1", mock.UniformPNG(100000, 1)},
		// 4000x4000 colour JPEG is over the budget, as it might be progressive
		{"jpeg 4000x4000", jpegHeaderWithDimensions(t, 4000, 4000)},
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
		Expect(messages[0]).Equals("The image dimensions are too large. Please upload an image with a lower resolution.")
		Expect(upload.Upload.Content).Equals(original)

		// Must not have decoded the image (which would allocate hundreds of MB)
		allocated := after.TotalAlloc - before.TotalAlloc
		Expect(allocated < 10*1024*1024).IsTrue()
	}
}

func TestValidateImageUpload_TooManyBytes_DoesNotResize(t *testing.T) {
	RegisterT(t)

	content := mock.UniformPNG(3000, 2000)
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
	Expect(width).Equals(3000)
	Expect(height).Equals(2000)
}

func TestValidateImageUpload_Resize(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		// Both dimensions over the limit: resized, keeping the aspect ratio
		{3000, 2000, 1500, 1000},
		{2000, 3000, 1000, 1500},
		// Only one dimension over the limit: not resized (keeps GIF animations / JPEG EXIF)
		{2000, 1000, 2000, 1000},
		{3000, 500, 3000, 500},
		// Within limits, not resized
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

func TestEstimateDecodeBytes(t *testing.T) {
	RegisterT(t)

	// Invalid or over the per-side cap
	Expect(validate.EstimateDecodeBytes(0, 100, "png", color.GrayModel)).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(100, 0, "png", color.GrayModel)).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(-100, -100, "png", color.GrayModel)).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(validate.MaxImageSide+1, 1, "png", color.GrayModel)).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(1, validate.MaxImageSide+1, "png", color.GrayModel)).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(40_000_000, 1, "gif", color.Palette{})).Equals(int64(-1))
	Expect(validate.EstimateDecodeBytes(math.MaxInt, math.MaxInt, "png", color.GrayModel)).Equals(int64(-1))

	// Bytes per pixel include 4 bytes for a full-size RGBA intermediate
	const mp = 1000 * 1000
	Expect(validate.EstimateDecodeBytes(1000, 1000, "gif", color.Palette{})).Equals(int64(5 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "png", color.Palette{})).Equals(int64(6 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "png", color.GrayModel)).Equals(int64(6 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "png", color.Gray16Model)).Equals(int64(8 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "png", color.NRGBAModel)).Equals(int64(12 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "png", color.NRGBA64Model)).Equals(int64(20 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "jpeg", color.GrayModel)).Equals(int64(9 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "jpeg", color.YCbCrModel)).Equals(int64(19 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "jpeg", color.RGBAModel)).Equals(int64(23 * mp))
	Expect(validate.EstimateDecodeBytes(1000, 1000, "jpeg", color.CMYKModel)).Equals(int64(24 * mp))

	// Worst case at the per-side cap doesn't overflow
	Expect(validate.EstimateDecodeBytes(validate.MaxImageSide, validate.MaxImageSide, "jpeg", color.CMYKModel) > 0).IsTrue()
}

func TestCheckDecodeBudget(t *testing.T) {
	RegisterT(t)

	Expect(validate.CheckDecodeBudget(mock.UniformPNG(2000, 1000))).IsNil()
	// 36MP * 6 bytes = 216MB
	Expect(validate.CheckDecodeBudget(mock.UniformPNG(6000, 6000))).IsNil()
	// 49MP * 6 bytes = 294MB
	Expect(validate.CheckDecodeBudget(mock.UniformPNG(7000, 7000))).Equals(validate.ErrImageTooLarge)
	Expect(validate.CheckDecodeBudget(mock.UniformPNG(12000, 12000))).Equals(validate.ErrImageTooLarge)
	Expect(validate.CheckDecodeBudget(mock.UniformPNG(100000, 1))).Equals(validate.ErrImageTooLarge)
	Expect(validate.CheckDecodeBudget(mock.GIFHeader(40000, 40000))).Equals(validate.ErrImageTooLarge)

	// 12MP phone photo: 12.2MP * 19 bytes = 232MB
	Expect(validate.CheckDecodeBudget(jpegHeaderWithDimensions(t, 4032, 3024))).IsNil()
	// 16MP: 304MB
	Expect(validate.CheckDecodeBudget(jpegHeaderWithDimensions(t, 4000, 4000))).Equals(validate.ErrImageTooLarge)

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo2.jpg"))
	Expect(validate.CheckDecodeBudget(img)).IsNil()

	Expect(validate.CheckDecodeBudget([]byte("not an image"))).Equals(imagic.ErrNotSupported)
	favicon, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/favicon.ico"))
	Expect(validate.CheckDecodeBudget(favicon)).Equals(imagic.ErrNotSupported)
}

func TestSafeApply(t *testing.T) {
	RegisterT(t)

	resized, err := validate.SafeApply(context.Background(), mock.UniformPNG(2000, 1000), imagic.Resize(500))
	Expect(err).IsNil()
	width, height := imageDimensions(t, resized)
	Expect(width).Equals(500)
	Expect(height).Equals(250)

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	result, err := validate.SafeApply(context.Background(), mock.UniformPNG(12000, 12000), imagic.Resize(500))
	runtime.ReadMemStats(&after)
	Expect(err).Equals(validate.ErrImageTooLarge)
	Expect(result).IsNil()
	Expect(after.TotalAlloc-before.TotalAlloc < 10*1024*1024).IsTrue()

	_, err = validate.SafeApply(context.Background(), []byte("not an image"), imagic.Resize(500))
	Expect(err).Equals(imagic.ErrNotSupported)
}
