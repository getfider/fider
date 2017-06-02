package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTenantStorage_First(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, err := tenants.First()

	Expect(err).To(BeNil())
	Expect(tenant.ID).To(Equal(300))
}

func TestTenantStorage_Empty_First(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	trx.Execute("TRUNCATE tenants CASCADE")

	tenants := postgres.NewTenantStorage(trx)
	tenant, err := tenants.First()

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(tenant).To(BeNil())
}

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantStorage_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenants.Add(&models.Tenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
	})

	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.CNAME).To(Equal(""))
	Expect(err).To(BeNil())
}

func TestTenantStorage_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenants.Add(&models.Tenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
		CNAME:     "mydomain.anydomain.com",
	})

	tenant, err := tenants.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.CNAME).To(Equal("mydomain.anydomain.com"))
	Expect(err).To(BeNil())
}
