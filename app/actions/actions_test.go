package actions_test

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
	Users:   &inmemory.UserStorage{},
}

func ExpectFailed(result *validate.Result, fields ...string) {
	Expect(result.Ok).To(BeFalse())
	Expect(result.Error).To(BeNil())
	for _, field := range fields {
		if field == "" {
			Expect(len(result.Messages) > 0).To(BeTrue())
		} else {
			Expect(len(result.Failures[field]) > 0).To(BeTrue(), "Failures should contain field: "+field)
		}
	}

	for field, _ := range result.Failures {
		Expect(fields).To(ContainElement(field), "Failures should not contain field: "+field)
	}
}

func ExpectSuccess(result *validate.Result) {
	Expect(result.Ok).To(BeTrue())
	Expect(len(result.Messages)).To(Equal(0))
	Expect(len(result.Failures)).To(Equal(0))
	Expect(result.Error).To(BeNil())
}
