package dbx_test

import (
	"context"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
)

func TestTryLock_MultipleProcesses_SameKey(t *testing.T) {
	RegisterT(t)
	ctx := context.Background()
	trx1, _ := dbx.BeginTx(context.Background())
	locked1, unlock1 := dbx.TryLock(ctx, trx1, "KEY_1")
	defer trx1.MustRollback()

	trx2, _ := dbx.BeginTx(context.Background())
	locked2, unlock2 := dbx.TryLock(ctx, trx2, "KEY_1")
	defer trx2.MustRollback()

	Expect(locked1).IsTrue()
	Expect(unlock1).IsNotNil()
	Expect(locked2).IsFalse()
	Expect(unlock2).IsNotNil()

	unlock1()

	locked2, unlock2 = dbx.TryLock(ctx, trx2, "KEY_1")
	Expect(locked2).IsTrue()
	Expect(unlock2).IsNotNil()

	unlock2()
}

func TestTryLock_MultipleProcesses_DifferentKey(t *testing.T) {
	RegisterT(t)
	ctx := context.Background()

	trx1, _ := dbx.BeginTx(context.Background())
	locked1, unlock1 := dbx.TryLock(ctx, trx1, "KEY_1")
	defer trx1.MustRollback()

	trx2, _ := dbx.BeginTx(context.Background())
	locked2, unlock2 := dbx.TryLock(ctx, trx2, "KEY_2")
	defer trx2.MustRollback()

	Expect(locked1).IsTrue()
	Expect(unlock1).IsNotNil()
	Expect(locked2).IsTrue()
	Expect(unlock2).IsNotNil()

	unlock1()
	unlock2()
}
