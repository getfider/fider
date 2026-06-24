package entity

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

// Tenant represents a tenant
type Tenant struct {
	ID                  int               `json:"id"`
	Name                string            `json:"name"`
	Subdomain           string            `json:"subdomain"`
	Invitation          string            `json:"invitation"`
	WelcomeMessage      string            `json:"welcomeMessage"`
	WelcomeHeader       string            `json:"welcomeHeader"`
	DescriptionTemplate string            `json:"descriptionTemplate"`
	CNAME               string            `json:"cname"`
	Status              enum.TenantStatus `json:"status"`
	Locale              string            `json:"locale"`
	IsPrivate           bool              `json:"isPrivate"`
	LogoBlobKey         string            `json:"logoBlobKey"`
	CustomCSS           string            `json:"-"`
	AllowedSchemes      string            `json:"allowedSchemes"`
	IsEmailAuthAllowed  bool              `json:"isEmailAuthAllowed"`
	IsFeedEnabled       bool              `json:"isFeedEnabled"`
	PreventIndexing     bool              `json:"preventIndexing"`
	IsModerationEnabled bool              `json:"isModerationEnabled"`
	IsPro               bool              `json:"isPro"`
	Statuses            []*Status         `json:"statuses,omitempty"`
	// ScheduledDeletionAt is set when the account owner has requested deletion of the whole
	// site. The tenant stays active during the grace window; a background job performs the
	// hard delete once this time passes. Not exposed to clients.
	ScheduledDeletionAt *time.Time `json:"-"`
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
