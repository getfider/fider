package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
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
		result := action.Validate(context.Background(), nil)
		ExpectFailed(result, testCase.expected...)
	}
}

func TestCreateEditOAuthConfig_Initialize(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditOAuthConfig{}
	action.Initialize()
	Expect(action.Model.Logo.BlobKey).Equals("")
}

func TestCreateEditOAuthConfig_AddNew_ValidInput(t *testing.T) {
	RegisterT(t)

	input := &models.CreateEditOAuthConfig{
		DisplayName:       "My Provider",
		Status:            enum.OAuthConfigEnabled,
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
	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(input.ID).Equals(0)
	Expect(input.Provider).HasLen(11)
	Expect(string(input.Provider[0])).Equals("_")
}

func TestCreateEditOAuthConfig_EditExisting_NewSecret(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_NAME" {
			q.Result = &models.OAuthConfig{
				ID:          4,
				Provider:    q.Provider,
				LogoBlobKey: "hello-world.png",
			}
			return nil
		}
		return app.ErrNotFound
	})

	action := &actions.CreateEditOAuthConfig{}
	action.Initialize()
	action.Model.Provider = "_NAME"
	action.Model.DisplayName = "My Provider"
	action.Model.Status = enum.OAuthConfigDisabled
	action.Model.ClientID = "823187ahjjfdha8fds7yfdashfjkdsa"
	action.Model.ClientSecret = "jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij"
	action.Model.AuthorizeURL = "http://provider/oauth/authorize"
	action.Model.TokenURL = "http://provider/oauth/token"
	action.Model.Scope = "profile email"
	action.Model.ProfileURL = "http://provider/profile/me"
	action.Model.JSONUserIDPath = "user.id"
	action.Model.JSONUserNamePath = "user.name"
	action.Model.JSONUserEmailPath = "user.email"

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.Model.ID).Equals(4)
	Expect(action.Model.Logo.BlobKey).Equals("hello-world.png")
	Expect(action.Model.ClientSecret).Equals("jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
}

func TestCreateEditOAuthConfig_EditExisting_OmitSecret(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		if q.Provider == "_NAME2" {
			q.Result = &models.OAuthConfig{
				ID:           5,
				Provider:     q.Provider,
				DisplayName:  "My Provider",
				ClientSecret: "MY_OLD_SECRET",
			}
			return nil
		}
		return app.ErrNotFound
	})

	action := &actions.CreateEditOAuthConfig{}
	action.Initialize()
	action.Model.Provider = "_NAME2"
	action.Model.DisplayName = "My Provider"
	action.Model.Status = enum.OAuthConfigDisabled
	action.Model.ClientID = "823187ahjjfdha8fds7yfdashfjkdsa"
	action.Model.AuthorizeURL = "http://provider/oauth/authorize"
	action.Model.TokenURL = "http://provider/oauth/token"
	action.Model.Scope = "profile email"
	action.Model.ProfileURL = "http://provider/profile/me"
	action.Model.JSONUserIDPath = "user.id"
	action.Model.JSONUserNamePath = "user.name"
	action.Model.JSONUserEmailPath = "user.email"

	result := action.Validate(context.Background(), nil)
	ExpectSuccess(result)
	Expect(action.Model.ID).Equals(5)
	Expect(action.Model.ClientSecret).Equals("MY_OLD_SECRET")
}

func TestCreateEditOAuthConfig_EditNonExisting(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
		return app.ErrNotFound
	})

	action := &actions.CreateEditOAuthConfig{}
	action.Initialize()
	action.Model.Provider = "_MY_NEW_PROVIDER"
	action.Model.DisplayName = "My Provider"
	action.Model.Status = enum.OAuthConfigDisabled
	action.Model.ClientID = "823187ahjjfdha8fds7yfdashfjkdsa"
	action.Model.AuthorizeURL = "http://provider/oauth/authorize"
	action.Model.TokenURL = "http://provider/oauth/token"
	action.Model.Scope = "profile email"
	action.Model.ProfileURL = "http://provider/profile/me"
	action.Model.JSONUserIDPath = "user.id"
	action.Model.JSONUserNamePath = "user.name"
	action.Model.JSONUserEmailPath = "user.email"
	result := action.Validate(context.Background(), nil)
	Expect(result.Err).Equals(app.ErrNotFound)
	Expect(result.Ok).IsFalse()
}
