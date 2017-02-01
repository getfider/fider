package services_test

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy/services"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

var _ = Describe("TenantService.GetByDomain", func() {
	It("TestTenantByDomain_WhenSubdomainIsUnknown_ShouldReturnError", func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnError(sql.ErrNoRows)

		svc := &services.PostgresTenantService{DB: db}
		tenant, err := svc.GetByDomain("mydomain")

		Expect(tenant).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

	It("TestTenantByDomain_WhenSubdomainIsKnown_ShouldReturnTenantInformation", func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
		mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

		svc := &services.PostgresTenantService{DB: db}
		tenant, err := svc.GetByDomain("mydomain")

		Expect(tenant.ID).To(Equal(234))
		Expect(tenant.Name).To(Equal("My Domain Inc."))
		Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
		Expect(err).To(BeNil())
	})

	It("TestTenantByDomain_WhenCompleteDomainIsUsed_ShouldReturnTenantInformation", func() {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
		mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

		svc := &services.PostgresTenantService{DB: db}
		tenant, err := svc.GetByDomain("mydomain.anydomain.com")

		Expect(tenant.ID).To(Equal(234))
		Expect(tenant.Name).To(Equal("My Domain Inc."))
		Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
		Expect(err).To(BeNil())
	})
})
