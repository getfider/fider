package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestStatusStorage_ListActive_SeededRows(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	q := &query.ListActiveStatusesForTenant{}
	err := bus.Dispatch(demoTenantCtx, q)
	Expect(err).IsNil()
	Expect(len(q.Result) >= 6).IsTrue()

	bySlug := map[string]bool{}
	for _, s := range q.Result {
		bySlug[s.Slug] = true
	}
	Expect(bySlug["open"]).IsTrue()
	Expect(bySlug["planned"]).IsTrue()
	Expect(bySlug["started"]).IsTrue()
	Expect(bySlug["completed"]).IsTrue()
	Expect(bySlug["declined"]).IsTrue()
	Expect(bySlug["duplicate"]).IsTrue()
}

func TestStatusStorage_CreateUpdateDelete_CustomStatus(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	create := &cmd.CreateStatus{
		Slug:       "triage",
		Label:      "Triage",
		Kind:       "open",
		Color:      "blue",
		Icon:       "lightbulb",
		ShowOnHome: true,
		Filterable: true,
		SortOrder:  15,
	}
	err := bus.Dispatch(demoTenantCtx, create)
	Expect(err).IsNil()
	Expect(create.Result).IsNotNil()
	Expect(create.Result.Slug).Equals("triage")
	Expect(create.Result.IsSystem).IsFalse()

	update := &cmd.UpdateStatus{
		ID:         create.Result.ID,
		Label:      "Under Review",
		Color:      "yellow",
		Icon:       "lightbulb",
		ShowOnHome: true,
		Filterable: true,
		SortOrder:  18,
		IsActive:   true,
	}
	err = bus.Dispatch(demoTenantCtx, update)
	Expect(err).IsNil()

	getBySlug := &query.GetStatusBySlug{Slug: "triage"}
	err = bus.Dispatch(demoTenantCtx, getBySlug)
	Expect(err).IsNil()
	Expect(getBySlug.Result.Label).Equals("Under Review")
	Expect(getBySlug.Result.Color).Equals("yellow")
	Expect(getBySlug.Result.SortOrder).Equals(18)

	del := &cmd.DeleteStatus{ID: create.Result.ID}
	err = bus.Dispatch(demoTenantCtx, del)
	Expect(err).IsNil()

	missing := &query.GetStatusBySlug{Slug: "triage"}
	err = bus.Dispatch(demoTenantCtx, missing)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
}

func TestStatusStorage_Delete_RefusesSystemRow(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	getBySlug := &query.GetStatusBySlug{Slug: "completed"}
	err := bus.Dispatch(demoTenantCtx, getBySlug)
	Expect(err).IsNil()
	Expect(getBySlug.Result.IsSystem).IsTrue()

	del := &cmd.DeleteStatus{ID: getBySlug.Result.ID}
	err = bus.Dispatch(demoTenantCtx, del)
	Expect(err).IsNotNil()
}

func TestStatusStorage_Create_RejectsUnknownKind(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	create := &cmd.CreateStatus{
		Slug:  "bogus",
		Label: "Bogus",
		Kind:  "garbage",
		Color: "blue",
		Icon:  "lightbulb",
	}
	err := bus.Dispatch(demoTenantCtx, create)
	Expect(err).IsNotNil()
}

func TestStatusStorage_Create_UniquePerTenant(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	first := &cmd.CreateStatus{Slug: "triage", Label: "Triage", Kind: "open", Color: "blue", Icon: "lightbulb"}
	err := bus.Dispatch(demoTenantCtx, first)
	Expect(err).IsNil()

	dup := &cmd.CreateStatus{Slug: "triage", Label: "Triage 2", Kind: "open", Color: "blue", Icon: "lightbulb"}
	err = bus.Dispatch(demoTenantCtx, dup)
	Expect(err).IsNotNil()
}

func TestStatusStorage_SeedTenantStatuses_Idempotent(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	seed := &cmd.SeedTenantStatuses{TenantID: 1}
	err := bus.Dispatch(demoTenantCtx, seed)
	Expect(err).IsNil()

	// Second run is a no-op thanks to ON CONFLICT DO NOTHING.
	err = bus.Dispatch(demoTenantCtx, seed)
	Expect(err).IsNil()

	q := &query.ListActiveStatusesForTenant{}
	err = bus.Dispatch(demoTenantCtx, q)
	Expect(err).IsNil()
	Expect(len(q.Result) >= 6).IsTrue()
}
