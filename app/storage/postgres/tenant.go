package postgres

import (
	"strings"

	"database/sql"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	DB *dbx.Database
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := models.Tenant{}

	err := s.DB.Get(&tenant, "SELECT id, name, subdomain FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}

	return &tenant, nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
