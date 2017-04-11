package identity

import "github.com/WeCanHearYou/wechy/app"

// UserService is used for user operations
type UserService interface {
	GetByID(userID int) (*app.User, error)
	GetByEmail(tenantID int, email string) (*app.User, error)
	Register(user *app.User) error
	RegisterProvider(userID int, provider *app.UserProvider) error
}

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) (*app.Tenant, error)
}
