package postgres

import (
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
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

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	trx     *dbx.Trx
	current *models.Tenant
	user    *models.User
}

// NewTenantStorage creates a new TenantStorage
func NewTenantStorage(trx *dbx.Trx) *TenantStorage {
	return &TenantStorage{trx: trx}
}

// SetCurrentTenant to current context
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
}

// SetCurrentUser to current context
func (s *TenantStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// Current returns current context tenant
func (s *TenantStorage) Current() *models.Tenant {
	return s.current
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string, status int) (*models.Tenant, error) {
	now := time.Now()

	var id int
	err := s.trx.Get(&id,
		`INSERT INTO tenants (name, subdomain, created_at, cname, invitation, welcome_message, status, is_private, custom_css, logo_bkey) 
		 VALUES ($1, $2, $3, '', '', '', $4, false, '', '') 
		 RETURNING id`, name, subdomain, now, status)
	if err != nil {
		return nil, err
	}

	if env.IsBillingEnabled() {
		_, err = s.trx.Execute(
			`INSERT INTO tenants_billing (tenant_id, trial_ends_at) VALUES ($1, $2)`,
			id, now.Add(30*24*time.Hour),
		)
		if err != nil {
			return nil, err
		}
	}

	return s.GetByDomain(subdomain)
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, `
		SELECT t.id, t.name, t.subdomain, t.cname, t.invitation, t.welcome_message, t.status, t.is_private, t.logo_bkey, t.custom_css,
					 tb.trial_ends_at AS billing_trial_ends_at,
					 tb.subscription_ends_at AS billing_subscription_ends_at,
					 tb.stripe_customer_id AS billing_stripe_customer_id,
					 tb.stripe_plan_id AS billing_stripe_plan_id,
					 tb.stripe_subscription_id AS billing_stripe_subscription_id
		FROM tenants t
		LEFT JOIN tenants_billing tb
		ON tb.tenant_id = t.id
		ORDER BY t.id LIMIT 1
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get first tenant")
	}

	return tenant.toModel(), nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, `
		SELECT t.id, t.name, t.subdomain, t.cname, t.invitation, t.welcome_message, t.status, t.is_private, t.logo_bkey, t.custom_css,
					 tb.trial_ends_at AS billing_trial_ends_at,
					 tb.subscription_ends_at AS billing_subscription_ends_at,
					 tb.stripe_customer_id AS billing_stripe_customer_id,
					 tb.stripe_plan_id AS billing_stripe_plan_id,
					 tb.stripe_subscription_id AS billing_stripe_subscription_id
		FROM tenants t
		LEFT JOIN tenants_billing tb
		ON tb.tenant_id = t.id
		WHERE t.subdomain = $1 OR t.subdomain = $2 OR t.cname = $3 
		ORDER BY t.cname DESC
	`, env.Subdomain(domain), domain, domain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tenant with domain '%s'", domain)
	}

	return tenant.toModel(), nil
}

// UpdateSettings of current tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	if settings.Logo.Remove {
		settings.Logo.BlobKey = ""
	}

	query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3, cname = $4, logo_bkey = $5 WHERE id = $6"
	_, err := s.trx.Execute(query, settings.Title, settings.Invitation, settings.WelcomeMessage, settings.CNAME, settings.Logo.BlobKey, s.current.ID)
	if err != nil {
		return errors.Wrap(err, "failed update tenant settings")
	}

	s.current.Name = settings.Title
	s.current.Invitation = settings.Invitation
	s.current.CNAME = settings.CNAME
	s.current.WelcomeMessage = settings.WelcomeMessage

	return nil
}

// UpdateBillingSettings of current tenant
func (s *TenantStorage) UpdateBillingSettings(billing *models.TenantBilling) error {
	_, err := s.trx.Execute(`
		UPDATE tenants_billing 
		SET stripe_customer_id = $1, stripe_plan_id = $2, stripe_subscription_id = $3, 
			subscription_ends_at = $4 
		WHERE tenant_id = $5
	`, billing.StripeCustomerID, billing.StripePlanID, billing.StripeSubscriptionID,
		billing.SubscriptionEndsAt, s.current.ID)
	if err != nil {
		return errors.Wrap(err, "failed update tenant billing settings")
	}
	return nil
}

// UpdateAdvancedSettings of current tenant
func (s *TenantStorage) UpdateAdvancedSettings(settings *models.UpdateTenantAdvancedSettings) error {
	query := "UPDATE tenants SET custom_css = $1 WHERE id = $2"
	_, err := s.trx.Execute(query, settings.CustomCSS, s.current.ID)
	if err != nil {
		return errors.Wrap(err, "failed update tenant advanced settings")
	}
	s.current.CustomCSS = settings.CustomCSS
	return nil
}

// UpdatePrivacy settings of current tenant
func (s *TenantStorage) UpdatePrivacy(settings *models.UpdateTenantPrivacy) error {
	query := "UPDATE tenants SET is_private = $1 WHERE id = $2"
	_, err := s.trx.Execute(query, settings.IsPrivate, s.current.ID)
	if err != nil {
		return errors.Wrap(err, "failed update tenant privacy settings")
	}
	return nil
}

// IsSubdomainAvailable returns true if subdomain is available to use
func (s *TenantStorage) IsSubdomainAvailable(subdomain string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", subdomain)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if tenant exists with subdomain '%s'", subdomain)
	}
	return !exists, nil
}

// IsCNAMEAvailable returns true if cname is available to use
func (s *TenantStorage) IsCNAMEAvailable(cname string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE cname = $1 AND id <> $2", cname, s.current.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if tenant exists with CNAME '%s'", cname)
	}
	return !exists, nil
}

// Activate given tenant
func (s *TenantStorage) Activate(id int) error {
	query := "UPDATE tenants SET status = $1 WHERE id = $2"
	_, err := s.trx.Execute(query, models.TenantActive, id)
	if err != nil {
		return errors.Wrap(err, "failed to activate tenant with id '%d'", id)
	}
	return nil
}

// SaveVerificationKey used by email verification process
func (s *TenantStorage) SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error {
	var userID interface{}
	if request.GetUser() != nil {
		userID = request.GetUser().ID
	}

	query := "INSERT INTO email_verifications (tenant_id, email, created_at, expires_at, key, name, kind, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.trx.Execute(query, s.current.ID, request.GetEmail(), time.Now(), time.Now().Add(duration), key, request.GetName(), request.GetKind(), userID)
	if err != nil {
		return errors.Wrap(err, "failed to save verification key for kind '%d'", request.GetKind())
	}
	return nil
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error) {
	verification := dbEmailVerification{}

	query := "SELECT id, email, name, key, created_at, verified_at, expires_at, kind, user_id FROM email_verifications WHERE key = $1 AND kind = $2 LIMIT 1"
	err := s.trx.Get(&verification, query, key, kind)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get email verification by its key")
	}

	return verification.toModel(), nil
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	query := "UPDATE email_verifications SET verified_at = $1 WHERE tenant_id = $2 AND key = $3"
	_, err := s.trx.Execute(query, time.Now(), s.current.ID, key)
	if err != nil {
		return errors.Wrap(err, "failed to update verified date of email verification request")
	}
	return nil
}

// SaveOAuthConfig saves given config into database
func (s *TenantStorage) SaveOAuthConfig(config *models.CreateEditOAuthConfig) error {
	var err error

	if config.Logo.Remove {
		config.Logo.BlobKey = ""
	}

	if config.ID == 0 {
		query := `INSERT INTO oauth_providers (
			tenant_id, provider, display_name, status,
			client_id, client_secret, authorize_url,
			profile_url, token_url, scope, json_user_id_path,
			json_user_name_path, json_user_email_path, logo_bkey
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

		err = s.trx.Get(&config.ID, query, s.current.ID, config.Provider,
			config.DisplayName, config.Status, config.ClientID, config.ClientSecret,
			config.AuthorizeURL, config.ProfileURL, config.TokenURL,
			config.Scope, config.JSONUserIDPath, config.JSONUserNamePath,
			config.JSONUserEmailPath, config.Logo.BlobKey)
	} else {
		query := `
			UPDATE oauth_providers 
			SET display_name = $3, status = $4, client_id = $5, client_secret = $6, 
					authorize_url = $7, profile_url = $8, token_url = $9, scope = $10, 
					json_user_id_path = $11, json_user_name_path = $12, json_user_email_path = $13,
					logo_bkey = $14
		WHERE tenant_id = $1 AND id = $2`

		_, err = s.trx.Execute(query, s.current.ID, config.ID,
			config.DisplayName, config.Status, config.ClientID, config.ClientSecret,
			config.AuthorizeURL, config.ProfileURL, config.TokenURL,
			config.Scope, config.JSONUserIDPath, config.JSONUserNamePath,
			config.JSONUserEmailPath, config.Logo.BlobKey)
	}

	if err != nil {
		return errors.Wrap(err, "failed to save OAuth Provider")
	}

	return nil
}

// GetOAuthConfigByProvider returns a custom OAuth configuration by provider name
func (s *TenantStorage) GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error) {
	if s.current == nil {
		return nil, app.ErrNotFound
	}

	config := &dbOAuthConfig{}
	err := s.trx.Get(config, `
	SELECT id, provider, display_name, status, logo_bkey,
				 client_id, client_secret, authorize_url,
				 profile_url, token_url, scope, json_user_id_path,
				 json_user_name_path, json_user_email_path
	FROM oauth_providers
	WHERE tenant_id = $1 AND provider = $2
	`, s.current.ID, provider)
	if err != nil {
		return nil, err
	}
	return config.toModel(), nil
}

// ListOAuthConfig returns a list of all custom OAuth provider for current tenant
func (s *TenantStorage) ListOAuthConfig() ([]*models.OAuthConfig, error) {
	configs := []*dbOAuthConfig{}
	if s.current != nil {
		err := s.trx.Select(&configs, `
		SELECT id, provider, display_name, status, logo_bkey,
					 client_id, client_secret, authorize_url,
					 profile_url, token_url, scope, json_user_id_path,
					 json_user_name_path, json_user_email_path
		FROM oauth_providers
		WHERE tenant_id = $1
		ORDER BY id
		`, s.current.ID)
		if err != nil {
			return nil, err
		}
	}

	var result = make([]*models.OAuthConfig, len(configs))
	for i, config := range configs {
		result[i] = config.toModel()
	}
	return result, nil
}
