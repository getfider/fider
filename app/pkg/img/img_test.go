package img_test

import (
	"image/color"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/img"
)

var parseTestCases = []struct {
	fileName  string
	width     int
	height    int
	supported bool
}{
	{"/app/pkg/img/testdata/logo1.png", 300, 300, true},
	{"/app/pkg/img/testdata/logo2.jpg", 2624, 2184, true},
	{"/app/pkg/img/testdata/logo3.gif", 1165, 822, true},
	{"/app/pkg/img/testdata/logo4.png", 150, 150, true},
	{"/app/pkg/img/testdata/logo5.png", 200, 200, true},
	{"/app/pkg/img/testdata/logo6.jpg", 400, 400, true},
	{"/app/pkg/img/testdata/logo7.gif", 400, 400, true},
	{"/app/pkg/img/testdata/favicon.ico", 0, 0, false},
}

func TestImageParse(t *testing.T) {
	RegisterT(t)

	for _, testCase := range parseTestCases {
		bytes, err := ioutil.ReadFile(env.Path(testCase.fileName))
		Expect(err).IsNil()

		file, err := img.Parse(bytes)
		if testCase.supported {
			Expect(err).IsNil()
			Expect(file.Width).Equals(testCase.width)
			Expect(file.Height).Equals(testCase.height)
			Expect(file.Size).Equals(len(bytes))
		} else {
			Expect(err).Equals(img.ErrNotSupported)
			Expect(file).IsNil()
		}
	}
}

var resizeTestCases = []struct {
	fileName        string
	resizedFileName string
	size            int
	padding         int
	skipCI          bool
}{
	{"/app/pkg/img/testdata/logo1.png", "/app/pkg/img/testdata/logo1-200x200.png", 200, 0, false},
	{"/app/pkg/img/testdata/logo2.jpg", "/app/pkg/img/testdata/logo2-200w.jpg", 200, 0, true},
	{"/app/pkg/img/testdata/logo2.jpg", "/app/pkg/img/testdata/logo2-200w-pad20.jpg", 200, 20, true},
	{"/app/pkg/img/testdata/logo3.gif", "/app/pkg/img/testdata/logo3-200w.gif", 200, 0, true},
	{"/app/pkg/img/testdata/logo4.png", "/app/pkg/img/testdata/logo4-100x100.png", 100, 0, false},
	{"/app/pkg/img/testdata/logo5.png", "/app/pkg/img/testdata/logo5-200x200.png", 200, 0, false},
	{"/app/pkg/img/testdata/logo6.jpg", "/app/pkg/img/testdata/logo6-200x200.jpg", 200, 0, false},
	{"/app/pkg/img/testdata/logo7.gif", "/app/pkg/img/testdata/logo7-200x200.gif", 200, 0, false},
	{"/app/pkg/img/testdata/logo7.gif", "/app/pkg/img/testdata/logo7-1000-1000.gif", 1000, 0, false},
	{"/app/pkg/img/testdata/logo8.png", "/app/pkg/img/testdata/logo8-200h.gif", 200, 0, false},
}

func TestImageResize(t *testing.T) {
	RegisterT(t)

	for _, testCase := range resizeTestCases {
		// we don't run it on CI because it's way too slow, we need to investigate why
		if testCase.skipCI && os.Getenv("CI") == "true" {
			continue
		}

		bytes, err := ioutil.ReadFile(env.Path(testCase.fileName))
		Expect(err).IsNil()

		resized, err := img.Apply(
			bytes,
			img.Resize(testCase.size),
			img.Padding(testCase.padding),
		)
		Expect(err).IsNil()

		expected, err := ioutil.ReadFile(env.Path(testCase.resizedFileName))
		Expect(err).IsNil()

		Expect(resized).Equals(expected)
	}
}

var bgColorTestCases = []struct {
	fileName           string
	whiteColorFileName string
	bgColor            color.Color
}{
	{"/app/pkg/img/testdata/logo1.png", "/app/pkg/img/testdata/logo1-white.png", color.White},
	{"/app/pkg/img/testdata/logo1.png", "/app/pkg/img/testdata/logo1-black.png", color.Black},
}

func TestImageChangeBackground(t *testing.T) {
	RegisterT(t)

	for _, testCase := range bgColorTestCases {
		bytes, err := ioutil.ReadFile(env.Path(testCase.fileName))
		Expect(err).IsNil()

		withColor, err := img.Apply(bytes, img.ChangeBackground(testCase.bgColor))
		Expect(err).IsNil()

		expected, err := ioutil.ReadFile(env.Path(testCase.whiteColorFileName))
		Expect(err).IsNil()

		Expect(withColor).Equals(expected)
	}
}
