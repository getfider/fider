package identity

// UserService is used for user operations
type UserService interface {
	GetByEmail(email string) (*User, error)
	Register(user *User) error
}

// TenantService contains read and write operations for tenants
type TenantService interface {
	GetByDomain(domain string) (*Tenant, error)
}
