package postgres_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
)

var db *dbx.Database
var trx *dbx.Trx

var tenants *postgres.TenantStorage
var users *postgres.UserStorage
var posts *postgres.PostStorage
var tags *postgres.TagStorage
var notifications *postgres.NotificationStorage
var events *postgres.EventStorage

var demoTenant *models.Tenant
var avengersTenant *models.Tenant
var gotTenant *models.Tenant
var jonSnow *models.User
var aryaStark *models.User
var sansaStark *models.User
var tonyStark *models.User

func SetupDatabaseTest(t *testing.T) {
	RegisterT(t)
	trx, _ = db.Begin()

	tenants = postgres.NewTenantStorage(trx)
	users = postgres.NewUserStorage(trx)
	posts = postgres.NewPostStorage(trx)
	tags = postgres.NewTagStorage(trx)
	notifications = postgres.NewNotificationStorage(trx)
	events = postgres.NewEventStorage(trx)

	demoTenant, _ = tenants.GetByDomain("demo")
	avengersTenant, _ = tenants.GetByDomain("avengers")

	users.SetCurrentTenant(demoTenant)
	jonSnow, _ = users.GetByID(1)
	aryaStark, _ = users.GetByID(2)
	sansaStark, _ = users.GetByID(3)

	users.SetCurrentTenant(avengersTenant)
	tonyStark, _ = users.GetByID(4)
}

func TeardownDatabaseTest() {
	trx.Rollback()
}

func TestMain(m *testing.M) {
	db = dbx.New()
	db.Seed()
	defer db.Close()

	code := m.Run()
	os.Exit(code)
}
