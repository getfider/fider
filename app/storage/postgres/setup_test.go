package postgres_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/web"
	postgres2 "github.com/getfider/fider/app/services/sqlstore/postgres"
	"github.com/getfider/fider/app/storage/postgres"
)

var trx *dbx.Trx

var tenants *postgres.TenantStorage
var users *postgres.UserStorage
var posts *postgres.PostStorage
var tags *postgres.TagStorage

var demoTenant *models.Tenant
var avengersTenant *models.Tenant
var gotTenant *models.Tenant
var jonSnow *models.User
var aryaStark *models.User
var sansaStark *models.User
var tonyStark *models.User

var demoTenantCtx context.Context
var avengersTenantCtx context.Context
var gotTenantCtx context.Context

var jonSnowCtx context.Context
var aryaStarkCtx context.Context
var sansaStarkCtx context.Context
var tonyStarkCtx context.Context

func SetupDatabaseTest(t *testing.T) context.Context {
	RegisterT(t)
	bus.Init(postgres2.Service{})

	u, _ := url.Parse("http://cdn.test.fider.io")
	req := web.Request{URL: u}
	ctx := context.WithValue(context.Background(), app.RequestCtxKey, req)

	trx, _ = dbx.BeginTx(ctx)
	tenants = postgres.NewTenantStorage(trx, ctx)
	users = postgres.NewUserStorage(trx, ctx)
	posts = postgres.NewPostStorage(trx, ctx)
	tags = postgres.NewTagStorage(trx, ctx)

	demoTenant, _ = tenants.GetByDomain("demo")
	avengersTenant, _ = tenants.GetByDomain("avengers")

	users.SetCurrentTenant(demoTenant)
	jonSnow, _ = users.GetByID(1)
	aryaStark, _ = users.GetByID(2)
	sansaStark, _ = users.GetByID(3)

	users.SetCurrentTenant(avengersTenant)
	tonyStark, _ = users.GetByID(4)

	trxCtx := context.WithValue(ctx, app.TransactionCtxKey, trx)

	demoTenantCtx = withTenant(trxCtx, demoTenant)
	avengersTenantCtx = withTenant(trxCtx, avengersTenant)
	gotTenantCtx = withTenant(trxCtx, gotTenant)

	jonSnowCtx = withUser(trxCtx, jonSnow)
	aryaStarkCtx = withUser(trxCtx, aryaStark)
	sansaStarkCtx = withUser(trxCtx, sansaStark)
	tonyStarkCtx = withUser(trxCtx, tonyStark)
	return trxCtx
}

func TeardownDatabaseTest() {
	trx.Rollback()
}

func TestMain(m *testing.M) {
	dbx.Seed()

	code := m.Run()
	os.Exit(code)
}
