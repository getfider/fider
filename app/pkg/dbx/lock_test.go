package dbx_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
)

func TestLock(t *testing.T) {
	RegisterT(t)
	db1 := dbx.New()
	defer db1.Close()
	db2 := dbx.New()
	defer db2.Close()

	locked, err := db1.TryLock()
	Expect(locked).IsTrue()
	Expect(err).IsNil()

	locked, err = db2.TryLock()
	Expect(locked).IsFalse()
	Expect(err).IsNil()

	Expect(db1.Unlock()).IsNil()

	locked, err = db2.TryLock()
	Expect(locked).IsTrue()
	Expect(err).IsNil()
}
