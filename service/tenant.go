package service

import (
	"database/sql"
	"strings"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/WeCanHearYou/wchy/model"
)

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) (*model.Tenant, error)
}

// PostgresTenantService contains read and write operations for tenants
type PostgresTenantService struct {
	DB *sql.DB
}

// GetByDomain returns a tenant based on its domain
func (svc PostgresTenantService) GetByDomain(domain string) (*model.Tenant, error) {
	tenant := &model.Tenant{}

	row := svc.DB.QueryRow("SELECT id, name, subdomain FROM tenants WHERE subdomain = $1", extractSubdomain(domain))
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	tenant.Domain = tenant.Domain + "." + env.GetCurrentDomain()
	return tenant, nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
