package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTagStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))
	tag, err := tags.Add("Feature Request", "FF0000", true)
	Expect(err).To(BeNil())
	Expect(tag.ID).To(Equal(1))

	dbTag, err := tags.GetBySlug("feature-request")

	Expect(err).To(BeNil())
	Expect(dbTag.ID).To(Equal(1))
	Expect(dbTag.Name).To(Equal("Feature Request"))
	Expect(dbTag.Slug).To(Equal("feature-request"))
	Expect(dbTag.Color).To(Equal("FF0000"))
	Expect(dbTag.IsPublic).To(BeTrue())
}

func TestTagStorage_AddUpdateAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))
	tag, err := tags.Add("Feature Request", "FF0000", true)
	tag, err = tags.Update(tag.ID, "Bug", "000000", false)

	dbTag, err := tags.GetBySlug("bug")

	Expect(err).To(BeNil())
	Expect(dbTag.ID).To(Equal(tag.ID))
	Expect(dbTag.Name).To(Equal("Bug"))
	Expect(dbTag.Slug).To(Equal("bug"))
	Expect(dbTag.Color).To(Equal("000000"))
	Expect(dbTag.IsPublic).To(BeFalse())
}

func TestTagStorage_AddDeleteAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))
	tag, err := tags.Add("Bug", "FFFFFF", true)

	err = tags.Delete(tag.ID)
	Expect(err).To(BeNil())

	dbTag, err := tags.GetBySlug("bug")

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbTag).To(BeNil())
}

func TestTagStorage_Assign_Unassign(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)

	ideas := postgres.NewIdeaStorage(trx)
	ideas.SetCurrentTenant(demoTenant(tenants))

	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))

	idea, _ := ideas.Add("My great idea", "with a great description", 2)
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag.ID, idea.ID, 2)
	Expect(err).To(BeNil())

	assigned, err := tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(1))
	Expect(assigned[0].ID).To(Equal(tag.ID))
	Expect(assigned[0].Name).To(Equal("Bug"))
	Expect(assigned[0].Slug).To(Equal("bug"))
	Expect(assigned[0].Color).To(Equal("FFFFFF"))
	Expect(assigned[0].IsPublic).To(BeTrue())

	err = tags.UnassignTag(tag.ID, idea.ID)
	Expect(err).To(BeNil())

	assigned, err = tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(0))
}

func TestTagStorage_Assign_DeleteTag(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)

	ideas := postgres.NewIdeaStorage(trx)
	ideas.SetCurrentTenant(demoTenant(tenants))

	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))

	idea, _ := ideas.Add("My great idea", "with a great description", 2)
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag.ID, idea.ID, 2)
	Expect(err).To(BeNil())

	err = tags.Delete(tag.ID)
	Expect(err).To(BeNil())

	assigned, err := tags.GetAssigned(idea.ID)
	Expect(err).To(BeNil())
	Expect(len(assigned)).To(Equal(0))
}

func TestTagStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenants := postgres.NewTenantStorage(trx)
	users := postgres.NewUserStorage(trx)
	users.SetCurrentTenant(demoTenant(tenants))

	tags := postgres.NewTagStorage(trx)
	tags.SetCurrentTenant(demoTenant(tenants))
	tags.SetCurrentUser(jonSnow(users))
	tags.Add("Feature Request", "FF0000", true)
	tags.Add("Bug", "0F0F0F", false)

	allTags, err := tags.GetAll()

	Expect(err).To(BeNil())
	Expect(len(allTags)).To(Equal(2))

	Expect(allTags[0].ID).NotTo(BeZero())
	Expect(allTags[0].Name).To(Equal("Feature Request"))
	Expect(allTags[0].Slug).To(Equal("feature-request"))
	Expect(allTags[0].Color).To(Equal("FF0000"))
	Expect(allTags[0].IsPublic).To(BeTrue())

	Expect(allTags[1].ID).NotTo(BeZero())
	Expect(allTags[1].Name).To(Equal("Bug"))
	Expect(allTags[1].Slug).To(Equal("bug"))
	Expect(allTags[1].Color).To(Equal("0F0F0F"))
	Expect(allTags[1].IsPublic).To(BeFalse())

	tags.SetCurrentUser(aryaStark(users))

	visitorTags, err := tags.GetAll()
	Expect(err).To(BeNil())
	Expect(len(visitorTags)).To(Equal(1))
	Expect(visitorTags[0].Name).To(Equal("Feature Request"))
}
