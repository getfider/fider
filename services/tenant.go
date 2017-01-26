package services

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy-api/models"
)

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) (*models.Tenant, error)
}

// InMemoryTenantService contains read and write operations for tenants
type InMemoryTenantService struct {
	Tenants []*models.Tenant
}

// GetByDomain returns a tenant based on its domain
func (svc InMemoryTenantService) GetByDomain(domain string) (*models.Tenant, error) {
	for _, tenant := range svc.Tenants {
		if tenant.Domain == domain {
			return tenant, nil
		}
	}
	return nil, ErrNotFound
}

// PostgresTenantService contains read and write operations for tenants
type PostgresTenantService struct {
	DB *sql.DB
}

// GetByDomain returns a tenant based on its domain
func (svc PostgresTenantService) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	row := svc.DB.QueryRow("SELECT id, name, domain FROM tenants WHERE domain = $1", domain)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return tenant, nil
}
