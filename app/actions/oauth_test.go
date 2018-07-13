package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
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
				DisplayName:       mock.RandomString(51),
				ClientID:          mock.RandomString(101),
				ClientSecret:      mock.RandomString(501),
				AuthorizeURL:      mock.RandomString(301),
				TokenURL:          mock.RandomString(301),
				Scope:             mock.RandomString(101),
				ProfileURL:        mock.RandomString(301),
				JSONUserIDPath:    mock.RandomString(101),
				JSONUserNamePath:  mock.RandomString(101),
				JSONUserEmailPath: mock.RandomString(101),
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
