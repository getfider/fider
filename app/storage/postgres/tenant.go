package postgres

import (
	"database/sql"
	"net/http"
	"strings"
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
		`INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message, status, is_private) 
		 VALUES ($1, $2, $3, '', '', '', $4, false) 
		 RETURNING id`, name, subdomain, time.Now(), status)
	if err != nil {
		return nil, err
	}

	return s.GetByDomain(subdomain)
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status, is_private, logo_id FROM tenants ORDER BY id LIMIT 1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get first tenant")
	}

	return tenant.toModel(), nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status, is_private, logo_id FROM tenants WHERE subdomain = $1 OR cname = $2 ORDER BY cname DESC", extractSubdomain(domain), domain)
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

	if settings.Logo != nil && !settings.Logo.Ignore {
		var newLogoID sql.NullInt64

		if !settings.Logo.Remove && len(settings.Logo.Upload.Content) > 0 {
			err := s.trx.Get(&newLogoID, `
				INSERT INTO uploads (tenant_id, size, content_type, file, created_on)
				VALUES ($1, $2, $3, $4, $5) RETURNING id
				`, s.current.ID, len(settings.Logo.Upload.Content), http.DetectContentType(settings.Logo.Upload.Content), settings.Logo.Upload.Content, time.Now(),
			)
			if err != nil {
				return errors.Wrap(err, "failed to upload new tenant logo")
			}
		}

		query := "UPDATE tenants SET logo_id = $1 WHERE id = $2"
		_, err = s.trx.Execute(query, newLogoID, s.current.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant logo")
		}

		if s.current.LogoID > 0 {
			query := "DELETE FROM uploads WHERE id = $1 AND tenant_id = $2"
			_, err = s.trx.Execute(query, s.current.LogoID, s.current.ID)
			if err != nil {
				return errors.Wrap(err, "failed delete old tenant logo")
			}
		}
	}

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

// GetLogo returns tenant logo by id
func (s *TenantStorage) GetLogo(id int) (*models.Upload, error) {
	upload := &models.Upload{}
	err := s.trx.Get(upload, `
		SELECT content_type, size, file FROM tenants
		INNER JOIN uploads
		ON uploads.tenant_id = tenants.id
		AND uploads.id = tenants.logo_id
		WHERE tenants.id = $1 AND uploads.id = $2
	`, s.current.ID, id)
	if err == app.ErrNotFound {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logo from tenant")
	}
	return upload, nil
}

func extractSubdomain(hostname string) string {
	domain := env.MultiTenantDomain()
	if domain == "" {
		return hostname
	}

	return strings.Replace(hostname, domain, "", -1)
}
