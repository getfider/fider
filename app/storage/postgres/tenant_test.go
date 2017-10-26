package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestTenantStorage_Add_Activate(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantInactive)

	Expect(err).To(BeNil())
	Expect(tenant).NotTo(BeNil())

	tenant, err = tenants.GetByDomain("mydomain")
	Expect(err).To(BeNil())
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.Status).To(Equal(models.TenantInactive))

	err = tenants.Activate(tenant.ID)
	Expect(err).To(BeNil())

	tenant, err = tenants.GetByDomain("mydomain")
	Expect(err).To(BeNil())
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.Status).To(Equal(models.TenantActive))
}

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
	Expect(tenant.ID).To(Equal(1))
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
	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantActive)
	Expect(err).To(BeNil())

	tenant, err = tenants.GetByDomain("mydomain")

	Expect(err).To(BeNil())
	Expect(tenant.ID).To(Equal(int(3)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Subdomain).To(Equal("mydomain"))
	Expect(tenant.CNAME).To(Equal(""))
	Expect(tenant.Status).To(Equal(models.TenantActive))
}

func TestTenantStorage_IsSubdomainAvailable_ExistingDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	available, err := tenants.IsSubdomainAvailable("demo")

	Expect(available).To(BeFalse())
	Expect(err).To(BeNil())
}

func TestTenantStorage_IsSubdomainAvailable_NewDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	available, err := tenants.IsSubdomainAvailable("thisisanewdomain")

	Expect(available).To(BeTrue())
	Expect(err).To(BeNil())
}

func TestTenantStorage_UpdateSettings(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	settings := &models.UpdateTenantSettings{
		Title:          "New Demonstration",
		Invitation:     "Leave us your suggestion",
		WelcomeMessage: "Welcome!",
		CNAME:          "demo.company.com",
	}
	err := tenants.UpdateSettings(settings)
	Expect(err).To(BeNil())

	tenant, err = tenants.GetByDomain("demo")
	Expect(err).To(BeNil())
	Expect(tenant.ID).To(Equal(1))
	Expect(tenant.Name).To(Equal("New Demonstration"))
	Expect(tenant.Invitation).To(Equal("Leave us your suggestion"))
	Expect(tenant.WelcomeMessage).To(Equal("Welcome!"))
	Expect(tenant.CNAME).To(Equal("demo.company.com"))
}

func TestTenantStorage_SaveFindSet_VerificationKey(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	//Save new Key
	err := tenants.SaveVerificationKey("s3cr3tk3y", 15*time.Minute, "jon.snow@got.com", "Jon Snow")
	Expect(err).To(BeNil())

	//Find and check values
	result, err := tenants.FindVerificationByKey("s3cr3tk3y")
	Expect(err).To(BeNil())
	Expect(result.CreatedOn).NotTo(Equal(time.Time{}))
	Expect(result.VerifiedOn).To(BeNil())
	Expect(result.Email).To(Equal("jon.snow@got.com"))
	Expect(result.Name).To(Equal("Jon Snow"))
	Expect(result.Key).To(Equal("s3cr3tk3y"))
	Expect(result.ExpiresOn).To(BeTemporally("~", result.CreatedOn.Add(15*time.Minute), 1*time.Second))

	//Set as verified check values
	err = tenants.SetKeyAsVerified("s3cr3tk3y")
	Expect(err).To(BeNil())

	//Find and check that VerifiedOn is now set
	result, err = tenants.FindVerificationByKey("s3cr3tk3y")
	Expect(err).To(BeNil())
	Expect(time.Now().After(result.CreatedOn)).To(BeTrue())
	Expect(result.VerifiedOn.After(result.CreatedOn)).To(BeTrue())
	Expect(result.Email).To(Equal("jon.snow@got.com"))
	Expect(result.Name).To(Equal("Jon Snow"))
	Expect(result.Key).To(Equal("s3cr3tk3y"))
	Expect(result.ExpiresOn).To(BeTemporally("~", result.CreatedOn.Add(15*time.Minute), 1*time.Second))
}

func TestTenantStorage_FindUnknownVerificationKey(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)
	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	//Find and check values
	result, err := tenants.FindVerificationByKey("blahblahblah")
	Expect(err).To(Equal(app.ErrNotFound))
	Expect(result).To(BeNil())
}
