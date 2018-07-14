package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/rand"
)

func TestCreateEditOAuthConfig_InvalidInput(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		expected []string
		input    *models.CreateEditOAuthConfig
	}{
		{
			expected: []string{"displayName", "tokenURL", "clientId", "clientSecret", "scope", "authorizeURL", "tokenURL", "profileURL", "jsonUserIdPath"},
			input:    &models.CreateEditOAuthConfig{},
		},
		{
			expected: []string{"displayName", "tokenURL", "clientId", "clientSecret", "scope", "authorizeURL", "tokenURL", "profileURL", "jsonUserIdPath", "jsonUserNamePath", "jsonUserEmailPath"},
			input: &models.CreateEditOAuthConfig{
				DisplayName:       rand.String(51),
				ClientID:          rand.String(101),
				ClientSecret:      rand.String(501),
				AuthorizeURL:      rand.String(301),
				TokenURL:          rand.String(301),
				Scope:             rand.String(101),
				ProfileURL:        rand.String(301),
				JSONUserIDPath:    rand.String(101),
				JSONUserNamePath:  rand.String(101),
				JSONUserEmailPath: rand.String(101),
			},
		},
	}

	for _, testCase := range testCases {
		action := &actions.CreateEditOAuthConfig{
			Model: testCase.input,
		}
		result := action.Validate(nil, services)
		ExpectFailed(result, testCase.expected...)
	}
}

func TestCreateEditOAuthConfig_ValidInput(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditOAuthConfig{
		Model: &models.CreateEditOAuthConfig{
			DisplayName:       "My Provider",
			ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
			ClientSecret:      "jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij",
			AuthorizeURL:      "http://provider/oauth/authorize",
			TokenURL:          "http://provider/oauth/token",
			Scope:             "profile email",
			ProfileURL:        "http://provider/profile/me",
			JSONUserIDPath:    "user.id",
			JSONUserNamePath:  "user.name",
			JSONUserEmailPath: "user.email",
		},
	}
	result := action.Validate(nil, services)
	ExpectSuccess(result)
}
