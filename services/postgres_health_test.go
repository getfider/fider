package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestIsDatabaseOnline_ShouldReturnFalseIfQueryFails(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT now()").WillReturnError(fmt.Errorf("some error"))

	svc := NewPostgresHealthCheckService(db)
	isOnline := svc.IsDatabaseOnline()

	assert.Equal(t, false, isOnline)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestIsDatabaseOnline_ShouldReturnTrueIfQuerySucceed(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT now()").WillReturnRows(sqlmock.NewRows([]string{}))

	svc := NewPostgresHealthCheckService(db)
	isOnline := svc.IsDatabaseOnline()

	assert.Equal(t, true, isOnline)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
