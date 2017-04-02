package postgres_test

import (
	"testing"
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/dbx"
	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

func TestIdeaService_GetAll(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	now := time.Now()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	db.Execute("INSERT INTO users (name, email, created_on) VALUES ('Jon Snow','jon.snow@got.com', now())")
	db.Execute("INSERT INTO users (name, email, created_on) VALUES ('Arya Start','arya.start@got.com', now())")
	db.Execute("INSERT INTO ideas (title, description, created_on, tenant_id, user_id) VALUES ('Idea #1','Description #1', $1, 1, 1)", now)
	db.Execute("INSERT INTO ideas (title, description, created_on, tenant_id, user_id) VALUES ('Idea #2','Description #2', $1, 1, 2)", now)

	svc := &postgres.IdeaService{DB: db}
	ideas, err := svc.GetAll(1)

	Expect(err).To(BeNil())
	Expect(ideas).To(HaveLen(2))

	Expect(ideas[0].Title).To(Equal("Idea #1"))
	Expect(ideas[0].Description).To(Equal("Description #1"))
	Expect(ideas[0].User.Name).To(Equal("Jon Snow"))

	Expect(ideas[1].Title).To(Equal("Idea #2"))
	Expect(ideas[1].Description).To(Equal("Description #2"))
	Expect(ideas[1].User.Name).To(Equal("Arya Start"))
}

func TestIdeaService_SaveAndGet(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	db.Execute("INSERT INTO users (name, email, created_on) VALUES ('Jon Snow','jon.snow@got.com', now())")

	svc := &postgres.IdeaService{DB: db}
	idea, err := svc.Save(1, 1, "My new idea", "with this description")
	Expect(err).To(BeNil())
	Expect(idea.ID).To(Equal(1))

	dbIdea, err := svc.GetByID(1, 1)

	Expect(err).To(BeNil())
	Expect(dbIdea.ID).To(Equal(1))
	Expect(dbIdea.Title).To(Equal("My new idea"))
	Expect(dbIdea.Description).To(Equal("with this description"))
	Expect(dbIdea.User.ID).To(Equal(1))
	Expect(dbIdea.User.Name).To(Equal("Jon Snow"))
	Expect(dbIdea.User.Email).To(Equal("jon.snow@got.com"))
}

func TestIdeaService_GetInvalid(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.IdeaService{DB: db}
	dbIdea, err := svc.GetByID(1, 1)

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(dbIdea).To(BeNil())
}

func TestIdeaService_AddAndReturnComments(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	db.Execute("INSERT INTO users (name, email, created_on) VALUES ('Jon Snow','jon.snow@got.com', now())")

	svc := &postgres.IdeaService{DB: db}
	idea, _ := svc.Save(1, 1, "My new idea", "with this description")
	svc.AddComment(1, idea.ID, "Comment #1")
	svc.AddComment(1, idea.ID, "Comment #2")

	comments, err := svc.GetCommentsByIdeaID(1, idea.ID)
	Expect(err).To(BeNil())
	Expect(len(comments)).To(Equal(2))

	Expect(comments[0].Content).To(Equal("Comment #2"))
	Expect(comments[1].Content).To(Equal("Comment #1"))
}
