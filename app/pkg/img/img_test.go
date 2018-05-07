package img_test

import (
	"io/ioutil"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/img"
)

var testCases = []struct {
	fileName  string
	width     int
	height    int
	supported bool
}{
	{
		"/app/pkg/img/testdata/logo1.png",
		300,
		300,
		true,
	},
	{
		"/app/pkg/img/testdata/logo2.jpg",
		2624,
		2184,
		true,
	},
	{
		"/app/pkg/img/testdata/logo3.gif",
		1165,
		822,
		true,
	},
	{
		"/favicon.ico",
		0,
		0,
		false,
	},
}

func TestImageParse(t *testing.T) {
	RegisterT(t)

	for _, testCase := range testCases {
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
