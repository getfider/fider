package inmemory

import (
	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
)

var demo = &models.Tenant{ID: 1, Name: "Demonstration"}
var orange = &models.Tenant{ID: 2, Name: "The Orange Inc."}

// TenantStorage contains read and write operations for tenants
type TenantStorage struct{}

// GetByDomain returns a tenant based on its domain
func (t *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	if extractSubdomain(domain) == "orange" {
		return orange, nil
	}
	if extractSubdomain(domain) == "demo" {
		return demo, nil
	}
	return nil, app.ErrNotFound
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
