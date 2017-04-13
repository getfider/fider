package identity

import "github.com/WeCanHearYou/wechy/app/models"

// UserService is used for user operations
type UserService interface {
	GetByID(userID int) (*models.User, error)
	GetByEmail(tenantID int, email string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
}

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) (*models.Tenant, error)
}
