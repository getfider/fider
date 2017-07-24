package im_test

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

var jonSnowToken, _ = jwt.Encode(&models.OAuthClaims{
	OAuthID:       "123",
	OAuthName:     "Jon Snow",
	OAuthEmail:    "jon.snow@got.com",
	OAuthProvider: "facebook",
})

var services = &app.Services{
	Tenants: &inmemory.TenantStorage{},
}

func ExpectFailed(result *validate.Result, field string) {
	Expect(result.Ok).To(BeFalse())
	Expect(result.Error).To(BeNil())
	if field == "" {
		Expect(len(result.Messages) > 0).To(BeTrue())
	} else {
		Expect(len(result.Failures[field]) > 0).To(BeTrue())
	}
}

func ExpectSuccess(result *validate.Result) {
	Expect(result.Ok).To(BeTrue())
	Expect(len(result.Messages)).To(Equal(0))
	Expect(len(result.Failures)).To(Equal(0))
	Expect(result.Error).To(BeNil())
}
