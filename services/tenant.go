package services

import "github.com/WeCanHearYou/wchy-api/models"

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) models.Tenant
}
