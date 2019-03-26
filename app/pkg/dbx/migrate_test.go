package dbx_test

import (
	"context"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
)

func setupMigrationTest(t *testing.T) {
	RegisterT(t)
	ctx := context.Background()

	trx, _ := dbx.BeginTx(ctx)
	trx.Execute("DELETE FROM migrations_history WHERE version >= 210001010000")
	trx.Execute("DROP TABLE IF EXISTS dummy")
	trx.Execute("DROP TABLE IF EXISTS foo")
	trx.Commit()
}

func TestMigrate_Success(t *testing.T) {
	setupMigrationTest(t)
	ctx := context.Background()

	err := dbx.Migrate(ctx, "/app/pkg/dbx/testdata/migration_success")
	Expect(err).IsNil()

	trx, _ := dbx.BeginTx(ctx)
	var value string
	err = trx.Scalar(&value, "SELECT description FROM dummy WHERE id = 200 LIMIT 1")
	Expect(err).IsNil()
	Expect(value).Equals("Description 200Y")

	var count int
	err = trx.Scalar(&count, "SELECT COUNT(*) FROM dummy")
	Expect(err).IsNil()
	Expect(count).Equals(2)
	trx.Rollback()
}

func TestMigrate_Failure(t *testing.T) {
	setupMigrationTest(t)
	ctx := context.Background()

	trx, _ := dbx.BeginTx(ctx)
	defer trx.Rollback()

	err := dbx.Migrate(context.Background(), "/app/pkg/dbx/testdata/migration_failure")
	Expect(err).IsNotNil()

	_, err = trx.Execute("SELECT description FROM dummy")
	Expect(err).IsNotNil()
}
