package postgres

import (
	"strings"
	"time"

	"database/sql"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	DB *dbx.Database
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(tenant *models.Tenant) error {
	row := s.DB.QueryRow(`INSERT INTO tenants (name, subdomain, cname, created_on) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`, tenant.Name, tenant.Subdomain, tenant.CNAME, time.Now())
	if err := row.Scan(&tenant.ID); err != nil {
		return err
	}

	return nil
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := models.Tenant{}

	err := s.DB.Get(&tenant, "SELECT id, name, subdomain, cname FROM tenants LIMIT 1")
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := models.Tenant{}

	err := s.DB.Get(&tenant, "SELECT id, name, subdomain, cname FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
