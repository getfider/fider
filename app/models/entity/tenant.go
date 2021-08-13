package entity

import "github.com/getfider/fider/app/models/enum"

//Tenant represents a tenant
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
