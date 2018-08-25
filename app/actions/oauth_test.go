package actions_test

import (
	"testing"

	"github.com/getfider/fider/app"
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
			expected: []string{"displayName", "status", "tokenURL", "clientID", "clientSecret", "scope", "authorizeURL", "tokenURL", "jsonUserIDPath"},
			input:    &models.CreateEditOAuthConfig{},
		},
		{
			expected: []string{"displayName", "status", "tokenURL", "clientID", "clientSecret", "scope", "authorizeURL", "tokenURL", "profileURL", "jsonUserIDPath", "jsonUserNamePath", "jsonUserEmailPath"},
			input: &models.CreateEditOAuthConfig{
				DisplayName:       rand.String(51),
				ClientID:          rand.String(101),
				Status:            0,
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

func TestCreateEditOAuthConfig_AddNew_ValidInput(t *testing.T) {
	RegisterT(t)

	input := &models.CreateEditOAuthConfig{
		DisplayName:       "My Provider",
		Status:            models.OAuthConfigEnabled,
		ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
		ClientSecret:      "jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij",
		AuthorizeURL:      "http://provider/oauth/authorize",
		TokenURL:          "http://provider/oauth/token",
		Scope:             "profile email",
		ProfileURL:        "http://provider/profile/me",
		JSONUserIDPath:    "user.id",
		JSONUserNamePath:  "user.name",
		JSONUserEmailPath: "user.email",
	}
	action := &actions.CreateEditOAuthConfig{
		Model: input,
	}
	result := action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(input.ID).Equals(0)
	Expect(input.Provider).HasLen(11)
	Expect(string(input.Provider[0])).Equals("_")
}

func TestCreateEditOAuthConfig_EditExisting_NewSecret(t *testing.T) {
	RegisterT(t)

	services.Tenants.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		ID:           4,
		Provider:     "_NAME",
		DisplayName:  "My Provider",
		ClientSecret: "MY_OLD_SECRET",
	})

	input := &models.CreateEditOAuthConfig{
		Provider:          "_NAME",
		DisplayName:       "My Provider",
		Status:            models.OAuthConfigDisabled,
		ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
		ClientSecret:      "jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij",
		AuthorizeURL:      "http://provider/oauth/authorize",
		TokenURL:          "http://provider/oauth/token",
		Scope:             "profile email",
		ProfileURL:        "http://provider/profile/me",
		JSONUserIDPath:    "user.id",
		JSONUserNamePath:  "user.name",
		JSONUserEmailPath: "user.email",
	}
	action := &actions.CreateEditOAuthConfig{
		Model: input,
	}
	result := action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(input.ID).Equals(4)
	Expect(input.ClientSecret).Equals("jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
}

func TestCreateEditOAuthConfig_EditExisting_OmitSecret(t *testing.T) {
	RegisterT(t)

	services.Tenants.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		ID:           5,
		Provider:     "_NAME2",
		DisplayName:  "My Provider",
		ClientSecret: "MY_OLD_SECRET",
	})

	input := &models.CreateEditOAuthConfig{
		Provider:          "_NAME2",
		DisplayName:       "My Provider",
		Status:            models.OAuthConfigDisabled,
		ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
		AuthorizeURL:      "http://provider/oauth/authorize",
		TokenURL:          "http://provider/oauth/token",
		Scope:             "profile email",
		ProfileURL:        "http://provider/profile/me",
		JSONUserIDPath:    "user.id",
		JSONUserNamePath:  "user.name",
		JSONUserEmailPath: "user.email",
	}
	action := &actions.CreateEditOAuthConfig{
		Model: input,
	}
	result := action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(input.ID).Equals(5)
	Expect(input.ClientSecret).Equals("MY_OLD_SECRET")
}

func TestCreateEditOAuthConfig_EditNonExisting(t *testing.T) {
	RegisterT(t)

	input := &models.CreateEditOAuthConfig{
		Provider:          "_MY_NEW_PROVIDER",
		DisplayName:       "My Provider",
		Status:            models.OAuthConfigDisabled,
		ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
		AuthorizeURL:      "http://provider/oauth/authorize",
		TokenURL:          "http://provider/oauth/token",
		Scope:             "profile email",
		ProfileURL:        "http://provider/profile/me",
		JSONUserIDPath:    "user.id",
		JSONUserNamePath:  "user.name",
		JSONUserEmailPath: "user.email",
	}
	action := &actions.CreateEditOAuthConfig{
		Model: input,
	}
	result := action.Validate(nil, services)
	Expect(result.Err).Equals(app.ErrNotFound)
	Expect(result.Ok).IsFalse()
}
