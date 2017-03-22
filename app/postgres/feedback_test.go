package postgres_test

import (
	"database/sql/driver"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/WeCanHearYou/wechy/app/feedback"
	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

func TestIdeaService_GetAll_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT (.*) FROM ideas (.*) JOIN users").WithArgs(2134).WillReturnError(driver.ErrBadConn)

	svc := &postgres.IdeaService{DB: db}
	ideas, err := svc.GetAll(2134)

	Expect(ideas).To(BeEmpty())
	Expect(err).NotTo(BeNil())
}

func TestIdeaService_GetAll(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "created_on", "u_id", "u_name", "u_email"})
	rows.AddRow(1, "Idea #1", "Description #1", now, 1, "Jon Snow", "jon.snow@got.com")
	rows.AddRow(2, "Idea #2", "Description #2", now, 2, "Arya Start", "arya.start@got.com")
	mock.ExpectQuery("SELECT (.*) FROM ideas (.*) JOIN users").WithArgs(2134).WillReturnRows(rows)

	svc := &postgres.IdeaService{DB: db}
	ideas, err := svc.GetAll(2134)

	Expect(err).To(BeNil())
	Expect(ideas).To(HaveLen(2))
	Expect(ideas[0]).To(Equal(&feedback.Idea{
		ID:          1,
		Title:       "Idea #1",
		Description: "Description #1",
		CreatedOn:   now,
		User: identity.User{
			ID:    1,
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		},
	}))
	Expect(ideas[1]).To(Equal(&feedback.Idea{
		ID:          2,
		Title:       "Idea #2",
		Description: "Description #2",
		CreatedOn:   now,
		User: identity.User{
			ID:    2,
			Name:  "Arya Start",
			Email: "arya.start@got.com",
		},
	}))
}
