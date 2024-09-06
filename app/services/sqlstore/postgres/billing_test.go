package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestLockExpiredTenants_ShouldTriggerForOneTenant(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// There is a tenant with an expired trial setup in the seed for the test database.
	q := &cmd.LockExpiredTenants{}

	err := bus.Dispatch(ctx, q)
	Expect(err).IsNil()
	Expect(q.NumOfTenantsLocked).Equals(int64(1))
	Expect(q.TenantsLocked).Equals([]int{3})

}
