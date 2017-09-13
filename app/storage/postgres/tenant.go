package postgres

import (
	"strings"
	"time"

	"database/sql"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	trx     *dbx.Trx
	current *models.Tenant
}

// NewTenantStorage creates a new TenantStorage
func NewTenantStorage(trx *dbx.Trx) *TenantStorage {
	return &TenantStorage{trx: trx}
}

type dbTenant struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	Subdomain      string `db:"subdomain"`
	CNAME          string `db:"cname"`
	Invitation     string `db:"invitation"`
	WelcomeMessage string `db:"welcome_message"`
}

func (t *dbTenant) toModel() *models.Tenant {
	return &models.Tenant{
		ID:             t.ID,
		Name:           t.Name,
		Subdomain:      t.Subdomain,
		CNAME:          t.CNAME,
		Invitation:     t.Invitation,
		WelcomeMessage: t.WelcomeMessage,
	}
}

type dbEmailVerification struct {
	ID         int          `db:"id"`
	Email      string       `db:"email"`
	Key        string       `db:"key"`
	CreatedOn  time.Time    `db:"created_on"`
	VerifiedOn dbx.NullTime `db:"verified_on"`
}

func (t *dbEmailVerification) toModel() *models.EmailVerification {
	model := &models.EmailVerification{
		Email:      t.Email,
		Key:        t.Key,
		CreatedOn:  t.CreatedOn,
		VerifiedOn: nil,
	}

	if t.VerifiedOn.Valid {
		model.VerifiedOn = &t.VerifiedOn.Time
	}

	return model
}

// SetCurrentTenant tenant
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string) (*models.Tenant, error) {
	var id int
	row := s.trx.QueryRow(`INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message) 
						VALUES ($1, $2, $3, '', '', '') 
						RETURNING id`, name, subdomain, time.Now())
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return s.GetByDomain(subdomain)
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message FROM tenants LIMIT 1")
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message FROM tenants WHERE subdomain = $1 OR cname = $2 ORDER BY cname DESC", extractSubdomain(domain), domain)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
}

// UpdateSettings of given tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3 WHERE id = $4"
	return s.trx.Execute(query, settings.Title, settings.Invitation, settings.WelcomeMessage, s.current.ID)
}

// IsSubdomainAvailable returns true if subdomain is available to use
func (s *TenantStorage) IsSubdomainAvailable(subdomain string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", subdomain)
	return !exists, err
}

// SaveVerificationKey used by e-mail verification
func (s *TenantStorage) SaveVerificationKey(email, key string) error {
	query := "INSERT INTO email_verifications (tenant_id, email, created_on, key) VALUES ($1, $2, $3, $4)"
	return s.trx.Execute(query, s.current.ID, email, time.Now(), key)
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(key string) (*models.EmailVerification, error) {
	verification := dbEmailVerification{}

	err := s.trx.Get(&verification, "SELECT id, email, key, created_on, verified_on FROM email_verifications WHERE key = $1 LIMIT 1", key)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return verification.toModel(), nil
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	query := "UPDATE email_verifications SET verified_on = $1 WHERE tenant_id = $2 AND key = $3"
	return s.trx.Execute(query, time.Now(), s.current.ID, key)
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
