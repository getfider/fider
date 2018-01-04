package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	. "github.com/onsi/gomega"
)

func TestIdeaStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	now := time.Now()

	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('Idea #1', 'idea-1', 1, 'Description #1', $1, 1, 1, 0, 1)", now)
	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('Idea #2', 'idea-2', 2, 'Description #2', $1, 1, 2, 5, 2)", now)

	ideas.SetCurrentTenant(demoTenant)

	dbIdeas, err := ideas.GetAll()
	Expect(err).To(BeNil())
	Expect(dbIdeas).To(HaveLen(2))

	Expect(dbIdeas[0].Title).To(Equal("Idea #1"))
	Expect(dbIdeas[0].Slug).To(Equal("idea-1"))
	Expect(dbIdeas[0].Number).To(Equal(1))
	Expect(dbIdeas[0].Description).To(Equal("Description #1"))
	Expect(dbIdeas[0].User.Name).To(Equal("Jon Snow"))
	Expect(dbIdeas[0].TotalSupporters).To(Equal(0))
	Expect(dbIdeas[0].Status).To(Equal(models.IdeaStarted))

	Expect(dbIdeas[1].Title).To(Equal("Idea #2"))
	Expect(dbIdeas[1].Slug).To(Equal("idea-2"))
	Expect(dbIdeas[1].Number).To(Equal(2))
	Expect(dbIdeas[1].Description).To(Equal("Description #2"))
	Expect(dbIdeas[1].User.Name).To(Equal("Arya Stark"))
	Expect(dbIdeas[1].TotalSupporters).To(Equal(5))
	Expect(dbIdeas[1].Status).To(Equal(models.IdeaCompleted))
}

func TestIdeaStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, err := ideas.Add("My new idea", "with this description", 1)
	Expect(err).To(BeNil())

	dbIdeaById, err := ideas.GetByID(idea.ID)
	Expect(err).To(BeNil())
	Expect(dbIdeaById.ID).To(Equal(idea.ID))
	Expect(dbIdeaById.Number).To(Equal(1))
	Expect(dbIdeaById.ViewerSupported).To(BeFalse())
	Expect(dbIdeaById.TotalSupporters).To(Equal(0))
	Expect(dbIdeaById.Status).To(Equal(models.IdeaOpen))
	Expect(dbIdeaById.Title).To(Equal("My new idea"))
	Expect(dbIdeaById.Description).To(Equal("with this description"))
	Expect(dbIdeaById.User.ID).To(Equal(1))
	Expect(dbIdeaById.User.Name).To(Equal("Jon Snow"))
	Expect(dbIdeaById.User.Email).To(Equal("jon.snow@got.com"))

	dbIdeaBySlug, err := ideas.GetBySlug("my-new-idea")

	Expect(err).To(BeNil())
	Expect(dbIdeaBySlug.ID).To(Equal(idea.ID))
	Expect(dbIdeaBySlug.Number).To(Equal(1))
	Expect(dbIdeaBySlug.ViewerSupported).To(BeFalse())
	Expect(dbIdeaBySlug.TotalSupporters).To(Equal(0))
	Expect(dbIdeaBySlug.Status).To(Equal(models.IdeaOpen))
	Expect(dbIdeaBySlug.Title).To(Equal("My new idea"))
	Expect(dbIdeaBySlug.Description).To(Equal("with this description"))
	Expect(dbIdeaBySlug.User.ID).To(Equal(1))
	Expect(dbIdeaBySlug.User.Name).To(Equal("Jon Snow"))
	Expect(dbIdeaBySlug.User.Email).To(Equal("jon.snow@got.com"))
}

func TestIdeaStorage_GetInvalid(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)

	dbIdea, err := ideas.GetByID(1)
	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbIdea).To(BeNil())
}

func TestIdeaStorage_AddAndReturnComments(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, err := ideas.Add("My new idea", "with this description", 1)
	Expect(err).To(BeNil())

	ideas.AddComment(idea.Number, "Comment #1", 1)
	ideas.AddComment(idea.Number, "Comment #2", 2)

	comments, err := ideas.GetCommentsByIdea(idea.Number)
	Expect(err).To(BeNil())
	Expect(len(comments)).To(Equal(2))

	Expect(comments[0].Content).To(Equal("Comment #1"))
	Expect(comments[0].User.Name).To(Equal("Jon Snow"))
	Expect(comments[1].Content).To(Equal("Comment #2"))
	Expect(comments[1].User.Name).To(Equal("Arya Stark"))
}

func TestIdeaStorage_AddAndGet_DifferentTenants(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	demoIdea, _ := ideas.Add("My new idea", "with this description", 1)

	ideas.SetCurrentTenant(avengersTenant)
	avengersIdea, _ := ideas.Add("My other idea", "with other description", 3)

	ideas.SetCurrentTenant(demoTenant)
	dbIdea, err := ideas.GetByNumber(1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(demoIdea.ID))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My new idea"))
	Expect(dbIdea.Slug).To(Equal("my-new-idea"))

	ideas.SetCurrentTenant(avengersTenant)
	dbIdea, err = ideas.GetByNumber(1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(avengersIdea.ID))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My other idea"))
	Expect(dbIdea.Slug).To(Equal("my-other-idea"))
}

func TestIdeaStorage_Update(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)

	idea, err := ideas.Add("My new idea", "with this description", 1)
	Expect(err).To(BeNil())

	idea, err = ideas.Update(idea.Number, "The new comment", "With the new description")
	Expect(err).To(BeNil())
	Expect(idea.Title).To(Equal("The new comment"))
	Expect(idea.Description).To(Equal("With the new description"))
	Expect(idea.Slug).To(Equal("the-new-comment"))
}

