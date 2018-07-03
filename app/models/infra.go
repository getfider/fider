package models

import (
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
	ID             int    `db:"id"`
	Provider       string `db:"provider"`
	DisplayName    string `db:"display_name"`
	Status         int    `db:"status"`
	ClientID       string `db:"client_id"`
	ClientSecret   string `db:"client_secret"`
	AuthorizeURL   string `db:"authorize_url"`
	TokenURL       string `db:"token_url"`
	ProfileURL     string `db:"profile_url"`
	Scope          string `db:"scope"`
	JSONUserIDPath string `db:"json_user_id_path"`
	JSONNamePath   string `db:"json_name_path"`
	JSONEmailPath  string `db:"json_email_path"`
}
