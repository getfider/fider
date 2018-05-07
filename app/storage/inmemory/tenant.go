package inmemory

import (
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	lastID        int
	lastLogoID    int
	tenants       []*models.Tenant
	current       *models.Tenant
	user          *models.User
	verifications []*models.EmailVerification
	tenantLogos   map[int]*models.Upload
}

// SetCurrentTenant tenant
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
}

// SetCurrentUser to current context
func (s *TenantStorage) SetCurrentUser(user *models.User) {
	s.user = user
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
		if tenant.Subdomain == extractSubdomain(domain) || tenant.CNAME == domain {
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
	for _, tenant := range s.tenants {
		if tenant.ID == s.current.ID {

			if len(settings.Logo) > 0 {
				s.lastLogoID = s.lastLogoID + 1
				if s.tenantLogos == nil {
					s.tenantLogos = make(map[int]*models.Upload, 0)
				}
				tenant.LogoID = s.lastLogoID
				s.tenantLogos[s.lastLogoID] = &models.Upload{
					Content:     settings.Logo,
					Size:        len(settings.Logo),
					ContentType: http.DetectContentType(settings.Logo),
				}
			}

			tenant.Invitation = settings.Invitation
			tenant.WelcomeMessage = settings.WelcomeMessage
			tenant.Name = settings.Title
			tenant.CNAME = settings.CNAME
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
		CreatedOn:  time.Now(),
		ExpiresOn:  time.Now().Add(duration),
		VerifiedOn: nil,
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
			verification.VerifiedOn = &now
		}
	}
	return nil
}

// GetLogo returns tenant logo by id
func (s *TenantStorage) GetLogo(id int) (*models.Upload, error) {
	if s.tenantLogos != nil {
		logo, ok := s.tenantLogos[id]
		if !ok {
			return nil, app.ErrNotFound
		}
		return logo, nil
	}
	return nil, app.ErrNotFound
}

func extractSubdomain(hostname string) string {
	domain := env.MultiTenantDomain()
	if domain == "" {
		return hostname
	}

	return strings.Replace(hostname, domain, "", -1)
}
