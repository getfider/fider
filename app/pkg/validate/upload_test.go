package validate_test

import (
	"io/ioutil"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestValidateImageUpload(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		fileName string
		count    int
	}{
		{"/app/pkg/img/testdata/logo1.png", 0},
		{"/app/pkg/img/testdata/logo2.jpg", 2},
		{"/app/pkg/img/testdata/logo3.gif", 1},
		{"/app/pkg/img/testdata/logo4.png", 1},
		{"/app/pkg/img/testdata/logo5.png", 0},
		{"/README.md", 1},
		{"/favicon.ico", 1},
	}

	for _, testCase := range testCases {
		img, _ := ioutil.ReadFile(env.Path(testCase.fileName))

		upload := &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content: img,
			},
		}
		messages, err := validate.ImageUpload(upload, 200, 200, 100)
		Expect(messages).HasLen(testCase.count)
		Expect(err).IsNil()
	}
}

func TestValidateImageUpload_Nil(t *testing.T) {
	RegisterT(t)

	messages, err := validate.ImageUpload(nil, 200, 200, 50)
	Expect(messages).HasLen(0)
	Expect(err).IsNil()

	messages, err = validate.ImageUpload(&models.ImageUpload{}, 200, 200, 50)
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}
