package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestIdeaStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	now := time.Now()

	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('add twitter integration', 'add-twitter-integration', 1, 'Would be great to see it integrated with twitter', $1, 1, 1, 0, 1)", now)
	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('this is my idea', 'this-is-my-idea', 2, 'no description', $1, 1, 2, 5, 2)", now)

	ideas.SetCurrentTenant(demoTenant)

	dbIdeas, err := ideas.GetAll()
	Expect(err).IsNil()
	Expect(dbIdeas).HasLen(2)

	Expect(dbIdeas[0].Title).Equals("this is my idea")
	Expect(dbIdeas[0].Slug).Equals("this-is-my-idea")
	Expect(dbIdeas[0].Number).Equals(2)
	Expect(dbIdeas[0].Description).Equals("no description")
	Expect(dbIdeas[0].User.Name).Equals("Arya Stark")
	Expect(dbIdeas[0].TotalSupporters).Equals(5)
	Expect(dbIdeas[0].Status).Equals(models.IdeaCompleted)

	Expect(dbIdeas[1].Title).Equals("add twitter integration")
	Expect(dbIdeas[1].Slug).Equals("add-twitter-integration")
	Expect(dbIdeas[1].Number).Equals(1)
	Expect(dbIdeas[1].Description).Equals("Would be great to see it integrated with twitter")
	Expect(dbIdeas[1].User.Name).Equals("Jon Snow")
	Expect(dbIdeas[1].TotalSupporters).Equals(0)
	Expect(dbIdeas[1].Status).Equals(models.IdeaStarted)

	dbIdeas, err = ideas.Search("twitter", "trending", []string{})
	Expect(err).IsNil()
	Expect(dbIdeas).HasLen(1)
	Expect(dbIdeas[0].Slug).Equals("add-twitter-integration")
}

func TestIdeaStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	dbIdeaById, err := ideas.GetByID(idea.ID)
	Expect(err).IsNil()
	Expect(dbIdeaById.ID).Equals(idea.ID)
	Expect(dbIdeaById.Number).Equals(1)
	Expect(dbIdeaById.ViewerSupported).IsFalse()
	Expect(dbIdeaById.TotalSupporters).Equals(0)
	Expect(dbIdeaById.Status).Equals(models.IdeaOpen)
	Expect(dbIdeaById.Title).Equals("My new idea")
	Expect(dbIdeaById.Description).Equals("with this description")
	Expect(dbIdeaById.User.ID).Equals(1)
	Expect(dbIdeaById.User.Name).Equals("Jon Snow")
	Expect(dbIdeaById.User.Email).Equals("jon.snow@got.com")

	dbIdeaBySlug, err := ideas.GetBySlug("my-new-idea")

	Expect(err).IsNil()
	Expect(dbIdeaBySlug.ID).Equals(idea.ID)
	Expect(dbIdeaBySlug.Number).Equals(1)
	Expect(dbIdeaBySlug.ViewerSupported).IsFalse()
	Expect(dbIdeaBySlug.TotalSupporters).Equals(0)
	Expect(dbIdeaBySlug.Status).Equals(models.IdeaOpen)
	Expect(dbIdeaBySlug.Title).Equals("My new idea")
	Expect(dbIdeaBySlug.Description).Equals("with this description")
	Expect(dbIdeaBySlug.User.ID).Equals(1)
	Expect(dbIdeaBySlug.User.Name).Equals("Jon Snow")
	Expect(dbIdeaBySlug.User.Email).Equals("jon.snow@got.com")
}

func TestIdeaStorage_GetInvalid(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)

	dbIdea, err := ideas.GetByID(1)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(dbIdea).IsNil()
}

func TestIdeaStorage_AddAndReturnComments(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	ideas.SetCurrentUser(jonSnow)
	ideas.AddComment(idea, "Comment #1")
	ideas.SetCurrentUser(aryaStark)
	ideas.AddComment(idea, "Comment #2")

	comments, err := ideas.GetCommentsByIdea(idea)
	Expect(err).IsNil()
	Expect(comments).HasLen(2)

	Expect(comments[0].Content).Equals("Comment #1")
	Expect(comments[0].User.Name).Equals("Jon Snow")
	Expect(comments[1].Content).Equals("Comment #2")
	Expect(comments[1].User.Name).Equals("Arya Stark")
}

