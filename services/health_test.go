package services_test

import (
	"fmt"

	"github.com/WeCanHearYou/wchy/services"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HealthCheckService.IsDatabaseOnline", func() {
	It("should return false if query fails", func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		mock.ExpectQuery("SELECT now()").WillReturnError(fmt.Errorf("some error"))

		svc := &services.PostgresHealthCheckService{DB: db}

		Expect(svc.IsDatabaseOnline()).To(BeFalse())

		if err := mock.ExpectationsWereMet(); err != nil {
			Fail(err.Error())
		}
	})

	It("should return true if query succeed", func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		mock.ExpectQuery("SELECT now()").WillReturnRows(sqlmock.NewRows([]string{}))

		svc := &services.PostgresHealthCheckService{DB: db}

		Expect(svc.IsDatabaseOnline()).To(BeTrue())

		if err := mock.ExpectationsWereMet(); err != nil {
			Fail(err.Error())
		}
	})
})
