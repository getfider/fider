package postgres_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func orangeTenant(tenants *postgres.TenantStorage) *models.Tenant {
	tenant, _ := tenants.GetByDomain("orange")
	return tenant
}

func demoTenant(tenants *postgres.TenantStorage) *models.Tenant {
	tenant, _ := tenants.GetByDomain("demo")
	return tenant
}

var db *dbx.Database
var trx *dbx.Trx

func SetupDatabaseTest(t *testing.T) {
	RegisterTestingT(t)
	trx, _ = db.Begin()
}

func TeardownDatabaseTest() {
	trx.Rollback()
}

func TestMain(m *testing.M) {
	db, _ = dbx.New()
	db.Migrate()
	db.Seed()
	defer db.Close()

	code := m.Run()
	os.Exit(code)
}
