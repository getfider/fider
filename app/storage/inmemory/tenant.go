package inmemory

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	lastID  int
	tenants []*models.Tenant
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string) (*models.Tenant, error) {
	s.lastID = s.lastID + 1
	tenant := &models.Tenant{ID: s.lastID, Name: name, Subdomain: subdomain}
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
func (s *TenantStorage) UpdateSettings(tenantID int, title, invitation, welcomeMessage string) error {
	for _, tenant := range s.tenants {
		if tenant.ID == tenantID {
			tenant.Invitation = invitation
			tenant.WelcomeMessage = welcomeMessage
			tenant.Name = title
			return nil
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
