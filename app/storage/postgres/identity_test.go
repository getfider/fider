package postgres_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantStorage_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(err).To(BeNil())
}

func TestTenantStorage_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, cname, created_on) VALUES ('My Domain Inc.','mydomain', 'mydomain.anydomain.com', now())")

	tenants := &postgres.TenantStorage{DB: db}
	tenant, err := tenants.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(err).To(BeNil())
}
