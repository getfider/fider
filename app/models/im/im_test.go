package im_test

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
	. "github.com/onsi/gomega"
)

var jonSnowToken, _ = jwt.Encode(&models.OAuthClaims{
	OAuthID:       "123",
	OAuthName:     "Jon Snow",
	OAuthEmail:    "jon.snow@got.com",
	OAuthProvider: "facebook",
})

func ExpectFailed(result *validate.Result) {
	Expect(result.Ok).To(BeFalse())
	Expect(len(result.Messages) > 0).To(BeTrue())
	Expect(result.Error).To(BeNil())
}

func ExpectSuccess(result *validate.Result) {
	Expect(result.Ok).To(BeTrue())
	Expect(len(result.Messages)).To(Equal(0))
	Expect(result.Error).To(BeNil())
}
