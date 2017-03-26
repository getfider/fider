package postgres_test

import (
	"testing"
	"time"

	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

func TestIdeaService_GetAll(t *testing.T) {
	RegisterTestingT(t)
	db := setup()
	defer teardown(db)

	now := time.Now()

	execute(db, "INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")
	execute(db, "INSERT INTO users (name, email, created_on) VALUES ('Jon Snow','jon.snow@got.com', now())")
	execute(db, "INSERT INTO users (name, email, created_on) VALUES ('Arya Start','arya.start@got.com', now())")
	execute(db, "INSERT INTO ideas (title, description, created_on, tenant_id, user_id) VALUES ('Idea #1','Description #1', $1, 1, 1)", now)
	execute(db, "INSERT INTO ideas (title, description, created_on, tenant_id, user_id) VALUES ('Idea #2','Description #2', $1, 1, 2)", now)

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
