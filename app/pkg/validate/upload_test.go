package validate_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
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
// Go always encodes baseline JPEGs (SOF0); progressive ones are made by patching it to SOF2.
func jpegHeaderWithDimensions(t *testing.T, width, height int, progressive bool) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, image.NewYCbCr(image.Rect(0, 0, 16, 16), image.YCbCrSubsampleRatio444), nil)
	Expect(err).IsNil()
	content := buf.Bytes()
	sof := bytes.Index(content, []byte{0xFF, 0xC0})
	Expect(sof > 0).IsTrue()
	// FF C0, length (2), precision (1), height (2), width (2)
	binary.BigEndian.PutUint16(content[sof+5:], uint16(height))
	binary.BigEndian.PutUint16(content[sof+7:], uint16(width))
	if progressive {
		content[sof+1] = 0xC2
	}
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
		// 6000x6000 8-bit gray with tRNS: decodes to NRGBA
		{"png gray+tRNS 6000x6000", mock.UniformGrayPNG(6000, 6000, 8, true)},
		// 4000x4000 progressive colour JPEG
		{"jpeg 4000x4000", jpegHeaderWithDimensions(t, 4000, 4000, true)},
		// 8000x8000 baseline colour JPEG
		{"jpeg 8000x8000", jpegHeaderWithDimensions(t, 8000, 8000, false)},
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

func header(width, height int, format string, model color.Model, progressive bool) *validate.ImageHeader {
	return &validate.ImageHeader{Width: width, Height: height, Format: format, ColorModel: model, Progressive: progressive}
}

