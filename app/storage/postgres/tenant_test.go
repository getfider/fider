package postgres_test

import (
	"testing"

	"database/sql"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTenantStorage_First(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.First()

	Expect(err).To(BeNil())
	Expect(tenant.ID).To(Equal(300))
}

func TestTenantStorage_Empty_First(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	db.Execute("TRUNCATE tenants CASCADE")

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.First()

	Expect(err).To(Equal(app.ErrNotFound))
	Expect(tenant).To(BeNil())
}

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantStorage_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	tenants := &postgres.TenantStorage{DB: db}
	tenants.Add(&models.Tenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
	})

	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.CNAME.Valid).To(BeFalse())
	Expect(err).To(BeNil())
}

func TestTenantStorage_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	tenants := &postgres.TenantStorage{DB: db}
	tenants.Add(&models.Tenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
		CNAME:     sql.NullString{String: "mydomain.anydomain.com", Valid: true},
	})

	tenant, err := tenants.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.CNAME.String).To(Equal("mydomain.anydomain.com"))
	Expect(tenant.CNAME.Valid).To(BeTrue())
	Expect(err).To(BeNil())
}
