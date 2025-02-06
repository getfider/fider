package postgres_test

import (
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"

	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
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
