package service_test

import (
	"fmt"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/WeCanHearYou/wchy/app/service"
	. "github.com/onsi/gomega"
)

func TestHealthCheckService_IsDatabaseOnline_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT now()").WillReturnError(fmt.Errorf("some error"))

	svc := &service.PostgresHealthCheckService{DB: db}

	Expect(svc.IsDatabaseOnline()).To(BeFalse())

	Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
}

func TestHealthCheckService_IsDatabaseOnline_Success(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT now()").WillReturnRows(sqlmock.NewRows([]string{}))

	svc := &service.PostgresHealthCheckService{DB: db}

	Expect(svc.IsDatabaseOnline()).To(BeTrue())
	Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
}
