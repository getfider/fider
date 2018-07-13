package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// SystemSettings is the system-wide settings
type SystemSettings struct {
	Mode            string
	BuildTime       string
	Version         string
	Environment     string
	GoogleAnalytics string
	Compiler        string
	Domain          string
	HasLegal        bool
}

// Notification is the system generated notification entity
type Notification struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Link      string    `json:"link" db:"link"`
	Read      bool      `json:"read" db:"read"`
	CreatedOn time.Time `json:"createdOn" db:"created_on"`
}

// OAuthConfig is the configuration of a custom OAuth provider
type OAuthConfig struct {
	ID             int    `json:"id" db:"id"`
	Provider       string `json:"provider" db:"provider"`
	DisplayName    string `json:"displayName" db:"display_name"`
	Status         int    `json:"status" db:"status"`
	ClientID       string `json:"clientId" db:"client_id"`
	ClientSecret   string `json:"-" db:"client_secret"`
	AuthorizeURL   string `json:"authorizeUrl" db:"authorize_url"`
	TokenURL       string `json:"tokenUrl" db:"token_url"`
	ProfileURL     string `json:"profileUrl" db:"profile_url"`
	Scope          string `json:"scope" db:"scope"`
	JSONUserIDPath string `json:"jsonUserIdPath" db:"json_user_id_path"`
	JSONNamePath   string `json:"jsonUserNamePath" db:"json_user_name_path"`
	JSONEmailPath  string `json:"jsonUserEmailPath" db:"json_user_email_path"`
}

// MarshalJSON converts model into a JSON string
func (o *OAuthConfig) MarshalJSON() ([]byte, error) {
	secret := o.ClientSecret
	if len(secret) >= 10 {
		secret = fmt.Sprintf("%s...%s", secret[0:3], secret[len(secret)-3:])
	} else {
		secret = "..."
	}

	return json.Marshal(&struct {
		ID             int    `json:"id"`
		Provider       string `json:"provider"`
		DisplayName    string `json:"displayName"`
		Status         int    `json:"status"`
		ClientID       string `json:"clientId"`
		ClientSecret   string `json:"clientSecret"`
		AuthorizeURL   string `json:"authorizeUrl"`
		TokenURL       string `json:"tokenUrl"`
		ProfileURL     string `json:"profileUrl"`
		Scope          string `json:"scope"`
		JSONUserIDPath string `json:"jsonUserIdPath"`
		JSONNamePath   string `json:"jsonUserNamePath"`
		JSONEmailPath  string `json:"jsonUserEmailPath"`
	}{
		ID:             o.ID,
		Provider:       o.Provider,
		DisplayName:    o.DisplayName,
		Status:         o.Status,
		ClientID:       o.ClientID,
		ClientSecret:   secret,
		AuthorizeURL:   o.AuthorizeURL,
		TokenURL:       o.TokenURL,
		ProfileURL:     o.ProfileURL,
		Scope:          o.Scope,
		JSONUserIDPath: o.JSONUserIDPath,
		JSONNamePath:   o.JSONNamePath,
		JSONEmailPath:  o.JSONEmailPath,
	})
}
