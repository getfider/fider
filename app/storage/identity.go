package storage

import "github.com/getfider/fider/app/models"

// User is used for user operations
type User interface {
	GetByID(userID int) (*models.User, error)
	GetByEmail(tenantID int, email string) (*models.User, error)
	GetByProvider(tenantID int, provider string, uid string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
}

// Tenant contains read and write operations for tenants
type Tenant interface {
	Add(name string, subdomain string) (*models.Tenant, error)
	First() (*models.Tenant, error)
	GetByDomain(domain string) (*models.Tenant, error)
	IsSubdomainAvailable(subdomain string) (bool, error)
}
