package cmd

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
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
	IsTrusted         bool
	JSONUserIDPath    string
	JSONUserNamePath  string
	JSONUserEmailPath string
}

type ParseOAuthRawProfile struct {
	Provider string
	Body     string

	Result *dto.OAuthUserProfile
}
