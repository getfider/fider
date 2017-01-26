package services

import "github.com/WeCanHearYou/wchy-api/models"

// InMemoryTenantService contains read and write operations for tenants
type InMemoryTenantService struct {
	tenants []*models.Tenant
}

// NewInMemoryTenantService creates a new InMemoryTenantService
func NewInMemoryTenantService(tenants []*models.Tenant) *InMemoryTenantService {
	return &InMemoryTenantService{tenants}
}

// GetByDomain returns a tenant based on its domain
func (svc InMemoryTenantService) GetByDomain(domain string) *models.Tenant {
	for _, tenant := range svc.tenants {
		if tenant.Domain == domain {
			return tenant
		}
	}
	return nil
}
