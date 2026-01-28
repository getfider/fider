package dbEntities

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services"
)

type Tenant struct {
	ID                  int    `db:"id"`
	Name                string `db:"name"`
	Subdomain           string `db:"subdomain"`
	CNAME               string `db:"cname"`
	Invitation          string `db:"invitation"`
	WelcomeMessage      string `db:"welcome_message"`
	WelcomeHeader       string `db:"welcome_header"`
	Status              int    `db:"status"`
	Locale              string `db:"locale"`
	IsPrivate           bool   `db:"is_private"`
	LogoBlobKey         string `db:"logo_bkey"`
	CustomCSS           string `db:"custom_css"`
	AllowedSchemes      string `db:"allowed_schemes"`
	IsEmailAuthAllowed  bool   `db:"is_email_auth_allowed"`
	IsFeedEnabled       bool   `db:"is_feed_enabled"`
	PreventIndexing     bool   `db:"prevent_indexing"`
	IsModerationEnabled bool   `db:"is_moderation_enabled"`
	IsPro               bool   `db:"is_pro"`
}

func (t *Tenant) ToModel() *entity.Tenant {
	if t == nil {
		return nil
	}

	// Compute hasCommercialFeatures based on hosting mode
	var hasCommercialFeatures bool
	if env.IsSingleHostMode() {
		// Self-hosted: check if license service validated successfully
		hasCommercialFeatures = services.IsCommercialFeatureEnabled(services.FeatureContentModeration)
	} else {
		// Hosted multi-tenant: check if this tenant has Pro subscription
		hasCommercialFeatures = t.IsPro
	}

	tenant := &entity.Tenant{
		ID:                    t.ID,
		Name:                  t.Name,
		Subdomain:             t.Subdomain,
		CNAME:                 t.CNAME,
		Invitation:            t.Invitation,
		WelcomeMessage:        t.WelcomeMessage,
		WelcomeHeader:         t.WelcomeHeader,
		Status:                enum.TenantStatus(t.Status),
		Locale:                t.Locale,
		IsPrivate:             t.IsPrivate,
		LogoBlobKey:           t.LogoBlobKey,
		CustomCSS:             t.CustomCSS,
		AllowedSchemes:        t.AllowedSchemes,
		IsEmailAuthAllowed:    t.IsEmailAuthAllowed,
		IsFeedEnabled:         t.IsFeedEnabled,
		PreventIndexing:       t.PreventIndexing,
		IsModerationEnabled:   t.IsModerationEnabled,
		HasCommercialFeatures: hasCommercialFeatures,
	}

	return tenant
}
