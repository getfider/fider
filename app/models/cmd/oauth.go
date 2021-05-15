package cmd

import (
	"github.com/getfider/fider/app/models/dto"
)

type SaveCustomOAuthConfig struct {
	ID                int
	Logo              *dto.ImageUpload
	Provider          string
	Status            int
	DisplayName       string
	ClientID          string
	ClientSecret      string
	AuthorizeURL      string
	TokenURL          string
	Scope             string
	ProfileURL        string
	JSONUserIDPath    string
	JSONUserNamePath  string
	JSONUserEmailPath string
}

type ParseOAuthRawProfile struct {
	Provider string
	Body     string

	Result *dto.OAuthUserProfile
}