func TestIdeaStorage_AddSupporter(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea, err := ideas.Add("My new idea", "with this description", 1)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(idea.Number, aryaStark.ID)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(dbIdea.ViewerSupported).To(BeFalse())
	Expect(dbIdea.TotalSupporters).To(Equal(1))

	err = ideas.AddSupporter(idea.Number, jonSnow.ID)
	Expect(err).To(BeNil())

	dbIdea, err = ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.ViewerSupported).To(BeTrue())
	Expect(dbIdea.TotalSupporters).To(Equal(2))
}

func TestIdeaStorage_AddSupporter_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)

	err := ideas.AddSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(1))
}

func TestIdeaStorage_RemoveSupporter(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)

	err := ideas.AddSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_RemoveSupporter_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)

	err := ideas.AddSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 1)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_SetResponse(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)
	err := ideas.SetResponse(idea.Number, "We liked this idea", 1, models.IdeaStarted)

	Expect(err).To(BeNil())

	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.Text).To(Equal("We liked this idea"))
	Expect(idea.Status).To(Equal(models.IdeaStarted))
	Expect(idea.Response.User.ID).To(Equal(1))
}

func TestIdeaStorage_SetResponse_KeepOpen(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)
	err := ideas.SetResponse(idea.Number, "We liked this idea", 1, models.IdeaOpen)
	Expect(err).To(BeNil())
}

func TestIdeaStorage_SetResponse_ChangeText(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)
	ideas.SetResponse(idea.Number, "We liked this idea", 1, models.IdeaStarted)
	idea, _ = ideas.GetByID(idea.ID)
	respondedOn := idea.Response.RespondedOn

	ideas.SetResponse(idea.Number, "We liked this idea and we'll work on it", 1, models.IdeaStarted)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.RespondedOn).To(Equal(respondedOn))

	ideas.SetResponse(idea.Number, "We finished it", 1, models.IdeaCompleted)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.RespondedOn).Should(BeTemporally(">", respondedOn))
}

func TestIdeaStorage_SetResponse_AsDuplicate(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea1, _ := ideas.Add("My new idea", "with this description", 1)
	ideas.AddSupporter(idea1.Number, 1)
	idea2, _ := ideas.Add("My other idea", "with similar description", 2)
	ideas.AddSupporter(idea2.Number, 2)

	ideas.MarkAsDuplicate(idea2.Number, idea1.Number, 1)
	idea1, _ = ideas.GetByID(idea1.ID)

	Expect(idea1.TotalSupporters).To(Equal(2))
	Expect(idea1.Status).To(Equal(models.IdeaOpen))
	Expect(idea1.Response).To(BeNil())

	idea2, _ = ideas.GetByID(idea2.ID)

	Expect(idea2.Response.Text).To(Equal(""))
	Expect(idea2.TotalSupporters).To(Equal(1))
	Expect(idea2.Status).To(Equal(models.IdeaDuplicate))
	Expect(idea2.Response.User.ID).To(Equal(1))
	Expect(idea2.Response.Original.Number).To(Equal(idea1.Number))
	Expect(idea2.Response.Original.Title).To(Equal(idea1.Title))
	Expect(idea2.Response.Original.Slug).To(Equal(idea1.Slug))
	Expect(idea2.Response.Original.Status).To(Equal(idea1.Status))
}

func TestIdeaStorage_AddSupporter_ClosedIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)
	ideas.SetResponse(idea.Number, "We liked this idea", 1, models.IdeaCompleted)
	ideas.AddSupporter(idea.Number, 1)

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_RemoveSupporter_ClosedIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea, _ := ideas.Add("My new idea", "with this description", 1)
	ideas.AddSupporter(idea.Number, 1)
	ideas.SetResponse(idea.Number, "We liked this idea", 1, models.IdeaCompleted)
	ideas.RemoveSupporter(idea.Number, 1)

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(1))
}

func TestIdeaStorage_ListSupportedIdeas(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea1, _ := ideas.Add("My new idea", "with this description", 1)
	idea2, _ := ideas.Add("My other idea", "with better description", 1)
	ideas.AddSupporter(idea1.Number, 2)
	ideas.AddSupporter(idea2.Number, 2)

	Expect(ideas.SupportedBy(1)).To(Equal([]int{}))
	Expect(ideas.SupportedBy(2)).To(Equal([]int{idea1.ID, idea2.ID}))
}

func TestIdeaStorage_WithTags(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	tags.SetCurrentTenant(demoTenant)

	idea, _ := ideas.Add("My new idea", "with this description", 1)
	bug, _ := tags.Add("Bug", "FF0000", true)
	featureRequest, _ := tags.Add("Feature Request", "00FF00", false)

	tags.AssignTag(bug.ID, idea.ID, 1)
	tags.AssignTag(featureRequest.ID, idea.ID, 1)

	idea, _ = ideas.GetByID(idea.ID)
	Expect(len(idea.Tags)).To(Equal(1))
	Expect(idea.Tags[0]).To(Equal(int64(bug.ID)))

	ideas.SetCurrentUser(jonSnow)
	idea, _ = ideas.GetByID(idea.ID)
	Expect(len(idea.Tags)).To(Equal(2))
	Expect(idea.Tags[0]).To(Equal(int64(bug.ID)))
	Expect(idea.Tags[1]).To(Equal(int64(featureRequest.ID)))
}
