package img_test

import (
	"image/color"
	"io/ioutil"
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
}{
	{"/app/pkg/img/testdata/logo1.png", "/app/pkg/img/testdata/logo1-200x200.png", 200, 0},
	{"/app/pkg/img/testdata/logo2.jpg", "/app/pkg/img/testdata/logo2-200x200.jpg", 200, 0},
	{"/app/pkg/img/testdata/logo3.gif", "/app/pkg/img/testdata/logo3-200x200.gif", 200, 0},
	{"/app/pkg/img/testdata/logo4.png", "/app/pkg/img/testdata/logo4-100x100.png", 100, 0},
	{"/app/pkg/img/testdata/logo5.png", "/app/pkg/img/testdata/logo5-200x200.png", 200, 0},
	{"/app/pkg/img/testdata/logo6.jpg", "/app/pkg/img/testdata/logo6-200x200.jpg", 200, 0},
	{"/app/pkg/img/testdata/logo7.gif", "/app/pkg/img/testdata/logo7-200x200.gif", 200, 0},
	{"/app/pkg/img/testdata/logo7.gif", "/app/pkg/img/testdata/logo7-1000-1000.gif", 1000, 0},
}

func TestImageResize(t *testing.T) {
	RegisterT(t)

	for _, testCase := range resizeTestCases {
		bytes, err := ioutil.ReadFile(env.Path(testCase.fileName))
		Expect(err).IsNil()

		resized, err := img.Resize(bytes, testCase.size, testCase.padding)
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
	{"/app/pkg/img/testdata/logo2.jpg", "/app/pkg/img/testdata/logo2.jpg", color.White},
	{"/app/pkg/img/testdata/logo7.gif", "/app/pkg/img/testdata/logo7.gif", color.White},
	{"/app/pkg/img/testdata/logo4.png", "/app/pkg/img/testdata/logo4-white.png", color.White},
	{"/app/pkg/img/testdata/logo4.png", "/app/pkg/img/testdata/logo4-black.png", color.Black},
	{"/app/pkg/img/testdata/logo5.png", "/app/pkg/img/testdata/logo5-white.png", color.White},
	{"/app/pkg/img/testdata/logo5.png", "/app/pkg/img/testdata/logo5-black.png", color.Black},
}

func TestImageChangeBackground(t *testing.T) {
	RegisterT(t)

	for _, testCase := range bgColorTestCases {
		bytes, err := ioutil.ReadFile(env.Path(testCase.fileName))
		Expect(err).IsNil()

		withColor, err := img.ChangeBackground(bytes, testCase.bgColor)
		Expect(err).IsNil()

		expected, err := ioutil.ReadFile(env.Path(testCase.whiteColorFileName))
		Expect(err).IsNil()

		Expect(withColor).Equals(expected)
	}
}
