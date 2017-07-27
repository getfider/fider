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
	trx *dbx.Trx
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

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
}

// UpdateSettings of given tenant
func (s *TenantStorage) UpdateSettings(tenantID int, title, invitation, welcomeMessage string) error {
	query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3 WHERE id = $4"
	return s.trx.Execute(query, title, invitation, welcomeMessage, tenantID)
}

// IsSubdomainAvailable returns true if subdomain is available to use
func (s *TenantStorage) IsSubdomainAvailable(subdomain string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", subdomain)
	return !exists, err
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
