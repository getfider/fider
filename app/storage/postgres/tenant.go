package postgres

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbTenant struct {
	ID             int         `db:"id"`
	Name           string      `db:"name"`
	Subdomain      string      `db:"subdomain"`
	CNAME          string      `db:"cname"`
	Invitation     string      `db:"invitation"`
	WelcomeMessage string      `db:"welcome_message"`
	Status         int         `db:"status"`
	IsPrivate      bool        `db:"is_private"`
	LogoID         dbx.NullInt `db:"logo_id"`
	CustomCSS      string      `db:"custom_css"`
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
		CustomCSS:      t.CustomCSS,
	}

	if t.LogoID.Valid {
		tenant.LogoID = int(t.LogoID.Int64)
	}

	return tenant
}

type dbEmailVerification struct {
	ID         int                          `db:"id"`
	Name       string                       `db:"name"`
	Email      string                       `db:"email"`
	Key        string                       `db:"key"`
	Kind       models.EmailVerificationKind `db:"kind"`
	UserID     dbx.NullInt                  `db:"user_id"`
	CreatedOn  time.Time                    `db:"created_on"`
	ExpiresOn  time.Time                    `db:"expires_on"`
	VerifiedOn dbx.NullTime                 `db:"verified_on"`
}

func (t *dbEmailVerification) toModel() *models.EmailVerification {
	model := &models.EmailVerification{
		Name:       t.Name,
		Email:      t.Email,
		Key:        t.Key,
		Kind:       t.Kind,
		CreatedOn:  t.CreatedOn,
		ExpiresOn:  t.ExpiresOn,
		VerifiedOn: nil,
	}

	if t.VerifiedOn.Valid {
		model.VerifiedOn = &t.VerifiedOn.Time
	}

	if t.UserID.Valid {
		model.UserID = int(t.UserID.Int64)
	}

	return model
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

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string, status int) (*models.Tenant, error) {
	var id int
	err := s.trx.Get(&id,
		`INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message, status, is_private, custom_css) 
		 VALUES ($1, $2, $3, '', '', '', $4, false, '') 
		 RETURNING id`, name, subdomain, time.Now(), status)
	if err != nil {
		return nil, err
	}

	return s.GetByDomain(subdomain)
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status, is_private, logo_id, custom_css FROM tenants ORDER BY id LIMIT 1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get first tenant")
	}

	return tenant.toModel(), nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status, is_private, logo_id, custom_css FROM tenants WHERE subdomain = $1 OR subdomain = $2 OR cname = $3 ORDER BY cname DESC", env.Subdomain(domain), domain, domain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tenant with domain '%s'", domain)
	}

	return tenant.toModel(), nil
}

// UpdateSettings of current tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3, cname = $4 WHERE id = $5"
	_, err := s.trx.Execute(query, settings.Title, settings.Invitation, settings.WelcomeMessage, settings.CNAME, s.current.ID)
	if err != nil {
		return errors.Wrap(err, "failed update tenant settings")
	}

	s.current.Name = settings.Title
	s.current.Invitation = settings.Invitation
	s.current.CNAME = settings.CNAME
	s.current.WelcomeMessage = settings.WelcomeMessage

	// Only update tenant logo if there's a change
	if settings.Logo != nil {
		var newLogoID sql.NullInt64

		// If there's an upload, save it
		if !settings.Logo.Remove && len(settings.Logo.Upload.Content) > 0 {
			uploadID, err := s.SaveNewUpload(settings.Logo.Upload.Content)
			if err != nil {
				return errors.Wrap(err, "failed to upload new tenant logo")
			}
			newLogoID.Scan(uploadID)
		}

		// Update current tenant logo to either new ID or null
		query := "UPDATE tenants SET logo_id = $1 WHERE id = $2"
		_, err = s.trx.Execute(query, newLogoID, s.current.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant logo")
		}

		// Remove upload based on old tenant logo ID
		err := s.RemoveUpload(s.current.LogoID)
		if err != nil {
			return errors.Wrap(err, "failed delete old tenant logo")
		}

		// Update reference to new logo ID
		s.current.LogoID = int(newLogoID.Int64)
	}

	return nil
}

