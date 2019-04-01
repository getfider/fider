package postgres

import (
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

type dbTenant struct {
	ID             int              `db:"id"`
	Name           string           `db:"name"`
	Subdomain      string           `db:"subdomain"`
	CNAME          string           `db:"cname"`
	Invitation     string           `db:"invitation"`
	WelcomeMessage string           `db:"welcome_message"`
	Status         int              `db:"status"`
	IsPrivate      bool             `db:"is_private"`
	LogoBlobKey    string           `db:"logo_bkey"`
	CustomCSS      string           `db:"custom_css"`
	Billing        *dbTenantBilling `db:"billing"`
}

func (t *dbTenant) toModel() *models.Tenant {
	if t == nil {
		return nil
	}

	tenant := &models.Tenant{
		ID:             t.ID,
		Name:           t.Name,
		Subdomain:      t.Subdomain,
		CNAME:          t.CNAME,
		Invitation:     t.Invitation,
		WelcomeMessage: t.WelcomeMessage,
		Status:         t.Status,
		IsPrivate:      t.IsPrivate,
		LogoBlobKey:    t.LogoBlobKey,
		CustomCSS:      t.CustomCSS,
	}

	if t.Billing != nil && t.Billing.TrialEndsAt.Valid {
		tenant.Billing = &models.TenantBilling{
			TrialEndsAt:          t.Billing.TrialEndsAt.Time,
			StripeCustomerID:     t.Billing.StripeCustomerID.String,
			StripeSubscriptionID: t.Billing.StripeSubscriptionID.String,
			StripePlanID:         t.Billing.StripePlanID.String,
		}
		if t.Billing.SubscriptionEndsAt.Valid {
			tenant.Billing.SubscriptionEndsAt = &t.Billing.SubscriptionEndsAt.Time
		}
	}

	return tenant
}

type dbTenantBilling struct {
	StripeCustomerID     dbx.NullString `db:"stripe_customer_id"`
	StripeSubscriptionID dbx.NullString `db:"stripe_subscription_id"`
	StripePlanID         dbx.NullString `db:"stripe_plan_id"`
	TrialEndsAt          dbx.NullTime   `db:"trial_ends_at"`
	SubscriptionEndsAt   dbx.NullTime   `db:"subscription_ends_at"`
}

type dbEmailVerification struct {
	ID         int                          `db:"id"`
	Name       string                       `db:"name"`
	Email      string                       `db:"email"`
	Key        string                       `db:"key"`
	Kind       models.EmailVerificationKind `db:"kind"`
	UserID     dbx.NullInt                  `db:"user_id"`
	CreatedAt  time.Time                    `db:"created_at"`
	ExpiresAt  time.Time                    `db:"expires_at"`
	VerifiedAt dbx.NullTime                 `db:"verified_at"`
}

func (t *dbEmailVerification) toModel() *models.EmailVerification {
	model := &models.EmailVerification{
		Name:       t.Name,
		Email:      t.Email,
		Key:        t.Key,
		Kind:       t.Kind,
		CreatedAt:  t.CreatedAt,
		ExpiresAt:  t.ExpiresAt,
		VerifiedAt: nil,
	}

	if t.VerifiedAt.Valid {
		model.VerifiedAt = &t.VerifiedAt.Time
	}

	if t.UserID.Valid {
		model.UserID = int(t.UserID.Int64)
	}

	return model
}

type dbOAuthConfig struct {
	ID                int    `db:"id"`
	Provider          string `db:"provider"`
	DisplayName       string `db:"display_name"`
	LogoBlobKey       string `db:"logo_bkey"`
	Status            int    `db:"status"`
	ClientID          string `db:"client_id"`
	ClientSecret      string `db:"client_secret"`
	AuthorizeURL      string `db:"authorize_url"`
	TokenURL          string `db:"token_url"`
	Scope             string `db:"scope"`
	ProfileURL        string `db:"profile_url"`
	JSONUserIDPath    string `db:"json_user_id_path"`
	JSONUserNamePath  string `db:"json_user_name_path"`
	JSONUserEmailPath string `db:"json_user_email_path"`
}

func (m *dbOAuthConfig) toModel() *models.OAuthConfig {
	return &models.OAuthConfig{
		ID:                m.ID,
		Provider:          m.Provider,
		DisplayName:       m.DisplayName,
		Status:            m.Status,
		LogoBlobKey:       m.LogoBlobKey,
		ClientID:          m.ClientID,
		ClientSecret:      m.ClientSecret,
		AuthorizeURL:      m.AuthorizeURL,
		TokenURL:          m.TokenURL,
		ProfileURL:        m.ProfileURL,
		Scope:             m.Scope,
		JSONUserIDPath:    m.JSONUserIDPath,
		JSONUserNamePath:  m.JSONUserNamePath,
		JSONUserEmailPath: m.JSONUserEmailPath,
	}
}
