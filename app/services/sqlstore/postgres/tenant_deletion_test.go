package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"

	. "github.com/getfider/fider/app/pkg/assert"
)

func countTenantRows(table string, tenantID int) int {
	var count int
	err := trx.Scalar(&count, "SELECT COUNT(*) FROM "+table+" WHERE tenant_id = $1", tenantID)
	Expect(err).IsNil()
	return count
}

// seedTenantContent creates a post, comment, reaction, vote and tag so deletion exercises
// the full FK-ordered teardown (including the tenant-less reactions table).
func seedTenantContent(tenantCtx, userCtx context.Context) {
	post := &cmd.AddNewPost{Title: "A feature request to delete", Description: "please"}
	Expect(bus.Dispatch(userCtx, post)).IsNil()

	comment := &cmd.AddNewComment{Post: post.Result, Content: "great idea"}
	Expect(bus.Dispatch(userCtx, comment)).IsNil()

	user, _ := userCtx.Value(app.UserCtxKey).(*entity.User)
	Expect(bus.Dispatch(userCtx, &cmd.ToggleCommentReaction{
		Comment: &entity.Comment{ID: comment.Result.ID}, Emoji: "👍", User: user,
	})).IsNil()
	Expect(bus.Dispatch(userCtx, &cmd.AddVote{Post: post.Result, User: user})).IsNil()

	tag := &cmd.AddNewTag{Name: "doomed", Color: "ffffff", IsPublic: true}
	Expect(bus.Dispatch(tenantCtx, tag)).IsNil()
	Expect(bus.Dispatch(userCtx, &cmd.AssignTag{Tag: tag.Result, Post: post.Result})).IsNil()
}

func TestTenantDeletion_DeletesEverythingForTenantOnly(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	seedTenantContent(demoTenantCtx, jonSnowCtx)
	seedTenantContent(avengersTenantCtx, tonyStarkCtx)

	Expect(countTenantRows("users", demoTenant.ID) > 0).IsTrue()
	Expect(countTenantRows("posts", demoTenant.ID) > 0).IsTrue()
	Expect(countTenantRows("comments", demoTenant.ID) > 0).IsTrue()

	avengersUsersBefore := countTenantRows("users", avengersTenant.ID)
	avengersPostsBefore := countTenantRows("posts", avengersTenant.ID)
	avengersCommentsBefore := countTenantRows("comments", avengersTenant.ID)
	Expect(avengersPostsBefore > 0).IsTrue()

	err := bus.Dispatch(ctx, &cmd.DeleteTenant{TenantID: demoTenant.ID})
	Expect(err).IsNil()

	// Every tenant-scoped table is empty for the deleted tenant.
	for _, table := range []string{
		"users", "posts", "comments", "post_votes", "post_tags", "tags",
		"post_subscribers", "notifications", "attachments", "user_settings",
		"user_providers", "email_verifications", "events", "blobs", "webhooks",
		"oauth_providers", "tenant_providers", "tenants_billing",
	} {
		Expect(countTenantRows(table, demoTenant.ID)).Equals(0)
	}

	// The tenant row itself is gone.
	getDemo := &query.GetTenantByDomain{Domain: "demo"}
	err = bus.Dispatch(ctx, getDemo)
	Expect(errors.Cause(err)).IsNotNil()

	// Another tenant is left completely intact.
	Expect(countTenantRows("users", avengersTenant.ID)).Equals(avengersUsersBefore)
	Expect(countTenantRows("posts", avengersTenant.ID)).Equals(avengersPostsBefore)
	Expect(countTenantRows("comments", avengersTenant.ID)).Equals(avengersCommentsBefore)
}

func TestTenantDeletion_ScheduleAndCancel(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	scheduledAt := time.Now().Add(time.Hour)
	err := bus.Dispatch(ctx, &cmd.ScheduleTenantDeletion{
		TenantID:          demoTenant.ID,
		RequestedByUserID: jonSnow.ID,
		CancelKey:         "test-cancel-key-123",
		ScheduledAt:       scheduledAt,
	})
	Expect(err).IsNil()

	// Status is untouched during the grace window.
	getDemo := &query.GetTenantByDomain{Domain: "demo"}
	Expect(bus.Dispatch(ctx, getDemo)).IsNil()
	Expect(getDemo.Result.Status).Equals(enum.TenantActive)
	Expect(getDemo.Result.ScheduledDeletionAt).IsNotNil()

	// The cancel key resolves back to the tenant.
	byKey := &query.GetTenantByCancelKey{Key: "test-cancel-key-123"}
	Expect(bus.Dispatch(ctx, byKey)).IsNil()
	Expect(byKey.Result.ID).Equals(demoTenant.ID)

	// Cancelling clears the schedule.
	Expect(bus.Dispatch(ctx, &cmd.CancelTenantDeletion{TenantID: demoTenant.ID})).IsNil()
	getDemo = &query.GetTenantByDomain{Domain: "demo"}
	Expect(bus.Dispatch(ctx, getDemo)).IsNil()
	Expect(getDemo.Result.ScheduledDeletionAt).IsNil()
}

func TestTenantDeletion_PendingDeletionAndOwner(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Owner is the lowest-id active administrator (the seed creator, Jon Snow).
	owner := &query.GetTenantOwner{TenantID: demoTenant.ID}
	Expect(bus.Dispatch(ctx, owner)).IsNil()
	Expect(owner.Result.ID).Equals(jonSnow.ID)
	Expect(owner.Result.Role).Equals(enum.RoleAdministrator)

	// A schedule in the past is returned; one in the future is not.
	Expect(bus.Dispatch(ctx, &cmd.ScheduleTenantDeletion{
		TenantID: demoTenant.ID, RequestedByUserID: jonSnow.ID,
		CancelKey: "past-key", ScheduledAt: time.Now().Add(-time.Minute),
	})).IsNil()
	Expect(bus.Dispatch(ctx, &cmd.ScheduleTenantDeletion{
		TenantID: avengersTenant.ID, RequestedByUserID: tonyStark.ID,
		CancelKey: "future-key", ScheduledAt: time.Now().Add(time.Hour),
	})).IsNil()

	pending := &query.GetTenantsPendingDeletion{}
	Expect(bus.Dispatch(ctx, pending)).IsNil()

	ids := map[int]bool{}
	for _, tn := range pending.Result {
		ids[tn.ID] = true
	}
	Expect(ids[demoTenant.ID]).IsTrue()
	Expect(ids[avengersTenant.ID]).IsFalse()
}
