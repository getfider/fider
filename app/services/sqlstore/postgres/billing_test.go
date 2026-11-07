package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestCancelStripeSubscription_OnlyCurrentSubscriptionDowngrades(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	err := bus.Dispatch(ctx, &cmd.ActivateStripeSubscription{TenantID: demoTenant.ID, CustomerID: "cus_1", SubscriptionID: "sub_current"})
	Expect(err).IsNil()

	isPro := func() bool {
		q := &query.GetTenantByDomain{Domain: "demo"}
		Expect(bus.Dispatch(ctx, q)).IsNil()
		return q.Result.IsPro
	}

	stale := &cmd.CancelStripeSubscription{TenantID: demoTenant.ID, SubscriptionID: "sub_other"}
	Expect(bus.Dispatch(ctx, stale)).IsNil()
	Expect(stale.Result).IsFalse()
	Expect(isPro()).IsTrue()

	current := &cmd.CancelStripeSubscription{TenantID: demoTenant.ID, SubscriptionID: "sub_current"}
	Expect(bus.Dispatch(ctx, current)).IsNil()
	Expect(current.Result).IsTrue()
	Expect(isPro()).IsFalse()

	state := &query.GetStripeBillingState{}
	Expect(bus.Dispatch(demoTenantCtx, state)).IsNil()
	Expect(state.Result.SubscriptionID).Equals("")
	Expect(state.Result.CustomerID).Equals("cus_1")
}