// SaveNewUpload saves given content and return upload id
func (s *TenantStorage) SaveNewUpload(content []byte) (int, error) {
	var newID int
	err := s.trx.Get(&newID, `
		INSERT INTO uploads (tenant_id, size, content_type, file, created_on)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, s.current.ID, len(content), http.DetectContentType(content), content, time.Now(),
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to save new upload")
	}
	return newID, nil
}

// RemoveUpload by given id
func (s *TenantStorage) RemoveUpload(uploadID int) error {
	if uploadID > 0 {
		query := "DELETE FROM uploads WHERE id = $1 AND tenant_id = $2"
		_, err := s.trx.Execute(query, uploadID, s.current.ID)
		if err != nil {
			return errors.Wrap(err, "failed delete upload")
		}
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

	query := "INSERT INTO email_verifications (tenant_id, email, created_on, expires_on, key, name, kind, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.trx.Execute(query, s.current.ID, request.GetEmail(), time.Now(), time.Now().Add(duration), key, request.GetName(), request.GetKind(), userID)
	if err != nil {
		return errors.Wrap(err, "failed to save verification key for kind '%d'", request.GetKind())
	}
	return nil
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error) {
	verification := dbEmailVerification{}

	query := "SELECT id, email, name, key, created_on, verified_on, expires_on, kind, user_id FROM email_verifications WHERE key = $1 AND kind = $2 LIMIT 1"
	err := s.trx.Get(&verification, query, key, kind)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get email verification by its key")
	}

	return verification.toModel(), nil
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	query := "UPDATE email_verifications SET verified_on = $1 WHERE tenant_id = $2 AND key = $3"
	_, err := s.trx.Execute(query, time.Now(), s.current.ID, key)
	if err != nil {
		return errors.Wrap(err, "failed to update verified date of email verification request")
	}
	return nil
}

// GetUpload returns upload by id
func (s *TenantStorage) GetUpload(id int) (*models.Upload, error) {
	upload := &models.Upload{}
	err := s.trx.Get(upload, `
		SELECT content_type, size, file 
		FROM uploads
		WHERE tenant_id = $1 AND id = $2
	`, s.current.ID, id)
	if err == app.ErrNotFound {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get upload from tenant")
	}
	return upload, nil
}

// SaveOAuthConfig saves given config into database
func (s *TenantStorage) SaveOAuthConfig(config *models.CreateEditOAuthConfig) error {
	if config.ID == 0 {
		query := `INSERT INTO oauth_providers (
			tenant_id, provider, display_name, status,
			client_id, client_secret, authorize_url,
			profile_url, token_url, scope, json_user_id_path,
			json_user_name_path, json_user_email_path
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

		_, err := s.trx.Execute(query, s.current.ID, config.Provider,
			config.DisplayName, 0, config.ClientID, config.ClientSecret,
			config.AuthorizeURL, config.ProfileURL, config.TokenURL,
			config.Scope, config.JSONUserIDPath, config.JSONUserNamePath,
			config.JSONUserEmailPath)
		if err != nil {
			return errors.Wrap(err, "failed to insert OAuth Provider")
		}
	} else {
		query := `
			UPDATE oauth_providers 
			SET display_name = $3, status = $4, client_id = $5, client_secret = $6, 
					authorize_url = $7, profile_url = $8, token_url = $9, scope = $10, 
					json_user_id_path = $11, json_user_name_path = $12, json_user_email_path = $13
		WHERE tenant_id = $1 AND id = $2`

		_, err := s.trx.Execute(query, s.current.ID, config.ID,
			config.DisplayName, 0, config.ClientID, config.ClientSecret,
			config.AuthorizeURL, config.ProfileURL, config.TokenURL,
			config.Scope, config.JSONUserIDPath, config.JSONUserNamePath,
			config.JSONUserEmailPath)
		if err != nil {
			return errors.Wrap(err, "failed to update OAuth Provider")
		}
	}
	return nil
}

// GetOAuthConfigByProvider returns a custom OAuth configuration by provider name
func (s *TenantStorage) GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error) {
	if s.current == nil {
		return nil, app.ErrNotFound
	}

	config := &models.OAuthConfig{}
	err := s.trx.Get(config, `
	SELECT id, provider, display_name, status,
				 client_id, client_secret, authorize_url,
				 profile_url, token_url, scope, json_user_id_path,
				 json_user_name_path, json_user_email_path
	FROM oauth_providers
	WHERE tenant_id = $1 AND provider = $2
	`, s.current.ID, provider)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ListOAuthConfig returns a list of all custom OAuth provider for current tenant
func (s *TenantStorage) ListOAuthConfig() ([]*models.OAuthConfig, error) {
	entries := []*models.OAuthConfig{}
	if s.current != nil {
		err := s.trx.Select(&entries, `
		SELECT id, provider, display_name, status,
					 client_id, client_secret, authorize_url,
					 profile_url, token_url, scope, json_user_id_path,
					 json_user_name_path, json_user_email_path
		FROM oauth_providers
		WHERE tenant_id = $1
		`, s.current.ID)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}
