package inmemory

import (
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	lastID        int
	tenants       []*models.Tenant
	current       *models.Tenant
	user          *models.User
	verifications []*models.EmailVerification
	oauthConfigs  []*models.OAuthConfig
}

//NewTenantStorage creates a new TenantStorage
func NewTenantStorage() *TenantStorage {
	return &TenantStorage{
		tenants:       make([]*models.Tenant, 0),
		verifications: make([]*models.EmailVerification, 0),
		oauthConfigs:  make([]*models.OAuthConfig, 0),
	}
}

// SetCurrentTenant tenant
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
}

// SetCurrentUser to current context
func (s *TenantStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// Current returns current context tenant
func (s *TenantStorage) Current() *models.Tenant {
	return s.current
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string, status int) (*models.Tenant, error) {
	s.lastID = s.lastID + 1
	tenant := &models.Tenant{ID: s.lastID, Name: name, Subdomain: subdomain, Status: status}
	s.tenants = append(s.tenants, tenant)
	return tenant, nil
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	for _, tenant := range s.tenants {
		return tenant, nil
	}

	return nil, app.ErrNotFound
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	for _, tenant := range s.tenants {
		if tenant.Subdomain == env.Subdomain(domain) || tenant.Subdomain == domain || tenant.CNAME == domain {
			return tenant, nil
		}
	}

	return nil, app.ErrNotFound
}

// IsSubdomainAvailable returns true if subdomain is available to use
func (s *TenantStorage) IsSubdomainAvailable(subdomain string) (bool, error) {
	for _, tenant := range s.tenants {
		if tenant.Subdomain == subdomain {
			return false, nil
		}
	}
	return true, nil
}

// IsCNAMEAvailable returns true if cname is available to use
func (s *TenantStorage) IsCNAMEAvailable(cname string) (bool, error) {
	for _, tenant := range s.tenants {
		if tenant.CNAME == cname {
			return false, nil
		}
	}
	return true, nil
}

// UpdateSettings of current tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	if settings.Logo.Remove {
		settings.Logo.BlobKey = ""
	}

	for _, tenant := range s.tenants {
		if tenant.ID == s.current.ID {
			tenant.Invitation = settings.Invitation
			tenant.WelcomeMessage = settings.WelcomeMessage
			tenant.Name = settings.Title
			tenant.CNAME = settings.CNAME
			tenant.LogoBlobKey = settings.Logo.BlobKey
			return nil
		}
	}
	return nil
}

// UpdateBillingSettings of current tenant
func (s *TenantStorage) UpdateBillingSettings(billing *models.TenantBilling) error {
	s.current.Billing = billing
	return nil
}

// UpdateAdvancedSettings of current tenant
func (s *TenantStorage) UpdateAdvancedSettings(settings *models.UpdateTenantAdvancedSettings) error {
	for _, tenant := range s.tenants {
		if tenant.ID == s.current.ID {
			tenant.CustomCSS = settings.CustomCSS
			return nil
		}
	}
	return nil
}

// UpdatePrivacy settings of current tenant
func (s *TenantStorage) UpdatePrivacy(settings *models.UpdateTenantPrivacy) error {
	for _, tenant := range s.tenants {
		if tenant.ID == s.current.ID {
			tenant.IsPrivate = settings.IsPrivate
			return nil
		}
	}
	return nil
}

// Activate given tenant
func (s *TenantStorage) Activate(id int) error {
	for _, tenant := range s.tenants {
		if tenant.ID == id {
			tenant.Status = models.TenantActive
			return nil
		}
	}
	return app.ErrNotFound
}

// SaveVerificationKey used by email verification
func (s *TenantStorage) SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error {
	userID := 0
	if request.GetUser() != nil {
		userID = request.GetUser().ID
	}
	s.verifications = append(s.verifications, &models.EmailVerification{
		Email:      request.GetEmail(),
		Name:       request.GetName(),
		Kind:       request.GetKind(),
		Key:        key,
		UserID:     userID,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(duration),
		VerifiedAt: nil,
	})
	return nil
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error) {
	for _, verification := range s.verifications {
		if verification.Key == key && verification.Kind == kind {
			return verification, nil
		}
	}
	return nil, app.ErrNotFound
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	for _, verification := range s.verifications {
		if verification.Key == key {
			now := time.Now()
			verification.VerifiedAt = &now
		}
	}
	return nil
}

// SaveOAuthConfig saves given config into database
func (s *TenantStorage) SaveOAuthConfig(config *models.CreateEditOAuthConfig) error {
	bkey := ""
	if config.Logo != nil {
		bkey = config.Logo.BlobKey
	}
	for _, c := range s.oauthConfigs {
		if c.ID == config.ID {
			c.Provider = config.Provider
			c.DisplayName = config.DisplayName
			c.ClientID = config.ClientID
			c.ClientSecret = config.ClientSecret
			c.AuthorizeURL = config.AuthorizeURL
			c.TokenURL = config.TokenURL
			c.Scope = config.Scope
			c.ProfileURL = config.ProfileURL
			c.JSONUserIDPath = config.JSONUserIDPath
			c.JSONUserNamePath = config.JSONUserNamePath
			c.JSONUserEmailPath = config.JSONUserEmailPath
			c.LogoBlobKey = bkey
			return nil
		}
	}
	s.oauthConfigs = append(s.oauthConfigs, &models.OAuthConfig{
		ID:                config.ID,
		Provider:          config.Provider,
		DisplayName:       config.DisplayName,
		LogoBlobKey:       bkey,
		ClientID:          config.ClientID,
		ClientSecret:      config.ClientSecret,
		AuthorizeURL:      config.AuthorizeURL,
		TokenURL:          config.TokenURL,
		Scope:             config.Scope,
		ProfileURL:        config.ProfileURL,
		JSONUserIDPath:    config.JSONUserIDPath,
		JSONUserNamePath:  config.JSONUserNamePath,
		JSONUserEmailPath: config.JSONUserEmailPath,
	})
	return nil
}

// GetOAuthConfigByProvider returns a custom OAuth configuration by provider name
func (s *TenantStorage) GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error) {
	for _, c := range s.oauthConfigs {
		if c.Provider == provider {
			return c, nil
		}
	}
	return nil, app.ErrNotFound
}

// ListOAuthConfig returns a list of all custom OAuth provider for current tenant
func (s *TenantStorage) ListOAuthConfig() ([]*models.OAuthConfig, error) {
	return s.oauthConfigs, nil
}
