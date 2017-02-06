package service_test

import (
	"database/sql"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/WeCanHearYou/wchy/service"
	. "github.com/onsi/gomega"
)

func TestTenantService_GetByDomain_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnError(sql.ErrNoRows)

	svc := &service.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantService_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := &service.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(234))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}

func TestTenantService_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := &service.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(234))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}
