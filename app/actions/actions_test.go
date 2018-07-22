package actions_test

import (
	"github.com/getfider/fider/app"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
)

var jonSnowToken, _ = jwt.Encode(jwt.OAuthClaims{
	OAuthID:       "123",
	OAuthName:     "Jon Snow",
	OAuthEmail:    "jon.snow@got.com",
	OAuthProvider: "facebook",
})

var services = &app.Services{
	Tenants:       inmemory.NewTenantStorage(),
	Users:         &inmemory.UserStorage{},
	Ideas:         inmemory.NewIdeaStorage(),
	Tags:          inmemory.NewTagStorage(),
	Notifications: inmemory.NewNotificationStorage(),
}

func ExpectFailed(result *validate.Result, fields ...string) {
	Expect(result.Ok).IsFalse()
	Expect(result.Error).IsNil()
	for _, field := range fields {
		if field == "" {
			Expect(len(result.Messages) > 0).IsTrue()
		} else {
			if len(result.Failures[field]) == 0 {
				Fail("Failures should contain field: %s", field)
			}
		}
	}

	for field, _ := range result.Failures {
		if !contains(fields, field) {
			Fail("Failures should not contain field: %s", field)
		}
	}
}

func ExpectSuccess(result *validate.Result) {
	Expect(result.Ok).IsTrue()
	Expect(result.Messages).HasLen(0)
	Expect(result.Failures).HasLen(0)
	Expect(result.Error).IsNil()
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
