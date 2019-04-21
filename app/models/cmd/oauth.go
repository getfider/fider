package cmd

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/dto"
)

type SaveCustomOAuthConfig struct {
	Config *models.CreateEditOAuthConfig
}

type ParseOAuthRawProfile struct {
	Provider string
	Body     string

	Result *dto.OAuthUserProfile
}
