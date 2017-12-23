package postgres_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

var db *dbx.Database
var trx *dbx.Trx

var tenants *postgres.TenantStorage
var users *postgres.UserStorage
var ideas *postgres.IdeaStorage
var tags *postgres.TagStorage

var orangeTenant *models.Tenant
var demoTenant *models.Tenant
var jonSnow *models.User
var aryaStark *models.User

func SetupDatabaseTest(t *testing.T) {
	RegisterTestingT(t)
	trx, _ = db.Begin()

	tenants = postgres.NewTenantStorage(trx)
	users = postgres.NewUserStorage(trx)
	ideas = postgres.NewIdeaStorage(trx)
	tags = postgres.NewTagStorage(trx)

	orangeTenant, _ = tenants.GetByDomain("orange")
	demoTenant, _ = tenants.GetByDomain("demo")
	jonSnow, _ = users.GetByID(1)
	aryaStark, _ = users.GetByID(2)
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
