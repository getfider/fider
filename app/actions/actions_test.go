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
	Posts:         inmemory.NewPostStorage(),
	Tags:          inmemory.NewTagStorage(),
	Notifications: inmemory.NewNotificationStorage(),
}

func ExpectFailed(result *validate.Result, expectedFields ...string) {
	Expect(result.Ok).IsFalse()
	Expect(result.Err).IsNil()

	errFields := make([]string, 0)
	for _, err := range result.Errors {
		if err.Field != "" {
			errFields = append(errFields, err.Field)
		}
	}

	for _, field := range expectedFields {
		if field == "" {
			Expect(len(result.Errors) > 0).IsTrue()
		} else {
			if !contains(errFields, field) {
				Fail("Failure should contain field: %s", field)
			}
		}
	}

	for _, field := range errFields {
		if !contains(expectedFields, field) {
			Fail("Failure should not contain field: %s", field)
		}
	}
}

func ExpectSuccess(result *validate.Result) {
	Expect(result.Ok).IsTrue()
	Expect(result.Errors).HasLen(0)
	Expect(result.Err).IsNil()
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
