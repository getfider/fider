package dbEntities

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type Tenant struct {
	ID                  int    `db:"id"`
	Name                string `db:"name"`
	Subdomain           string `db:"subdomain"`
	CNAME               string `db:"cname"`
	Invitation          string `db:"invitation"`
	WelcomeMessage      string `db:"welcome_message"`
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
}

func (t *Tenant) ToModel() *entity.Tenant {
	if t == nil {
		return nil
	}

	tenant := &entity.Tenant{
		ID:                  t.ID,
		Name:                t.Name,
		Subdomain:           t.Subdomain,
		CNAME:               t.CNAME,
		Invitation:          t.Invitation,
		WelcomeMessage:      t.WelcomeMessage,
		Status:              enum.TenantStatus(t.Status),
		Locale:              t.Locale,
		IsPrivate:           t.IsPrivate,
		LogoBlobKey:         t.LogoBlobKey,
		CustomCSS:           t.CustomCSS,
		AllowedSchemes:      t.AllowedSchemes,
		IsEmailAuthAllowed:  t.IsEmailAuthAllowed,
		IsFeedEnabled:       t.IsFeedEnabled,
		PreventIndexing:     t.PreventIndexing,
		IsModerationEnabled: t.IsModerationEnabled,
	}

	return tenant
}
