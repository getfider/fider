package postgres

import (
	"strings"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
)

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	DB *dbx.Database
}

// GetByDomain returns a tenant based on its domain
func (svc *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	row := svc.DB.QueryRow("SELECT id, name, subdomain FROM tenants WHERE subdomain = $1 OR cname = $2", extractSubdomain(domain), domain)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Subdomain)
	if err != nil {
		return nil, app.ErrNotFound
	}

	return tenant, nil
}

func extractSubdomain(domain string) string {
	if idx := strings.Index(domain, "."); idx != -1 {
		return domain[:idx]
	}
	return domain
}
