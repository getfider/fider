package inmemory

import (
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	lastID   int
	tenants  []*models.Tenant
	current  *models.Tenant
	requests []*models.SignInRequest
}

// SetCurrentTenant tenant
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
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
		if tenant.Subdomain == extractSubdomain(domain) {
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

// UpdateSettings of given tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	for _, tenant := range s.tenants {
		if tenant.ID == s.current.ID {
			tenant.Invitation = settings.Invitation
			tenant.WelcomeMessage = settings.WelcomeMessage
			tenant.Name = settings.Title
			return nil
		}
	}
	return nil
}

// SaveVerificationKey used by e-mail verification
func (s *TenantStorage) SaveVerificationKey(email, key string) error {
	s.requests = append(s.requests, &models.SignInRequest{
		Email:      email,
		Key:        key,
		CreatedOn:  time.Now(),
		VerifiedOn: nil,
	})
	return nil
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(key string) (*models.SignInRequest, error) {
	for _, request := range s.requests {
		if request.Key == key {
			return request, nil
		}
	}
	return nil, app.ErrNotFound
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	for _, request := range s.requests {
		if request.Key == key {
			now := time.Now()
			request.VerifiedOn = &now
		}
	}
	return nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
