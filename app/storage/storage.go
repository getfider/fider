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

// Post contains read and write operations for posts
type Post interface {
	Base
	GetByID(postID int) (*models.Post, error)
	GetBySlug(slug string) (*models.Post, error)
	GetByNumber(number int) (*models.Post, error)
	Search(query, view, limit string, tags []string) ([]*models.Post, error)
	GetAll() ([]*models.Post, error)
	Add(title, description string) (*models.Post, error)
	Update(post *models.Post, title, description string) (*models.Post, error)
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
	Block(userID int) error
	Unblock(userID int) error
	Count() (int, error)
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
