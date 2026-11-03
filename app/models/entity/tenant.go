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
	IsSignUpDisabled    bool              `json:"isSignUpDisabled"`
	IsPro               bool              `json:"isPro"`
	// ScheduledDeletionAt is set when the account owner has requested deletion of the whole
	// site. The tenant stays active during the grace window; a background job performs the
	// hard delete once this time passes. Not exposed to clients.
	ScheduledDeletionAt *time.Time `json:"-"`
	// DeletionCancelKey authorises the emailed cancel link, so it is only ever set when the
	// deletion was requested by the account owner. An administrative deletion leaves it empty,
	// which is how the deletion job tells the two apart. Not exposed to clients.
	DeletionCancelKey string `json:"-"`
}

func (t *Tenant) IsDisabled() bool {
	return t.Status == enum.TenantDisabled
}

// CanCreateAccountByEmail returns true if completing an email verification of the given kind
// may create a new account. Invited users can always create their account, everyone else
// only if the site allows sign-ups and email login is turned on.
func (t *Tenant) CanCreateAccountByEmail(kind enum.EmailVerificationKind) bool {
	if kind == enum.EmailVerificationKindUserInvitation {
		return true
	}
	return t.IsEmailAuthAllowed && t.isSignUpAllowed()
}

// CanCreateAccountByOAuth returns true if signing in with the given OAuth provider may create
// a new account. Users from trusted providers can always create their account, everyone else
// only if the site allows sign-ups. Built-in providers have no config and are never trusted.
func (t *Tenant) CanCreateAccountByOAuth(config *OAuthConfig) bool {
	if config != nil && config.IsTrusted {
		return true
	}
	return t.isSignUpAllowed()
}

// isSignUpAllowed returns true if visitors can create a new account on their own,
// which is neither the case on private sites nor when sign-ups are disabled.
func (t *Tenant) isSignUpAllowed() bool {
	return !t.IsPrivate && !t.IsSignUpDisabled
}

// TenantContact is a reference to an administrator account
type TenantContact struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subdomain string `json:"subdomain"`
}
