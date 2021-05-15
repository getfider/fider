package entity

import "encoding/json"

// OAuthConfig is the configuration of a custom OAuth provider
type OAuthConfig struct {
	ID                int
	Provider          string
	DisplayName       string
	LogoBlobKey       string
	Status            int
	ClientID          string
	ClientSecret      string
	AuthorizeURL      string
	TokenURL          string
	ProfileURL        string
	Scope             string
	JSONUserIDPath    string
	JSONUserNamePath  string
	JSONUserEmailPath string
}

// MarshalJSON returns the JSON encoding of OAuthConfig
func (o OAuthConfig) MarshalJSON() ([]byte, error) {
	secret := "..."
	if len(o.ClientSecret) >= 10 {
		secret = o.ClientSecret[0:3] + "..." + o.ClientSecret[len(o.ClientSecret)-3:]
	}
	return json.Marshal(map[string]interface{}{
		"id":                o.ID,
		"provider":          o.Provider,
		"displayName":       o.DisplayName,
		"logoBlobKey":       o.LogoBlobKey,
		"status":            o.Status,
		"clientID":          o.ClientID,
		"clientSecret":      secret,
		"authorizeURL":      o.AuthorizeURL,
		"tokenURL":          o.TokenURL,
		"profileURL":        o.ProfileURL,
		"scope":             o.Scope,
		"jsonUserIDPath":    o.JSONUserIDPath,
		"jsonUserNamePath":  o.JSONUserNamePath,
		"jsonUserEmailPath": o.JSONUserEmailPath,
	})
}
