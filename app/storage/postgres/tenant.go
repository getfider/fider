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
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Subdomain string         `db:"subdomain"`
	CNAME     sql.NullString `db:"cname"`
}

func (t *dbTenant) toModel() *models.Tenant {
	return &models.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		Subdomain: t.Subdomain,
		CNAME:     t.CNAME.String,
	}
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(tenant *models.Tenant) error {
	row := s.trx.QueryRow(`INSERT INTO tenants (name, subdomain, cname, created_on) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`, tenant.Name, tenant.Subdomain, tenant.CNAME, time.Now())
	if err := row.Scan(&tenant.ID); err != nil {
		return err
	}

	return nil
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname FROM tenants LIMIT 1")
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

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
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
