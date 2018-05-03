package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/oauth"
)

func TestUserStorage_GetByID(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	user, err := users.GetByID(1)
	Expect(err).IsNil()
	Expect(user.ID).Equals(int(1))
	Expect(user.Tenant.ID).Equals(1)
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].UID).Equals("FB1234")
	Expect(user.Providers[0].Name).Equals("facebook")
}

func TestUserStorage_GetByEmail_Error(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByEmail("unknown@got.com")
	Expect(user).IsNil()
	Expect(err).IsNotNil()
}

func TestUserStorage_GetByEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.ID).Equals(int(1))
	Expect(user.Tenant.ID).Equals(1)
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].UID).Equals("FB1234")
	Expect(user.Providers[0].Name).Equals("facebook")
}

func TestUserStorage_GetByProvider(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	user, err := users.GetByProvider("facebook", "FB1234")
	Expect(err).IsNil()
	Expect(user.ID).Equals(int(1))
	Expect(user.Tenant.ID).Equals(1)
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].UID).Equals("FB1234")
	Expect(user.Providers[0].Name).Equals("facebook")
}

func TestUserStorage_GetByProvider_WrongTenant(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(avengersTenant)
	user, err := users.GetByProvider("facebook", "FB1234")
	Expect(user).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
}

func TestUserStorage_GetByEmail_WrongTenant(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(avengersTenant)
	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(user).IsNil()
	Expect(err).IsNotNil()
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
	Expect(err).IsNil()

	user, err = users.GetByEmail("rob.stark@got.com")
	Expect(err).IsNil()
	Expect(user.ID).Equals(int(6))
	Expect(user.Role).Equals(models.RoleCollaborator)
	Expect(user.Name).Equals("Rob Stark")
	Expect(user.Email).Equals("rob.stark@got.com")
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
	Expect(err).IsNil()

	user, err = users.GetByID(user.ID)
	Expect(err).IsNil()
	Expect(user.Role).Equals(models.RoleCollaborator)
	Expect(user.Name).Equals("Rob Stark")
	Expect(user.Email).Equals("")
}

func TestUserStorage_Register_MultipleProviders(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	var tenantID int
	trx.Get(&tenantID, `
		INSERT INTO tenants (name, subdomain, created_on, status, is_private) 
		VALUES ('My Domain Inc.','mydomain', now(), 1, false)
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
	Expect(err).IsNil()
	Expect(user.ID).NotEquals(0)
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
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

	Expect(err).IsNil()
	Expect(user.ID).Equals(int(1))
	Expect(user.Name).Equals("Jon Snow")
	Expect(user.Email).Equals("jon.snow@got.com")
	Expect(user.Tenant.ID).Equals(1)

	Expect(user.Providers).HasLen(2)
	Expect(user.Providers[0].UID).Equals("FB1234")
	Expect(user.Providers[0].Name).Equals("facebook")
	Expect(user.Providers[1].UID).Equals("GO1234")
	Expect(user.Providers[1].Name).Equals("google")
}

func TestUserStorage_UpdateSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(jonSnow)
	err := users.Update(&models.UpdateUserSettings{Name: "Jon Stark"})
	Expect(err).IsNil()

	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Jon Stark")
}

func TestUserStorage_ChangeRole(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	err := users.ChangeRole(jonSnow.ID, models.RoleVisitor)
	Expect(err).IsNil()

	user, err := users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Role).Equals(models.RoleVisitor)
}

func TestUserStorage_ChangeEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	err := users.ChangeEmail(jonSnow.ID, "jon.stark@got.com")
	Expect(err).IsNil()

	user, err := users.GetByEmail("jon.stark@got.com")
	Expect(err).IsNil()
	Expect(user.Email).Equals("jon.stark@got.com")

	user, err = users.GetByEmail("jon.snow@got.com")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(user).IsNil()
}

func TestUserStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)

	all, err := users.GetAll()
	Expect(err).IsNil()
	Expect(all).HasLen(3)
	Expect(all[0].Name).Equals("Jon Snow")
	Expect(all[1].Name).Equals("Arya Stark")
	Expect(all[2].Name).Equals("Sansa Stark")
}

func TestUserStorage_DefaultUserSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentUser(jonSnow)
	settings, _ := users.GetUserSettings()
	Expect(settings).Equals(map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      models.NotificationEventNewIdea.DefaultSettingValue,
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: models.NotificationEventChangeStatus.DefaultSettingValue,
	})

	users.SetCurrentUser(aryaStark)
	settings, _ = users.GetUserSettings()
	Expect(settings).Equals(map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: models.NotificationEventChangeStatus.DefaultSettingValue,
	})
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
	Expect(settings).Equals(map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   models.NotificationEventNewComment.DefaultSettingValue,
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	})

	users.UpdateSettings(nil)
	newSettings, _ := users.GetUserSettings()
	Expect(newSettings).Equals(settings)
}
