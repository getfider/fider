package dbEntities

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
)

type Tenant struct {
	ID                    int    `db:"id"`
	Name                  string `db:"name"`
	Subdomain             string `db:"subdomain"`
	CNAME                 string `db:"cname"`
	Invitation            string `db:"invitation"`
	WelcomeMessage        string `db:"welcome_message"`
	WelcomeHeader         string `db:"welcome_header"`
	DescriptionTemplate   string `db:"description_template"`
	Status                int    `db:"status"`
	Locale                string `db:"locale"`
	IsPrivate             bool   `db:"is_private"`
	LogoBlobKey           string `db:"logo_bkey"`
	CustomCSS             string `db:"custom_css"`
	AllowedSchemes        string `db:"allowed_schemes"`
	IsEmailAuthAllowed    bool   `db:"is_email_auth_allowed"`
	IsFeedEnabled         bool   `db:"is_feed_enabled"`
	PreventIndexing       bool   `db:"prevent_indexing"`
	IsModerationEnabled   bool   `db:"is_moderation_enabled"`
	IsPro                 bool         `db:"is_pro"`
	HasPaddleSubscription bool         `db:"has_paddle_subscription"`
	ScheduledDeletionAt   dbx.NullTime `db:"scheduled_deletion_at"`
	SiteBannerEnabled     bool         `db:"site_banner_enabled"`
	SiteBannerMessage     string       `db:"site_banner_message"`
	SiteBannerVariant     string       `db:"site_banner_variant"`
}

func (t *Tenant) ToModel() *entity.Tenant {
	if t == nil {
		return nil
	}

	// Self-hosted: all features are available (isPro = true)
	// Hosted multi-tenant: isPro based on subscription status
	isPro := true
	if env.IsMultiHostMode() {
		isPro = t.IsPro || t.HasPaddleSubscription
	}

	tenant := &entity.Tenant{
		ID:                  t.ID,
		Name:                t.Name,
		Subdomain:           t.Subdomain,
		CNAME:               t.CNAME,
		Invitation:          t.Invitation,
		WelcomeMessage:      t.WelcomeMessage,
		WelcomeHeader:       t.WelcomeHeader,
		DescriptionTemplate: t.DescriptionTemplate,
		Status:              enum.TenantStatus(t.Status),
		Locale:              t.Locale,
		IsPrivate:           t.IsPrivate,
		LogoBlobKey:         t.LogoBlobKey,
		CustomCSS:           t.CustomCSS,
		AllowedSchemes:      t.AllowedSchemes,
		IsEmailAuthAllowed:  t.IsEmailAuthAllowed,
		IsFeedEnabled:       t.IsFeedEnabled,
		PreventIndexing:     t.PreventIndexing,
		IsModerationEnabled: isPro && t.IsModerationEnabled,
		IsPro:               isPro,
		SiteBannerEnabled:   t.SiteBannerEnabled,
		SiteBannerMessage:   t.SiteBannerMessage,
		SiteBannerVariant:   t.SiteBannerVariant,
	}

	if t.ScheduledDeletionAt.Valid {
		tenant.ScheduledDeletionAt = &t.ScheduledDeletionAt.Time
	}

	return tenant
}
