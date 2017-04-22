package postgres_test

import (
	"testing"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestIdeaStorage_GetAll(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	now := time.Now()

	db.Execute("INSERT INTO ideas (title, number, description, created_on, tenant_id, user_id, supporters) VALUES ('Idea #1', 1, 'Description #1', $1, 300, 300, 0)", now)
	db.Execute("INSERT INTO ideas (title, number, description, created_on, tenant_id, user_id, supporters) VALUES ('Idea #2', 2, 'Description #2', $1, 300, 301, 5)", now)

	ideas := &postgres.IdeaStorage{DB: db}
	dbIdeas, err := ideas.GetAll(300)

	Expect(err).To(BeNil())
	Expect(dbIdeas).To(HaveLen(2))

	Expect(dbIdeas[0].Title).To(Equal("Idea #1"))
	Expect(dbIdeas[0].Number).To(Equal(1))
	Expect(dbIdeas[0].Description).To(Equal("Description #1"))
	Expect(dbIdeas[0].User.Name).To(Equal("Jon Snow"))
	Expect(dbIdeas[0].TotalSupporters).To(Equal(0))

	Expect(dbIdeas[1].Title).To(Equal("Idea #2"))
	Expect(dbIdeas[1].Number).To(Equal(2))
	Expect(dbIdeas[1].Description).To(Equal("Description #2"))
	Expect(dbIdeas[1].User.Name).To(Equal("Arya Stark"))
	Expect(dbIdeas[1].TotalSupporters).To(Equal(5))
}

func TestIdeaStorage_SaveAndGet(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	db.Execute("INSERT INTO users (name, email, created_on, role) VALUES ('Jon Snow','jon.snow@got.com', now(), 2)")

	ideas := &postgres.IdeaStorage{DB: db}
	idea, err := ideas.Save(1, 1, "My new idea", "with this description")
	Expect(err).To(BeNil())
	Expect(idea.ID).To(Equal(1))

	dbIdea, err := ideas.GetByID(1, 1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(1))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.TotalSupporters).To(Equal(0))
	Expect(dbIdea.Title).To(Equal("My new idea"))
	Expect(dbIdea.Description).To(Equal("with this description"))
	Expect(dbIdea.User.ID).To(Equal(1))
	Expect(dbIdea.User.Name).To(Equal("Jon Snow"))
	Expect(dbIdea.User.Email).To(Equal("jon.snow@got.com"))
}

func TestIdeaStorage_GetInvalid(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	dbIdea, err := ideas.GetByID(1, 1)

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbIdea).To(BeNil())
}

func TestIdeaStorage_AddAndReturnComments(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	db.Execute("INSERT INTO users (name, email, created_on, role) VALUES ('Jon Snow','jon.snow@got.com', now(), 2)")

	ideas := &postgres.IdeaStorage{DB: db}
	idea, _ := ideas.Save(1, 1, "My new idea", "with this description")
	ideas.AddComment(1, idea.ID, "Comment #1")
	ideas.AddComment(1, idea.ID, "Comment #2")

	comments, err := ideas.GetCommentsByIdeaID(1, idea.ID)
	Expect(err).To(BeNil())
	Expect(len(comments)).To(Equal(2))

	Expect(comments[0].Content).To(Equal("Comment #2"))
	Expect(comments[1].Content).To(Equal("Comment #1"))
}

func TestIdeaStorage_SaveAndGet_DifferentTenants(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	ideas.Save(300, 300, "My new idea", "with this description")
	ideas.Save(400, 400, "My other idea", "with other description")

	dbIdea, err := ideas.GetByNumber(300, 1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(1))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My new idea"))

	dbIdea, err = ideas.GetByNumber(400, 1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(2))
	Expect(dbIdea.Number).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My other idea"))
}

func TestIdeaStorage_AddSupporter(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	idea, _ := ideas.Save(300, 300, "My new idea", "with this description")

	err := ideas.AddSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(301, idea.ID)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(300, 1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(2))
}

func TestIdeaStorage_AddSupporter_Twice(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	idea, _ := ideas.Save(300, 300, "My new idea", "with this description")

	err := ideas.AddSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	err = ideas.AddSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(300, 1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(1))
}

func TestIdeaStorage_RemoveSupporter(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	idea, _ := ideas.Save(300, 300, "My new idea", "with this description")

	err := ideas.AddSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(300, 1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}

func TestIdeaStorage_RemoveSupporter_Twice(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	ideas := &postgres.IdeaStorage{DB: db}
	idea, _ := ideas.Save(300, 300, "My new idea", "with this description")

	err := ideas.AddSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	err = ideas.RemoveSupporter(300, idea.ID)
	Expect(err).To(BeNil())

	dbIdea, err := ideas.GetByNumber(300, 1)
	Expect(err).To(BeNil())
	Expect(dbIdea.TotalSupporters).To(Equal(0))
}
