package entity

import "github.com/getfider/fider/app/models/enum"

// Tenant represents a tenant
type Tenant struct {
	ID                 int               `json:"id"`
	Name               string            `json:"name"`
	Subdomain          string            `json:"subdomain"`
	Invitation         string            `json:"invitation"`
	WelcomeMessage     string            `json:"welcomeMessage"`
	CNAME              string            `json:"cname"`
	Status             enum.TenantStatus `json:"status"`
	Locale             string            `json:"locale"`
	IsPrivate          bool              `json:"isPrivate"`
	LogoBlobKey        string            `json:"logoBlobKey"`
	CustomCSS          string            `json:"-"`
	IsEmailAuthAllowed bool              `json:"isEmailAuthAllowed"`
}

func (t *Tenant) IsDisabled() bool {
	return t.Status == enum.TenantDisabled
}

// TenantContact is a reference to an administrator account
type TenantContact struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subdomain string `json:"subdomain"`
}
