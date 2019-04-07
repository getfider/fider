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
	SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error
	FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error)
	SetKeyAsVerified(key string) error
	SaveOAuthConfig(config *models.CreateEditOAuthConfig) error
	GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error)
	ListOAuthConfig() ([]*models.OAuthConfig, error)
}
