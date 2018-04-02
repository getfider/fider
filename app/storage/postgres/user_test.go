package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/oauth"
	. "github.com/onsi/gomega"
)

func TestUserStorage_GetByID(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

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
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByEmail("unknown@got.com")
	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_GetByEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByEmail("jon.snow@got.com")
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
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByProvider("facebook", "FB1234")
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
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(avengersTenant)
	user, err := users.GetByProvider("facebook", "FB1234")
	Expect(user).To(BeNil())
	Expect(errors.Cause(err)).To(Equal(app.ErrNotFound))
}

func TestUserStorage_GetByEmail_WrongTenant(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(avengersTenant)
	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserStorage_Register(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user := &models.User{
		Name:  "Rob Stark",
		Email: "rob.stark@got.com",
		Role:  models.RoleCollaborator,
		Providers: []*models.UserProvider{
			{
				UID:  "123123123",
				Name: oauth.FacebookProvider,
			},
		},
	}
	err := users.Register(user)
	Expect(err).To(BeNil())

	user, err = users.GetByEmail("rob.stark@got.com")
	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int(6)))
	Expect(user.Role).To(Equal(models.RoleCollaborator))
	Expect(user.Name).To(Equal("Rob Stark"))
	Expect(user.Email).To(Equal("rob.stark@got.com"))
}

func TestUserStorage_Register_WhiteSpaceEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user := &models.User{
		Name:  "Rob Stark",
		Email: "   ",
		Role:  models.RoleCollaborator,
	}
	err := users.Register(user)
	Expect(err).To(BeNil())

	user, err = users.GetByID(user.ID)
	Expect(err).To(BeNil())
	Expect(user.Role).To(Equal(models.RoleCollaborator))
	Expect(user.Name).To(Equal("Rob Stark"))
	Expect(user.Email).To(Equal(""))
}

func TestUserStorage_Register_MultipleProviders(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	var tenantID int
	trx.Get(&tenantID, `
		INSERT INTO tenants (name, subdomain, created_on, status) 
		VALUES ('My Domain Inc.','mydomain', now(), 1)
		RETURNING id
	`)
	users.SetCurrentTenant(&models.Tenant{
		ID: tenantID,
	})

	user := &models.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Role:  models.RoleCollaborator,
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
	Expect(user.ID).NotTo(BeZero())
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
}

func TestUserStorage_RegisterProvider(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
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
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(jonSnow)
	err := users.Update(&models.UpdateUserSettings{Name: "Jon Stark"})
	Expect(err).To(BeNil())

	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(user.Name).To(Equal("Jon Stark"))
}

func TestUserStorage_ChangeRole(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	err := users.ChangeRole(jonSnow.ID, models.RoleVisitor)
	Expect(err).To(BeNil())

	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(err).To(BeNil())
	Expect(user.Role).To(Equal(models.RoleVisitor))
}

func TestUserStorage_ChangeEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	err := users.ChangeEmail(jonSnow.ID, "jon.stark@got.com")
	Expect(err).To(BeNil())

	user, err := users.GetByEmail("jon.stark@got.com")
	Expect(err).To(BeNil())
	Expect(user.Email).To(Equal("jon.stark@got.com"))

	user, err = users.GetByEmail("jon.snow@got.com")
	Expect(errors.Cause(err)).To(Equal(app.ErrNotFound))
	Expect(user).To(BeNil())
}

func TestUserStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)

	all, err := users.GetAll()
	Expect(err).To(BeNil())
	Expect(len(all)).To(Equal(3))
	Expect(all[0].Name).To(Equal("Jon Snow"))
	Expect(all[1].Name).To(Equal("Arya Stark"))
	Expect(all[2].Name).To(Equal("Sansa Stark"))
}

func TestUserStorage_DefaultUserSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentUser(jonSnow)
	settings, _ := users.GetUserSettings()
	Expect(settings).To(Equal(map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      models.NotificationEventNewIdea.DefaultSettingValue,
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: models.NotificationEventChangeStatus.DefaultSettingValue,
	}))

	users.SetCurrentUser(aryaStark)
	settings, _ = users.GetUserSettings()
	Expect(settings).To(Equal(map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: models.NotificationEventChangeStatus.DefaultSettingValue,
	}))
}

func TestUserStorage_SaveGetUserSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)
	disableAll := map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      "0",
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}

	users.UpdateSettings(disableAll)
	settings, _ := users.GetUserSettings()
	Expect(settings).To(Equal(map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}))

	users.UpdateSettings(nil)
	newSettings, _ := users.GetUserSettings()
	Expect(newSettings).To(Equal(settings))
}
