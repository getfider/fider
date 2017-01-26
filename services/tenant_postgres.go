package services

import (
	"database/sql"

	"fmt"

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
func (svc PostgresTenantService) GetByDomain(domain string) models.Tenant {
	var tenant models.Tenant
	fmt.Println(domain)
	err := svc.db.QueryRow("SELECT id, name, domain FROM tenants WHERE domain = $1", domain).Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
	if err != nil {
		//TODO: proper error handling
		fmt.Println(err)
	}
	return tenant
}
