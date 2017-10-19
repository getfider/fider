package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestUserStorage_GetByID(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByID(1)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Tenant.ID).To(Equal(1))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(len(user.Providers)).To(Equal(1))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
}

func TestUserStorage_GetByEmail_Error(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByEmail(1, "unknown@got.com")

	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_GetByEmail(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByEmail(1, "jon.snow@got.com")

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Tenant.ID).To(Equal(1))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(len(user.Providers)).To(Equal(1))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
}

func TestUserStorage_GetByProvider(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByProvider(1, "facebook", "FB1234")

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Tenant.ID).To(Equal(1))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(len(user.Providers)).To(Equal(1))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
}

func TestUserStorage_GetByProvider_WrongTenant(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByProvider(2, "facebook", "FB1234")

	Expect(user).To(BeNil())
	Expect(err).To(Equal(app.ErrNotFound))
}

func TestUserStorage_GetByEmail_WrongTenant(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user, err := users.GetByEmail(2, "jon.snow@got.com")

	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_Register(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	user := &models.User{
		Name:  "Rob Stark",
		Email: "rob.stark@got.com",
		Tenant: &models.Tenant{
			ID: 1,
		},
		Role: models.RoleCollaborator,
		Providers: []*models.UserProvider{
			{
				UID:  "123123123",
				Name: oauth.FacebookProvider,
			},
		},
	}
	err := users.Register(user)
	Expect(err).To(BeNil())

	user, err = users.GetByEmail(1, "rob.stark@got.com")
	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(5)))
	Expect(user.Role).To(Equal(models.RoleCollaborator))
	Expect(user.Name).To(Equal("Rob Stark"))
	Expect(user.Email).To(Equal("rob.stark@got.com"))
}

func TestUserStorage_Register_MultipleProviders(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	trx.Execute("INSERT INTO tenants (name, subdomain, created_on, status) VALUES ('My Domain Inc.','mydomain', now(), 1)")

	users := postgres.NewUserStorage(nil, trx)
	user := &models.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Tenant: &models.Tenant{
			ID: 3,
		},
		Role: models.RoleCollaborator,
		Providers: []*models.UserProvider{
			{
				UID:  "123123123",
				Name: oauth.FacebookProvider,
			},
			{
				UID:  "456456456",
				Name: oauth.GoogleProvider,
			},
		},
	}
	err := users.Register(user)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(5)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
}

func TestUserStorage_RegisterProvider(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	users.RegisterProvider(1, &models.UserProvider{
		UID:  "GO1234",
		Name: "google",
	})
	user, err := users.GetByID(1)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(1)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
	Expect(user.Tenant.ID).To(Equal(1))

	Expect(len(user.Providers)).To(Equal(2))
	Expect(user.Providers[0].UID).To(Equal("FB1234"))
	Expect(user.Providers[0].Name).To(Equal("facebook"))
	Expect(user.Providers[1].UID).To(Equal("GO1234"))
	Expect(user.Providers[1].Name).To(Equal("google"))
}

func TestUserStorage_UpdateSettings(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	users := postgres.NewUserStorage(nil, trx)
	err := users.Update(1, &models.UpdateUserSettings{Name: "Jon Stark"})
	Expect(err).To(BeNil())

	user, err := users.GetByEmail(1, "jon.snow@got.com")
	Expect(user.Name).To(Equal("Jon Stark"))
}

func TestUserStorage_ChangeRole(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	db.Seed()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	tenants := postgres.NewTenantStorage(trx)

	users := postgres.NewUserStorage(demoTenant(tenants), trx)
	err := users.ChangeRole(1, models.RoleVisitor)
	Expect(err).To(BeNil())

	user, err := users.GetByEmail(1, "jon.snow@got.com")
	Expect(user.Role).To(Equal(models.RoleVisitor))
}
