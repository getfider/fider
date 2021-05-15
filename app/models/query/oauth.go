package query

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entities"
)

type GetCustomOAuthConfigByProvider struct {
	Provider string

	Result *entities.OAuthConfig
}

type ListCustomOAuthConfig struct {
	Result []*entities.OAuthConfig
}

type GetOAuthAuthorizationURL struct {
	Provider   string
	Redirect   string
	Identifier string

	Result string
}

type GetOAuthProfile struct {
	Provider string
	Code     string

	Result *dto.OAuthUserProfile
}

type GetOAuthRawProfile struct {
	Provider string
	Code     string

	Result string
}

type ListActiveOAuthProviders struct {
	Result []*dto.OAuthProviderOption
}

type ListAllOAuthProviders struct {
	Result []*dto.OAuthProviderOption
}
