package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
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

func isCNAMEAvailable(ctx context.Context, q *query.IsCNAMEAvailable) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		exists, err := trx.Exists("SELECT id FROM tenants WHERE cname = $1 AND id <> $2", q.CNAME, tenant.ID)
		if err != nil {
			q.Result = false
			return errors.Wrap(err, "failed to check if tenant exists with CNAME '%s'", q.CNAME)
		}
		q.Result = !exists
		return nil
	})
}

func isSubdomainAvailable(ctx context.Context, q *query.IsSubdomainAvailable) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		exists, err := trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", q.Subdomain)
		if err != nil {
			q.Result = false
			return errors.Wrap(err, "failed to check if tenant exists with subdomain '%s'", q.Subdomain)
		}
		q.Result = !exists
		return nil
	})
}

func updateTenantPrivacySettings(ctx context.Context, c *cmd.UpdateTenantPrivacySettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		_, err := trx.Execute("UPDATE tenants SET is_private = $1 WHERE id = $2", c.Settings.IsPrivate, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant privacy settings")
		}
		return nil
	})
}

func updateTenantSettings(ctx context.Context, c *cmd.UpdateTenantSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if c.Settings.Logo.Remove {
			c.Settings.Logo.BlobKey = ""
		}

		query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3, cname = $4, logo_bkey = $5 WHERE id = $6"
		_, err := trx.Execute(query, c.Settings.Title, c.Settings.Invitation, c.Settings.WelcomeMessage, c.Settings.CNAME, c.Settings.Logo.BlobKey, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant settings")
		}

		tenant.Name = c.Settings.Title
		tenant.Invitation = c.Settings.Invitation
		tenant.CNAME = c.Settings.CNAME
		tenant.WelcomeMessage = c.Settings.WelcomeMessage

		return nil
	})
}

func updateTenantBillingSettings(ctx context.Context, c *cmd.UpdateTenantBillingSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		_, err := trx.Execute(`
			UPDATE tenants_billing 
			SET stripe_customer_id = $1, 
					stripe_plan_id = $2, 
					stripe_subscription_id = $3, 
					subscription_ends_at = $4 
			WHERE tenant_id = $5
		`,
			c.Settings.StripeCustomerID,
			c.Settings.StripePlanID,
			c.Settings.StripeSubscriptionID,
			c.Settings.SubscriptionEndsAt,
			tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed update tenant billing settings")
		}
		return nil
	})
}

func updateTenantAdvancedSettings(ctx context.Context, c *cmd.UpdateTenantAdvancedSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		query := "UPDATE tenants SET custom_css = $1 WHERE id = $2"
		_, err := trx.Execute(query, c.Settings.CustomCSS, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant advanced settings")
		}

		tenant.CustomCSS = c.Settings.CustomCSS
		return nil
	})
}

func activateTenant(ctx context.Context, c *cmd.ActivateTenant) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		query := "UPDATE tenants SET status = $1 WHERE id = $2"
		_, err := trx.Execute(query, models.TenantActive, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to activate tenant with id '%d'", c.TenantID)
		}
		return nil
	})
}

func getVerificationByKey(ctx context.Context, q *query.GetVerificationByKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		verification := dbEmailVerification{}

		query := "SELECT id, email, name, key, created_at, verified_at, expires_at, kind, user_id FROM email_verifications WHERE key = $1 AND kind = $2 LIMIT 1"
		err := trx.Get(&verification, query, q.Key, q.Kind)
		if err != nil {
			return errors.Wrap(err, "failed to get email verification by its key")
		}

		q.Result = verification.toModel()
		return nil
	})
}

func saveVerificationKey(ctx context.Context, c *cmd.SaveVerificationKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		var userID interface{}
		if c.Request.GetUser() != nil {
			userID = c.Request.GetUser().ID
		}

		query := "INSERT INTO email_verifications (tenant_id, email, created_at, expires_at, key, name, kind, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err := trx.Execute(query, tenant.ID, c.Request.GetEmail(), time.Now(), time.Now().Add(c.Duration), c.Key, c.Request.GetName(), c.Request.GetKind(), userID)
		if err != nil {
			return errors.Wrap(err, "failed to save verification key for kind '%d'", c.Request.GetKind())
		}
		return nil
	})
}

func setKeyAsVerified(ctx context.Context, c *cmd.SetKeyAsVerified) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		query := "UPDATE email_verifications SET verified_at = $1 WHERE tenant_id = $2 AND key = $3"
		_, err := trx.Execute(query, time.Now(), tenant.ID, c.Key)
		if err != nil {
			return errors.Wrap(err, "failed to update verified date of email verification request")
		}
		return nil
	})
}
