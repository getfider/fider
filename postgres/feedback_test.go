package postgres_test

import (
	"database/sql/driver"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/WeCanHearYou/wechy/feedback"
	"github.com/WeCanHearYou/wechy/postgres"
	. "github.com/onsi/gomega"
)

func TestIdeaService_GetAll_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, description FROM ideas WHERE tenant_id = \\$1").WithArgs(2134).WillReturnError(driver.ErrBadConn)

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
	rows := sqlmock.NewRows([]string{"id", "title", "description", "created_on"})
	rows.AddRow(1, "Idea #1", "Description #1", now)
	rows.AddRow(2, "Idea #2", "Description #2", now)

	mock.ExpectQuery("SELECT id, title, description, created_on FROM ideas WHERE tenant_id = \\$1").WithArgs(2134).WillReturnRows(rows)

	svc := &postgres.IdeaService{DB: db}
	ideas, err := svc.GetAll(2134)

	Expect(ideas).To(HaveLen(2))
	Expect(ideas[0]).To(Equal(&feedback.Idea{ID: 1, Title: "Idea #1", Description: "Description #1", CreatedOn: now}))
	Expect(ideas[1]).To(Equal(&feedback.Idea{ID: 2, Title: "Idea #2", Description: "Description #2", CreatedOn: now}))
	Expect(err).To(BeNil())
}
