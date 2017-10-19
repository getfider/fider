package postgres_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
)

func orangeTenant(tenants *postgres.TenantStorage) *models.Tenant {
	tenant, _ := tenants.GetByDomain("orange")
	return tenant
}

func demoTenant(tenants *postgres.TenantStorage) *models.Tenant {
	tenant, _ := tenants.GetByDomain("demo")
	return tenant
}

func TestMain(m *testing.M) {
	db, _ := dbx.New()
	db.Migrate()
	db.Close()

	code := m.Run()
	os.Exit(code)
}
