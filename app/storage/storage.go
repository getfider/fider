package storage

import (
	"time"

	"github.com/getfider/fider/app/models"
)

// Base is a generic storage base interface
type Base interface {
	SetCurrentTenant(tenant *models.Tenant)
	SetCurrentUser(user *models.User)
}

// User is used for user operations
type User interface {
	Base
	GetByID(userID int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByProvider(provider string, uid string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
	Update(settings *models.UpdateUserSettings) error
	Delete() error
	ChangeEmail(userID int, email string) error
	ChangeRole(userID int, role models.Role) error
	GetAll() ([]*models.User, error)
	GetUserSettings() (map[string]string, error)
	UpdateSettings(settings map[string]string) error
	HasSubscribedTo(postID int) (bool, error)
	GetByAPIKey(apiKey string) (*models.User, error)
	RegenerateAPIKey() (string, error)
}

// Tenant contains read and write operations for tenants
type Tenant interface {
	Base
	Add(name string, subdomain string, status int) (*models.Tenant, error)
	First() (*models.Tenant, error)
	Activate(id int) error
	GetByDomain(domain string) (*models.Tenant, error)
	UpdateSettings(settings *models.UpdateTenantSettings) error
	UpdateBillingSettings(billing *models.TenantBilling) error
	UpdateAdvancedSettings(settings *models.UpdateTenantAdvancedSettings) error
	UpdatePrivacy(settings *models.UpdateTenantPrivacy) error
	IsSubdomainAvailable(subdomain string) (bool, error)
	IsCNAMEAvailable(cname string) (bool, error)
	SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error
	FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error)
	SetKeyAsVerified(key string) error
	SaveOAuthConfig(config *models.CreateEditOAuthConfig) error
	GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error)
	ListOAuthConfig() ([]*models.OAuthConfig, error)
}
