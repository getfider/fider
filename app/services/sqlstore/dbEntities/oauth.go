package dbEntities

import (
	"github.com/getfider/fider/app/models/entity"
)

type oauthConfig struct {
	ID                int    `db:"id"`
	Provider          string `db:"provider"`
	DisplayName       string `db:"display_name"`
	LogoBlobKey       string `db:"logo_bkey"`
	Status            int    `db:"status"`
	IsTrusted         bool   `db:"is_trusted"`
	ClientID          string `db:"client_id"`
	ClientSecret      string `db:"client_secret"`
	AuthorizeURL      string `db:"authorize_url"`
	TokenURL          string `db:"token_url"`
	Scope             string `db:"scope"`
	ProfileURL        string `db:"profile_url"`
	JSONUserIDPath    string `db:"json_user_id_path"`
	JSONUserNamePath  string `db:"json_user_name_path"`
	JSONUserEmailPath string `db:"json_user_email_path"`
}

func (m *oauthConfig) toModel() *entity.OAuthConfig {
	return &entity.OAuthConfig{
		ID:                m.ID,
		Provider:          m.Provider,
		DisplayName:       m.DisplayName,
		Status:            m.Status,
		IsTrusted:         m.IsTrusted,
		LogoBlobKey:       m.LogoBlobKey,
		ClientID:          m.ClientID,
		ClientSecret:      m.ClientSecret,
		AuthorizeURL:      m.AuthorizeURL,
		TokenURL:          m.TokenURL,
		ProfileURL:        m.ProfileURL,
		Scope:             m.Scope,
		JSONUserIDPath:    m.JSONUserIDPath,
		JSONUserNamePath:  m.JSONUserNamePath,
		JSONUserEmailPath: m.JSONUserEmailPath,
	}
}
