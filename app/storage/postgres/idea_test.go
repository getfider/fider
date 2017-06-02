package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

var demoTenant = &models.Tenant{
	ID:   300,
	Name: "Demonstration",
}

var orangeTenant = &models.Tenant{
	ID:   400,
	Name: "Orange Inc.",
}

func TestIdeaStorage_GetAll(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	now := time.Now()

	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('Idea #1', 'idea-1', 1, 'Description #1', $1, 300, 300, 0, 1)", now)
	trx.Execute("INSERT INTO ideas (title, slug, number, description, created_on, tenant_id, user_id, supporters, status) VALUES ('Idea #2', 'idea-2', 2, 'Description #2', $1, 300, 301, 5, 2)", now)

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
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
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, err := ideas.Add("My new idea", "with this description", 300)
	Expect(err).To(BeNil())
	Expect(idea.ID).To(Equal(1))

	dbIdea, err := ideas.GetByID(1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(1))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.TotalSupporters).To(Equal(0))
	Expect(dbIdea.Status).To(Equal(models.IdeaNew))
	Expect(dbIdea.Title).To(Equal("My new idea"))
	Expect(dbIdea.Description).To(Equal("with this description"))
	Expect(dbIdea.User.ID).To(Equal(300))
	Expect(dbIdea.User.Name).To(Equal("Jon Snow"))
	Expect(dbIdea.User.Email).To(Equal("jon.snow@got.com"))
}

func TestIdeaStorage_GetInvalid(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	dbIdea, err := ideas.GetByID(1)

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbIdea).To(BeNil())
}

func TestIdeaStorage_AddAndReturnComments(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, err := ideas.Add("My new idea", "with this description", 300)
	Expect(err).To(BeNil())

	ideas.AddComment(idea.Number, "Comment #1", 300)
	ideas.AddComment(idea.Number, "Comment #2", 300)

	comments, err := ideas.GetCommentsByIdea(idea.Number)
	Expect(err).To(BeNil())
	Expect(len(comments)).To(Equal(2))

	Expect(comments[0].Content).To(Equal("Comment #2"))
	Expect(comments[1].Content).To(Equal("Comment #1"))
}

func TestIdeaStorage_AddAndGet_DifferentTenants(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	demoIdeas := postgres.NewIdeaStorage(demoTenant, trx)
	demoIdeas.Add("My new idea", "with this description", 300)

	orangeIdeas := postgres.NewIdeaStorage(orangeTenant, trx)
	orangeIdeas.Add("My other idea", "with other description", 400)

	dbIdea, err := demoIdeas.GetByNumber(1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(1))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My new idea"))
	Expect(dbIdea.Slug).To(Equal("my-new-idea"))

	dbIdea, err = orangeIdeas.GetByNumber(1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(2))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My other idea"))
	Expect(dbIdea.Slug).To(Equal("my-other-idea"))

}

func TestIdeaStorage_AddSupporter(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)

	err := ideas.AddSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(idea.Number, 301)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(2))
}

func TestIdeaStorage_AddSupporter_Twice(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)

	err := ideas.AddSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(1))
}

func TestIdeaStorage_RemoveSupporter(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)

	err := ideas.AddSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_RemoveSupporter_Twice(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)

	err := ideas.AddSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(idea.Number, 300)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_SetResponse(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)
	err := ideas.SetResponse(idea.Number, "We liked this idea", 300, models.IdeaStarted)

	Expect(err).To(BeNil())

	idea, _ = ideas.GetByID(idea.ID)
	Expect(idea.Response.Text).To(Equal("We liked this idea"))
	Expect(idea.Status).To(Equal(models.IdeaStarted))
	Expect(idea.Response.User.ID).To(Equal(300))
}

func TestIdeaStorage_AddSupporter_ClosedIdea(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)
	ideas.SetResponse(idea.Number, "We liked this idea", 300, models.IdeaCompleted)
	ideas.AddSupporter(idea.Number, 300)

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_RemoveSupporter_ClosedIdea(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	ideas := postgres.NewIdeaStorage(demoTenant, trx)
	idea, _ := ideas.Add("My new idea", "with this description", 300)
	ideas.AddSupporter(idea.Number, 300)
	ideas.SetResponse(idea.Number, "We liked this idea", 300, models.IdeaCompleted)
	ideas.RemoveSupporter(idea.Number, 300)

	dbIdea, err := ideas.GetByNumber(1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(1))
}
