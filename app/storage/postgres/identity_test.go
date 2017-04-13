package postgres_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/identity"
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestUserStorage_GetByID(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.UserStorage{DB: db}
	user, err := svc.GetByID(300)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(300)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(len(user.Providers)).To(Equal(1))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
	Expect(user.Tenant.ID).To(Equal(300))
}

func TestUserStorage_GetByEmail_Error(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.UserStorage{DB: db}
	user, err := svc.GetByEmail(300, "unknown@got.com")

	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_GetByEmail(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.UserStorage{DB: db}
	user, err := svc.GetByEmail(300, "jon.snow@got.com")

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(300)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(len(user.Providers)).To(Equal(1))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
	Expect(user.Tenant.ID).To(Equal(300))
}

func TestUserStorage_GetByEmail_WrongTenant(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.UserStorage{DB: db}
	user, err := svc.GetByEmail(400, "jon.snow@got.com")

	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_Register(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.UserStorage{DB: db}
	user := &models.User{
		Name:  "Rob Stark",
		Email: "rob.stark@got.com",
		Tenant: &models.Tenant{
			ID: 300,
		},
		Role: models.RoleMember,
		Providers: []*models.UserProvider{
			{
				UID:  "123123123",
				Name: identity.OAuthFacebookProvider,
			},
		},
	}
	err := svc.Register(user)
	Expect(err).To(BeNil())

	user, err = svc.GetByEmail(300, "rob.stark@got.com")
	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Role).To(Equal(models.RoleMember))
	Expect(user.Name).To(Equal("Rob Stark"))
	Expect(user.Email).To(Equal("rob.stark@got.com"))
}

func TestUserStorage_Register_MultipleProviders(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")

	svc := &postgres.UserStorage{DB: db}
	user := &models.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Tenant: &models.Tenant{
			ID: 1,
		},
		Role: models.RoleMember,
		Providers: []*models.UserProvider{
			{
				UID:  "123123123",
				Name: identity.OAuthFacebookProvider,
			},
			{
				UID:  "456456456",
				Name: identity.OAuthGoogleProvider,
			},
		},
	}
	err := svc.Register(user)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
}

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.TenantStorage{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantStorage_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, created_on) VALUES ('My Domain Inc.','mydomain', now())")

	svc := &postgres.TenantStorage{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}

func TestTenantStorage_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	db.Execute("INSERT INTO tenants (name, subdomain, cname, created_on) VALUES ('My Domain Inc.','mydomain', 'mydomain.anydomain.com', now())")

	svc := &postgres.TenantStorage{DB: db}
	tenant, err := svc.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(int(1)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}