func TestImageHeader_EstimateBytes(t *testing.T) {
	RegisterT(t)

	// Invalid or over the per-side cap
	Expect(header(0, 100, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(100, 0, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(-100, -100, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(validate.MaxImageSide+1, 1, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(1, validate.MaxImageSide+1, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(40_000_000, 1, "gif", color.Palette{}, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(math.MaxInt, math.MaxInt, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(-1))
	Expect(header(100, 100, "bmp", color.RGBAModel, false).EstimateBytes(1)).Equals(int64(-1))

	// With one operation: bytes per pixel include 4 bytes for the operation's intermediate image
	const mp = 1000 * 1000
	Expect(header(1000, 1000, "gif", color.Palette{}, false).EstimateBytes(1)).Equals(int64(5 * mp))
	Expect(header(1000, 1000, "png", color.Palette{}, false).EstimateBytes(1)).Equals(int64(6 * mp))
	// Gray PNGs might have a tRNS chunk, decoding to NRGBA/NRGBA64
	Expect(header(1000, 1000, "png", color.GrayModel, false).EstimateBytes(1)).Equals(int64(12 * mp))
	Expect(header(1000, 1000, "png", color.Gray16Model, false).EstimateBytes(1)).Equals(int64(20 * mp))
	Expect(header(1000, 1000, "png", color.RGBAModel, false).EstimateBytes(1)).Equals(int64(12 * mp))
	Expect(header(1000, 1000, "png", color.NRGBAModel, false).EstimateBytes(1)).Equals(int64(12 * mp))
	Expect(header(1000, 1000, "png", color.RGBA64Model, false).EstimateBytes(1)).Equals(int64(20 * mp))
	Expect(header(1000, 1000, "png", color.NRGBA64Model, false).EstimateBytes(1)).Equals(int64(20 * mp))

	// JPEG images are allocated in whole MCUs (up to 32x32), so 1000 is rounded up to 1024
	const jpegPixels = 1024 * 1024
	Expect(header(1000, 1000, "jpeg", color.GrayModel, false).EstimateBytes(1)).Equals(int64(5 * jpegPixels))
	// 3 component JPEGs might be converted to RGBA by a late Adobe APP14 segment
	Expect(header(1000, 1000, "jpeg", color.YCbCrModel, false).EstimateBytes(1)).Equals(int64(11 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.RGBAModel, false).EstimateBytes(1)).Equals(int64(11 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.CMYKModel, false).EstimateBytes(1)).Equals(int64(12 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.GrayModel, true).EstimateBytes(1)).Equals(int64(9 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.YCbCrModel, true).EstimateBytes(1)).Equals(int64(23 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.RGBAModel, true).EstimateBytes(1)).Equals(int64(23 * jpegPixels))
	Expect(header(1000, 1000, "jpeg", color.CMYKModel, true).EstimateBytes(1)).Equals(int64(28 * jpegPixels))
	Expect(header(1024, 1024, "jpeg", color.YCbCrModel, false).EstimateBytes(1)).Equals(int64(11 * jpegPixels))

	// progressive is ignored for other formats
	Expect(header(1000, 1000, "png", color.NRGBAModel, true).EstimateBytes(1)).Equals(int64(12 * mp))

	// Each operation adds 4 bytes per pixel
	Expect(header(1000, 1000, "png", color.GrayModel, false).EstimateBytes(0)).Equals(int64(8 * mp))
	Expect(header(1000, 1000, "png", color.GrayModel, false).EstimateBytes(3)).Equals(int64(20 * mp))

	// Worst case at the per-side cap doesn't overflow
	Expect(header(validate.MaxImageSide, validate.MaxImageSide, "jpeg", color.CMYKModel, true).EstimateBytes(3) > 0).IsTrue()
}

// decodeAllocations decodes the image and returns it with the bytes allocated to decode it
func decodeAllocations(content []byte) (image.Image, int64) {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	img, _, err := image.Decode(bytes.NewReader(content))
	runtime.ReadMemStats(&after)
	Expect(err).IsNil()
	return img, int64(after.TotalAlloc - before.TotalAlloc)
}

func TestImageHeader_EstimateCoversPNGWithTransparency(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		bitDepth    int
		transparent bool
		decodedType string
	}{
		{1, false, "*image.Gray"},
		{8, false, "*image.Gray"},
		// image.DecodeConfig doesn't see the tRNS chunk, but image/png decodes to NRGBA/NRGBA64
		{1, true, "*image.NRGBA"},
		{8, true, "*image.NRGBA"},
		{16, false, "*image.Gray16"},
		{16, true, "*image.NRGBA64"},
	}

	for _, testCase := range testCases {
		content := mock.UniformGrayPNG(1000, 1000, testCase.bitDepth, testCase.transparent)
		header, err := validate.ReadImageHeader(content)
		Expect(err).IsNil()

		img, allocated := decodeAllocations(content)
		Expect(fmt.Sprintf("%T", img)).Equals(testCase.decodedType)
		Expect(header.EstimateBytes(0) >= allocated).IsTrue()
	}

	// A 36MP gray+tRNS PNG decodes to 144MB of NRGBA (288MB if interlaced), but is only a few KB
	content := mock.UniformGrayPNG(6000, 6000, 8, true)
	Expect(len(content) < 100*1024).IsTrue()
	_, err := validate.CheckDecodeBudget(content, 1)
	Expect(err).Equals(validate.ErrImageTooLarge)
}

// withLateAdobeAPP14 inserts an Adobe APP14 segment with transform=0 (RGB) just before EOI.
// image.DecodeConfig stops at the first SOS so never sees it, but a full decode
// converts the image to *image.RGBA.
func withLateAdobeAPP14(content []byte) []byte {
	Expect(bytes.HasSuffix(content, []byte{0xFF, 0xD9})).IsTrue()
	app14 := []byte{0xFF, 0xEE, 0x00, 0x0E, 'A', 'd', 'o', 'b', 'e', 0, 100, 0, 0, 0, 0, 0}
	result := append([]byte{}, content[:len(content)-2]...)
	result = append(result, app14...)
	return append(result, 0xFF, 0xD9)
}

func TestImageHeader_EstimateCoversJPEGWithLateAdobeAPP14(t *testing.T) {
	RegisterT(t)

	var buf bytes.Buffer
	Expect(jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 1000, 1000)), nil)).IsNil()
	content := withLateAdobeAPP14(buf.Bytes())

	header, err := validate.ReadImageHeader(content)
	Expect(err).IsNil()
	Expect(header.ColorModel).Equals(color.YCbCrModel)
	Expect(header.Progressive).IsFalse()

	img, allocated := decodeAllocations(content)
	Expect(fmt.Sprintf("%T", img)).Equals("*image.RGBA")
	Expect(header.EstimateBytes(0) >= allocated).IsTrue()

	// Without the APP14 segment, still decoded as YCbCr
	img, allocated = decodeAllocations(buf.Bytes())
	Expect(fmt.Sprintf("%T", img)).Equals("*image.YCbCr")
	Expect(header.EstimateBytes(0) >= allocated).IsTrue()
}

func TestCheckDecodeBudget(t *testing.T) {
	RegisterT(t)

	check := func(content []byte, operations int) error {
		_, err := validate.CheckDecodeBudget(content, operations)
		return err
	}

	h, err := validate.CheckDecodeBudget(mock.UniformPNG(2000, 1000), 1)
	Expect(err).IsNil()
	Expect(h.Width).Equals(2000)
	Expect(h.Height).Equals(1000)
	Expect(h.Format).Equals("png")

	// 16MP * 12 bytes = 192MB
	Expect(check(mock.UniformPNG(4000, 4000), 1)).IsNil()
	// ... but not with 3 operations (Favicon with background): 16MP * 20 bytes = 320MB
	Expect(check(mock.UniformPNG(4000, 4000), 3)).Equals(validate.ErrImageTooLarge)
	// 25MP * 12 bytes = 300MB
	Expect(check(mock.UniformPNG(5000, 5000), 1)).Equals(validate.ErrImageTooLarge)
	Expect(check(mock.UniformPNG(12000, 12000), 1)).Equals(validate.ErrImageTooLarge)
	Expect(check(mock.UniformPNG(100000, 1), 1)).Equals(validate.ErrImageTooLarge)
	Expect(check(mock.GIFHeader(40000, 40000), 1)).Equals(validate.ErrImageTooLarge)

	// The header is returned when over budget
	h, err = validate.CheckDecodeBudget(mock.UniformPNG(12000, 12000), 1)
	Expect(err).Equals(validate.ErrImageTooLarge)
	Expect(h.Width).Equals(12000)

	// 24MP baseline colour JPEG: 6016x4000 * 11 bytes = 265MB
	Expect(check(jpegHeaderWithDimensions(t, 6000, 4000, false), 1)).IsNil()
	// Same dimensions, but progressive: 24MP * 23 bytes = 553MB
	Expect(check(jpegHeaderWithDimensions(t, 6000, 4000, true), 1)).Equals(validate.ErrImageTooLarge)
	// 9MP progressive: 3008x3008 * 23 bytes = 208MB
	Expect(check(jpegHeaderWithDimensions(t, 3000, 3000, true), 1)).IsNil()
	// 12MP progressive phone photo: 4032x3040 * 23 bytes = 282MB
	Expect(check(jpegHeaderWithDimensions(t, 4032, 3024, true), 1)).Equals(validate.ErrImageTooLarge)
	// 40MP baseline: 8000x5024 * 11 bytes = 442MB
	Expect(check(jpegHeaderWithDimensions(t, 8000, 5000, false), 1)).Equals(validate.ErrImageTooLarge)

	img, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/logo2.jpg"))
	Expect(check(img, 1)).IsNil()

	h, err = validate.CheckDecodeBudget([]byte("not an image"), 1)
	Expect(err).Equals(imagic.ErrNotSupported)
	Expect(h).IsNil()
	favicon, _ := os.ReadFile(env.Path("/app/pkg/web/testdata/favicon.ico"))
	Expect(check(favicon, 1)).Equals(imagic.ErrNotSupported)
}

func TestReadImageHeader_SameFormatsAsImagic(t *testing.T) {
	RegisterT(t)

	files, err := filepath.Glob(env.Path("/app/pkg/web/testdata/*"))
	Expect(err).IsNil()
	Expect(len(files) > 5).IsTrue()
	files = append(files, env.Path("/README.md"), env.Path("/favicon.png"))

	for _, file := range files {
		content, _ := os.ReadFile(file)
		_, parseErr := imagic.Parse(content)
		_, headerErr := validate.ReadImageHeader(content)
		Expect(headerErr == nil).Equals(parseErr == nil)
	}
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

	// A valid header but corrupt pixel data fails to decode
	_, err = validate.SafeApply(context.Background(), mock.GIFHeader(100, 100), imagic.Resize(50))
	Expect(err).IsNotNil()
}

func TestSafeResize(t *testing.T) {
	RegisterT(t)

	content := mock.UniformPNG(2000, 1000)
	resized, err := validate.SafeResize(context.Background(), content, 500)
	Expect(err).IsNil()
	width, height := imageDimensions(t, resized)
	Expect(width).Equals(500)
	Expect(height).Equals(250)

	// Already within size: the original content is returned, without decoding it
	for _, size := range []int{2000, 3000} {
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		result, err := validate.SafeResize(context.Background(), content, size)
		runtime.ReadMemStats(&after)
		Expect(err).IsNil()
		Expect(&result[0] == &content[0]).IsTrue()
		Expect(after.TotalAlloc-before.TotalAlloc < 1024*1024).IsTrue()
	}

	_, err = validate.SafeResize(context.Background(), mock.UniformPNG(12000, 12000), 500)
	Expect(err).Equals(validate.ErrImageTooLarge)
}
