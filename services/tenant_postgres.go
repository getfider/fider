package services

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy-api/models"
)

// PostgresTenantService contains read and write operations for tenants
type PostgresTenantService struct {
	db *sql.DB
}

// NewPostgresTenantService creates a new PostgresTenantService
func NewPostgresTenantService(db *sql.DB) *PostgresTenantService {
	return &PostgresTenantService{db}
}

// GetByDomain returns a tenant based on its domain
func (svc PostgresTenantService) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	row := svc.db.QueryRow("SELECT id, name, domain FROM tenants WHERE domain = $1", domain)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return tenant, nil
}