func TestIdeaStorage_AddGetUpdateComment(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	commentId, err := ideas.AddComment(idea, "Comment #1")
	Expect(err).IsNil()

	comment, err := ideas.GetCommentByID(commentId)
	Expect(err).IsNil()
	Expect(comment.ID).Equals(commentId)
	Expect(comment.Content).Equals("Comment #1")
	Expect(comment.User.ID).Equals(jonSnow.ID)
	Expect(comment.EditedOn).IsNil()
	Expect(comment.EditedBy).IsNil()

	ideas.SetCurrentUser(aryaStark)
	err = ideas.UpdateComment(commentId, "Comment #1 with edit")
	Expect(err).IsNil()

	comment, err = ideas.GetCommentByID(commentId)
	Expect(err).IsNil()
	Expect(comment.ID).Equals(commentId)
	Expect(comment.Content).Equals("Comment #1 with edit")
	Expect(comment.User.ID).Equals(jonSnow.ID)
	Expect(comment.EditedOn).IsNotNil()
	Expect(comment.EditedBy.ID).Equals(aryaStark.ID)
}

func TestIdeaStorage_AddAndGet_DifferentTenants(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	demoIdea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	ideas.SetCurrentTenant(avengersTenant)
	ideas.SetCurrentUser(tonyStark)
	avengersIdea, err := ideas.Add("My other idea", "with other description")
	Expect(err).IsNil()

	ideas.SetCurrentTenant(demoTenant)
	dbIdea, err := ideas.GetByNumber(1)

	Expect(err).IsNil()
	Expect(dbIdea.ID).Equals(demoIdea.ID)
	Expect(dbIdea.Number).Equals(1)
	Expect(dbIdea.Title).Equals("My new idea")
	Expect(dbIdea.Slug).Equals("my-new-idea")

	ideas.SetCurrentTenant(avengersTenant)
	dbIdea, err = ideas.GetByNumber(1)

	Expect(err).IsNil()
	Expect(dbIdea.ID).Equals(avengersIdea.ID)
	Expect(dbIdea.Number).Equals(1)
	Expect(dbIdea.Title).Equals("My other idea")
	Expect(dbIdea.Slug).Equals("my-other-idea")
}

func TestIdeaStorage_Update(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	idea, err = ideas.Update(idea, "The new comment", "With the new description")
	Expect(err).IsNil()
	Expect(idea.Title).Equals("The new comment")
	Expect(idea.Description).Equals("With the new description")
	Expect(idea.Slug).Equals("the-new-comment")
}

func TestIdeaStorage_AddSupporter(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	err = ideas.AddSupporter(idea, aryaStark)
	Expect(err).IsNil()

	dbIdea, err := ideas.GetByNumber(1)
	Expect(dbIdea.ViewerSupported).IsFalse()
	Expect(dbIdea.TotalSupporters).Equals(1)

	err = ideas.AddSupporter(idea, jonSnow)
	Expect(err).IsNil()

	dbIdea, err = ideas.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbIdea.ViewerSupported).IsTrue()
	Expect(dbIdea.TotalSupporters).Equals(2)
}

func TestIdeaStorage_AddSupporter_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")

	err := ideas.AddSupporter(idea, jonSnow)
	Expect(err).IsNil()

	err = ideas.AddSupporter(idea, jonSnow)
	Expect(err).IsNil()

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbIdea.TotalSupporters).Equals(1)
}

func TestIdeaStorage_RemoveSupporter(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")

	err := ideas.AddSupporter(idea, jonSnow)
	Expect(err).IsNil()

	err = ideas.RemoveSupporter(idea, jonSnow)
	Expect(err).IsNil()

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbIdea.TotalSupporters).Equals(0)
}

func TestIdeaStorage_RemoveSupporter_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")

	err := ideas.AddSupporter(idea, jonSnow)
	Expect(err).IsNil()

	err = ideas.RemoveSupporter(idea, jonSnow)
	Expect(err).IsNil()

	err = ideas.RemoveSupporter(idea, jonSnow)
	Expect(err).IsNil()

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbIdea.TotalSupporters).Equals(0)
}

func TestIdeaStorage_SetResponse(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")
	err := ideas.SetResponse(idea, "We liked this idea", models.IdeaStarted)

	Expect(err).IsNil()

	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.Text).Equals("We liked this idea")
	Expect(idea.Status).Equals(models.IdeaStarted)
	Expect(idea.Response.User.ID).Equals(1)
}

func TestIdeaStorage_SetResponse_KeepOpen(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")
	err := ideas.SetResponse(idea, "We liked this idea", models.IdeaOpen)
	Expect(err).IsNil()
}

func TestIdeaStorage_SetResponse_ChangeText(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")
	ideas.SetResponse(idea, "We liked this idea", models.IdeaStarted)
	idea, _ = ideas.GetByID(idea.ID)
	respondedOn := idea.Response.RespondedOn

	ideas.SetResponse(idea, "We liked this idea and we'll work on it", models.IdeaStarted)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.RespondedOn).Equals(respondedOn)

	ideas.SetResponse(idea, "We finished it", models.IdeaCompleted)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.RespondedOn).TemporarilySimilar(respondedOn, time.Second)
}

