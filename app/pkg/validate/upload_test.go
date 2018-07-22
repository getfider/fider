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
		valid    bool
	}{
		{"/app/pkg/img/testdata/logo1.png", true},
		{"/app/pkg/img/testdata/logo2.jpg", false},
		{"/app/pkg/img/testdata/logo3.gif", false},
		{"/app/pkg/img/testdata/logo4.png", false},
		{"/app/pkg/img/testdata/logo5.png", true},
		{"/README.md", false},
		{"/favicon.ico", false},
	}

	for _, testCase := range testCases {
		img, _ := ioutil.ReadFile(env.Path(testCase.fileName))

		upload := &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content: img,
			},
		}
		result := validate.ImageUpload(upload, 200, 200, 100)
		Expect(result.Ok).Equals(testCase.valid)
	}
}

func TestValidateImageUpload_Nil(t *testing.T) {
	RegisterT(t)

	result := validate.ImageUpload(nil, 200, 200, 50)
	Expect(result.Ok).IsTrue()

	result = validate.ImageUpload(&models.ImageUpload{}, 200, 200, 50)
	Expect(result.Ok).IsTrue()
}
