package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestTagStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tags.SetCurrentTenant(demoTenant)
	tag, err := tags.Add("Feature Request", "FF0000", true)
	Expect(err).IsNil()
	Expect(tag.ID).NotEquals(0)

	dbTag, err := tags.GetBySlug("feature-request")

	Expect(err).IsNil()
	Expect(dbTag.ID).NotEquals(0)
	Expect(dbTag.Name).Equals("Feature Request")
	Expect(dbTag.Slug).Equals("feature-request")
	Expect(dbTag.Color).Equals("FF0000")
	Expect(dbTag.IsPublic).IsTrue()
}

func TestTagStorage_AddUpdateAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tags.SetCurrentTenant(demoTenant)
	tag, err := tags.Add("Feature Request", "FF0000", true)
	tag, err = tags.Update(tag, "Bug", "000000", false)

	dbTag, err := tags.GetBySlug("bug")

	Expect(err).IsNil()
	Expect(dbTag.ID).Equals(tag.ID)
	Expect(dbTag.Name).Equals("Bug")
	Expect(dbTag.Slug).Equals("bug")
	Expect(dbTag.Color).Equals("000000")
	Expect(dbTag.IsPublic).IsFalse()
}

func TestTagStorage_AddDeleteAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tags.SetCurrentTenant(demoTenant)
	tag, err := tags.Add("Bug", "FFFFFF", true)

	err = tags.Delete(tag)
	Expect(err).IsNil()

	dbTag, err := tags.GetBySlug("bug")

	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(dbTag).IsNil()
}

func TestTagStorage_Assign_Unassign(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)
	tags.SetCurrentTenant(demoTenant)
	tags.SetCurrentUser(aryaStark)

	post, _ := posts.Add("My great post", "with a great description")
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag, post)
	Expect(err).IsNil()

	assigned, err := tags.GetAssigned(post)
	Expect(err).IsNil()
	Expect(assigned).HasLen(1)
	Expect(assigned[0].ID).Equals(tag.ID)
	Expect(assigned[0].Name).Equals("Bug")
	Expect(assigned[0].Slug).Equals("bug")
	Expect(assigned[0].Color).Equals("FFFFFF")
	Expect(assigned[0].IsPublic).IsTrue()

	err = tags.UnassignTag(tag, post)
	Expect(err).IsNil()

	assigned, err = tags.GetAssigned(post)
	Expect(err).IsNil()
	Expect(assigned).HasLen(0)
}

func TestTagStorage_Assign_DeleteTag(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)
	tags.SetCurrentTenant(demoTenant)
	tags.SetCurrentUser(aryaStark)

	post, _ := posts.Add("My great post", "with a great description")
	tag, _ := tags.Add("Bug", "FFFFFF", true)

	err := tags.AssignTag(tag, post)
	Expect(err).IsNil()

	err = tags.Delete(tag)
	Expect(err).IsNil()

	assigned, err := tags.GetAssigned(post)
	Expect(err).IsNil()
	Expect(assigned).HasLen(0)
}

func TestTagStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	tags.SetCurrentTenant(demoTenant)
	tags.SetCurrentUser(jonSnow)

	tags.Add("Feature Request", "FF0000", true)
	tags.Add("Bug", "0F0F0F", false)

	allTags, err := tags.GetAll()

	Expect(err).IsNil()
	Expect(allTags).HasLen(2)

	Expect(allTags[0].ID).NotEquals(0)
	Expect(allTags[0].Name).Equals("Bug")
	Expect(allTags[0].Slug).Equals("bug")
	Expect(allTags[0].Color).Equals("0F0F0F")
	Expect(allTags[0].IsPublic).IsFalse()

	Expect(allTags[1].ID).NotEquals(0)
	Expect(allTags[1].Name).Equals("Feature Request")
	Expect(allTags[1].Slug).Equals("feature-request")
	Expect(allTags[1].Color).Equals("FF0000")
	Expect(allTags[1].IsPublic).IsTrue()

	tags.SetCurrentUser(aryaStark)

	visitorTags, err := tags.GetAll()
	Expect(err).IsNil()
	Expect(visitorTags).HasLen(1)
	Expect(visitorTags[0].Name).Equals("Feature Request")
}
