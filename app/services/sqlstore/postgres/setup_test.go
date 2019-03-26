package postgres_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/sqlstore/postgres"

	"github.com/getfider/fider/app"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/web"
)

var trx *dbx.Trx
var ctx context.Context

func SetupDatabaseTest(t *testing.T) {
	RegisterT(t)

	u, _ := url.Parse("http://cdn.test.fider.io")
	req := web.Request{URL: u}
	ctx = context.WithValue(context.Background(), app.RequestCtxKey, req)

	trx, _ = dbx.BeginTx(ctx)
	ctx = context.WithValue(ctx, app.TransactionCtxKey, trx)

	bus.Init(postgres.Service{})
}

func TeardownDatabaseTest() {
	trx.Rollback()
}

func TestMain(m *testing.M) {
	dbx.Seed()

	code := m.Run()
	os.Exit(code)
}
