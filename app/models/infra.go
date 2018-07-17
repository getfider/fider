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

// CreateEditOAuthConfig is used to create/edit an OAuth Configuration
type CreateEditOAuthConfig struct {
	ID                int
	Logo              *ImageUpload `json:"logo"`
	Provider          string       `json:"provider"`
	Status            int          `json:"status"`
	DisplayName       string       `json:"displayName"`
	ClientID          string       `json:"clientId"`
	ClientSecret      string       `json:"clientSecret"`
	AuthorizeURL      string       `json:"authorizeUrl" format:"lower"`
	TokenURL          string       `json:"tokenUrl" format:"lower"`
	Scope             string       `json:"scope"`
	ProfileURL        string       `json:"profileUrl" format:"lower"`
	JSONUserIDPath    string       `json:"jsonUserIdPath"`
	JSONUserNamePath  string       `json:"jsonUserNamePath"`
	JSONUserEmailPath string       `json:"jsonUserEmailPath"`
}

// OAuthConfig is the configuration of a custom OAuth provider
type OAuthConfig struct {
	ID                int    `json:"id"`
	Provider          string `json:"provider"`
	DisplayName       string `json:"displayName"`
	LogoURL           string `json:"logoUrl"`
	Status            int    `json:"status"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizeURL      string `json:"authorizeUrl"`
	TokenURL          string `json:"tokenUrl"`
	ProfileURL        string `json:"profileUrl"`
	Scope             string `json:"scope"`
	JSONUserIDPath    string `json:"jsonUserIdPath"`
	JSONUserNamePath  string `json:"jsonUserNamePath"`
	JSONUserEmailPath string `json:"jsonUserEmailPath"`
}

var (
	//OAuthConfigDisabled is used to disable an OAuthConfig for signin
	OAuthConfigDisabled = 1
	//OAuthConfigEnabled is used to enable an OAuthConfig for public use
	OAuthConfigEnabled = 2
)