func TestIdeaStorage_SetResponse_AsDuplicate(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea1, _ := ideas.Add("My new idea", "with this description")
	ideas.AddSupporter(idea1, jonSnow)

	ideas.SetCurrentUser(aryaStark)
	idea2, _ := ideas.Add("My other idea", "with similar description")
	ideas.AddSupporter(idea2, aryaStark)

	ideas.SetCurrentUser(jonSnow)
	ideas.MarkAsDuplicate(idea2, idea1)
	idea1, _ = ideas.GetByID(idea1.ID)

	Expect(idea1.TotalSupporters).Equals(2)
	Expect(idea1.Status).Equals(models.IdeaOpen)
	Expect(idea1.Response).IsNil()

	idea2, _ = ideas.GetByID(idea2.ID)

	Expect(idea2.Response.Text).Equals("")
	Expect(idea2.TotalSupporters).Equals(1)
	Expect(idea2.Status).Equals(models.IdeaDuplicate)
	Expect(idea2.Response.User.ID).Equals(1)
	Expect(idea2.Response.Original.Number).Equals(idea1.Number)
	Expect(idea2.Response.Original.Title).Equals(idea1.Title)
	Expect(idea2.Response.Original.Slug).Equals(idea1.Slug)
	Expect(idea2.Response.Original.Status).Equals(idea1.Status)
}

func TestIdeaStorage_SetResponse_AsDeleted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, err := ideas.Add("My new idea", "with this description")
	Expect(err).IsNil()

	ideas.SetResponse(idea, "Spam!", models.IdeaDeleted)

	idea1, err := ideas.GetByNumber(idea.Number)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(idea1).IsNil()

	idea2, err := ideas.GetByID(idea.ID)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(idea2).IsNil()
}

func TestIdeaStorage_AddSupporter_ClosedIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, _ := ideas.Add("My new idea", "with this description")
	ideas.SetResponse(idea, "We liked this idea", models.IdeaCompleted)
	ideas.AddSupporter(idea, jonSnow)

	dbIdea, err := ideas.GetByNumber(idea.Number)
	Expect(err).IsNil()
	Expect(dbIdea.TotalSupporters).Equals(0)
}

func TestIdeaStorage_RemoveSupporter_ClosedIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea, _ := ideas.Add("My new idea", "with this description")
	ideas.AddSupporter(idea, jonSnow)
	ideas.SetResponse(idea, "We liked this idea", models.IdeaCompleted)
	ideas.RemoveSupporter(idea, jonSnow)

	dbIdea, err := ideas.GetByNumber(idea.Number)
	Expect(err).IsNil()
	Expect(dbIdea.TotalSupporters).Equals(1)
}

func TestIdeaStorage_ListSupportedIdeas(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea1, _ := ideas.Add("My new idea", "with this description")
	idea2, _ := ideas.Add("My other idea", "with better description")
	ideas.AddSupporter(idea1, aryaStark)
	ideas.AddSupporter(idea2, aryaStark)

	ideas.SetCurrentUser(jonSnow)
	arr, err := ideas.SupportedBy()
	Expect(err).IsNil()
	Expect(arr).Equals([]int{})

	ideas.SetCurrentUser(aryaStark)
	arr, err = ideas.SupportedBy()
	Expect(err).IsNil()
	Expect(arr).Equals([]int{idea1.ID, idea2.ID})
}

func TestIdeaStorage_WithTags(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(aryaStark)
	tags.SetCurrentTenant(demoTenant)
	tags.SetCurrentUser(jonSnow)

	idea, _ := ideas.Add("My new idea", "with this description")
	bug, _ := tags.Add("Bug", "FF0000", true)
	featureRequest, _ := tags.Add("Feature Request", "00FF00", false)

	tags.AssignTag(bug, idea)
	tags.AssignTag(featureRequest, idea)

	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Tags).HasLen(1)
	Expect(idea.Tags[0]).Equals(bug.Slug)

	ideas.SetCurrentUser(jonSnow)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Tags).HasLen(2)
	Expect(idea.Tags[0]).Equals(bug.Slug)
	Expect(idea.Tags[1]).Equals(featureRequest.Slug)
}

func TestIdeaStorage_IsReferenced(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	idea1, _ := ideas.Add("My first idea", "with this description")
	idea2, _ := ideas.Add("My second idea", "with this description")
	idea3, _ := ideas.Add("My third idea", "with this description")

	ideas.MarkAsDuplicate(idea2, idea3)
	ideas.MarkAsDuplicate(idea3, idea1)

	referenced, err := ideas.IsReferenced(idea1)
	Expect(referenced).IsTrue()
	Expect(err).IsNil()

	referenced, err = ideas.IsReferenced(idea2)
	Expect(referenced).IsFalse()
	Expect(err).IsNil()

	referenced, err = ideas.IsReferenced(idea3)
	Expect(referenced).IsTrue()
	Expect(err).IsNil()
}
