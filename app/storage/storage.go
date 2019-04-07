package storage

import (
	"github.com/getfider/fider/app/models"
)

// Base is a generic storage base interface
type Base interface {
	SetCurrentTenant(tenant *models.Tenant)
	SetCurrentUser(user *models.User)
}

// Tenant contains read and write operations for tenants
type Tenant interface {
	Base

	Add(name string, subdomain string, status int) (*models.Tenant, error)
	First() (*models.Tenant, error)
	GetByDomain(domain string) (*models.Tenant, error)
}
