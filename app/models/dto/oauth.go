package dto

//OAuthUserProfile represents an OAuth user profile
type OAuthUserProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//OAuthProviderOption represents an OAuth provider that can be used to authenticate
type OAuthProviderOption struct {
	Provider         string `json:"provider"`
	DisplayName      string `json:"displayName"`
	ClientID         string `json:"clientID"`
	URL              string `json:"url"`
	CallbackURL      string `json:"callbackURL"`
	LogoBlobKey      string `json:"logoBlobKey"`
	IsCustomProvider bool   `json:"isCustomProvider"`
	IsEnabled        bool   `json:"isEnabled"`
}
