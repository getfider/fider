package postgres_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"

	"github.com/getfider/fider/app"
)

func TestTenantStorage_Add_Activate(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantPending)

	Expect(err).IsNil()
	Expect(tenant).IsNotNil()

	tenant, err = tenants.GetByDomain("mydomain")
	Expect(err).IsNil()
	Expect(tenant.Name).Equals("My Domain Inc.")
	Expect(tenant.Subdomain).Equals("mydomain")
	Expect(tenant.Status).Equals(models.TenantPending)
	Expect(tenant.IsPrivate).IsFalse()

	err = tenants.Activate(tenant.ID)
	Expect(err).IsNil()

	tenant, err = tenants.GetByDomain("mydomain")
	Expect(err).IsNil()
	Expect(tenant.Name).Equals("My Domain Inc.")
	Expect(tenant.Subdomain).Equals("mydomain")
	Expect(tenant.Status).Equals(models.TenantActive)
	Expect(tenant.IsPrivate).IsFalse()
}

func TestTenantStorage_SingleTenant_Add(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	os.Setenv("HOST_MODE", "single")

	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantPending)
	Expect(err).IsNil()
	Expect(tenant).IsNotNil()
}

func TestTenantStorage_First(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, err := tenants.First()
	Expect(err).IsNil()
	Expect(tenant.ID).Equals(1)
}

func TestTenantStorage_Empty_First(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	trx.Execute("TRUNCATE tenants CASCADE")

	tenant, err := tenants.First()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(tenant).IsNil()
}

func TestTenantStorage_UpdatePrivacy(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)
	Expect(tenant.IsPrivate).IsFalse()

	tenants.UpdatePrivacy(&models.UpdateTenantPrivacy{IsPrivate: true})
	tenant, _ = tenants.GetByDomain("demo")
	Expect(tenant.IsPrivate).IsTrue()
}

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, err := tenants.GetByDomain("mydomain")

	Expect(tenant).IsNil()
	Expect(err).IsNotNil()
}

func TestTenantStorage_GetByDomain_Subdomain(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantActive)
	Expect(err).IsNil()

	tenant, err = tenants.GetByDomain("mydomain")
	Expect(err).IsNil()
	Expect(tenant.ID).NotEquals(0)
	Expect(tenant.Name).Equals("My Domain Inc.")
	Expect(tenant.Subdomain).Equals("mydomain")
	Expect(tenant.CNAME).Equals("")
	Expect(tenant.Status).Equals(models.TenantActive)
}

func TestTenantStorage_GetByDomain_CNAME(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, err := tenants.Add("My Domain Inc.", "mydomain", models.TenantActive)
	Expect(err).IsNil()
	tenants.SetCurrentTenant(tenant)
	tenants.UpdateSettings(&models.UpdateTenantSettings{
		Title: "My Domain Inc.",
		CNAME: "feedback.mycompany.com",
	})

	tenant, err = tenants.GetByDomain("feedback.mycompany.com")
	Expect(err).IsNil()
	Expect(tenant.ID).NotEquals(0)
	Expect(tenant.Name).Equals("My Domain Inc.")
	Expect(tenant.Subdomain).Equals("mydomain")
	Expect(tenant.CNAME).Equals("feedback.mycompany.com")
	Expect(tenant.Status).Equals(models.TenantActive)
}

func TestTenantStorage_IsSubdomainAvailable_ExistingDomain(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	available, err := tenants.IsSubdomainAvailable("demo")
	Expect(available).IsFalse()
	Expect(err).IsNil()
}

func TestTenantStorage_IsSubdomainAvailable_NewDomain(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	available, err := tenants.IsSubdomainAvailable("thisisanewdomain")
	Expect(available).IsTrue()
	Expect(err).IsNil()
}

func TestTenantStorage_UpdateSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	settings := &models.UpdateTenantSettings{
		Title:          "New Demonstration",
		Invitation:     "Leave us your suggestion",
		WelcomeMessage: "Welcome!",
		CNAME:          "demo.company.com",
	}
	err := tenants.UpdateSettings(settings)
	Expect(err).IsNil()

	tenant, err = tenants.GetByDomain("demo")
	Expect(err).IsNil()
	Expect(tenant.ID).Equals(1)
	Expect(tenant.Name).Equals("New Demonstration")
	Expect(tenant.Invitation).Equals("Leave us your suggestion")
	Expect(tenant.WelcomeMessage).Equals("Welcome!")
	Expect(tenant.CNAME).Equals("demo.company.com")
	Expect(tenant.LogoID).Equals(0)
}

func TestTenantStorage_UpdateSettings_WithLogo(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	logo, _ := ioutil.ReadFile(env.Path("./favicon.ico"))

	settings := &models.UpdateTenantSettings{
		Logo: &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content: logo,
			},
		},
		Title:          "New Demonstration",
		Invitation:     "Leave us your suggestion",
		WelcomeMessage: "Welcome!",
		CNAME:          "demo.company.com",
	}
	err := tenants.UpdateSettings(settings)
	Expect(err).IsNil()

	upload, err := tenants.GetUpload(tenant.LogoID)
	Expect(err).IsNil()
	Expect(upload.Content).Equals(logo)
	Expect(upload.Size).Equals(len(logo))
	Expect(upload.ContentType).Equals("image/vnd.microsoft.icon")

	//Remove Logo
	settings.Logo.Upload = nil
	settings.Logo.Remove = true
	err = tenants.UpdateSettings(settings)
	Expect(err).IsNil()

	tenant, err = tenants.GetByDomain("demo")
	Expect(err).IsNil()
	Expect(tenant.LogoID).Equals(0)
}

func TestTenantStorage_UpdateSettings_ReplaceLogo(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	logo, _ := ioutil.ReadFile(env.Path("./favicon.ico"))

	settings := &models.UpdateTenantSettings{
		Logo: &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content: logo,
			},
		},
		Title:          "New Demonstration",
		Invitation:     "Leave us your suggestion",
		WelcomeMessage: "Welcome!",
		CNAME:          "demo.company.com",
	}
	err := tenants.UpdateSettings(settings)
	Expect(err).IsNil()

	firstLogoID := tenant.LogoID
	upload, err := tenants.GetUpload(firstLogoID)
	Expect(err).IsNil()
	Expect(upload.Content).Equals(logo)
	Expect(upload.Size).Equals(len(logo))
	Expect(upload.ContentType).Equals("image/vnd.microsoft.icon")

	//Replace logo with a new one
	newLogo, _ := ioutil.ReadFile(env.Path("./README.md"))
	settings.Logo.Upload.Content = newLogo
	err = tenants.UpdateSettings(settings)
	Expect(err).IsNil()

	Expect(tenant.LogoID).NotEquals(firstLogoID)

	upload, err = tenants.GetUpload(tenant.LogoID)
	Expect(err).IsNil()
	Expect(upload.Content).Equals(newLogo)
	Expect(upload.Size).Equals(len(newLogo))
	Expect(upload.ContentType).Equals("text/html; charset=utf-8")
}

func TestTenantStorage_AdvancedSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	settings := &models.UpdateTenantAdvancedSettings{
		CustomCSS: ".primary { color: red; }",
	}
	err := tenants.UpdateAdvancedSettings(settings)
	Expect(err).IsNil()

	tenant, err = tenants.GetByDomain("demo")
	Expect(err).IsNil()
	Expect(tenant.CustomCSS).Equals(".primary { color: red; }")
}

func TestTenantStorage_SaveFindSet_VerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	//Save new Key
	request := &models.CreateTenant{
		Email: "jon.snow@got.com",
		Name:  "Jon Snow",
	}
	err := tenants.SaveVerificationKey("s3cr3tk3y", 15*time.Minute, request)
	Expect(err).IsNil()

	//Find and check values
	result, err := tenants.FindVerificationByKey(models.EmailVerificationKindSignUp, "s3cr3tk3y")
	Expect(err).IsNil()
	Expect(result.CreatedAt).TemporarilySimilar(time.Now(), 1*time.Second)
	Expect(result.VerifiedAt).IsNil()
	Expect(result.Email).Equals("jon.snow@got.com")
	Expect(result.Name).Equals("Jon Snow")
	Expect(result.Kind).Equals(models.EmailVerificationKindSignUp)
	Expect(result.Key).Equals("s3cr3tk3y")
	Expect(result.UserID).Equals(0)
	Expect(result.ExpiresAt).TemporarilySimilar(result.CreatedAt.Add(15*time.Minute), 1*time.Second)

	//Set as verified check values
	err = tenants.SetKeyAsVerified("s3cr3tk3y")
	Expect(err).IsNil()

	//Find and check that VerifiedAt is now set
	result, err = tenants.FindVerificationByKey(models.EmailVerificationKindSignUp, "s3cr3tk3y")
	Expect(err).IsNil()
	Expect(time.Now().After(result.CreatedAt)).IsTrue()
	Expect(result.VerifiedAt.After(result.CreatedAt)).IsTrue()
	Expect(result.Email).Equals("jon.snow@got.com")
	Expect(result.Name).Equals("Jon Snow")
	Expect(result.Key).Equals("s3cr3tk3y")
	Expect(result.UserID).Equals(0)
	Expect(result.ExpiresAt).TemporarilySimilar(result.CreatedAt.Add(15*time.Minute), 1*time.Second)

	//Wrong kind should not find it
	result, err = tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "s3cr3tk3y")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(result).IsNil()
}

func TestTenantStorage_SaveFindSet_ChangeEmailVerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	//Save new Key
	request := &models.ChangeUserEmail{
		Email:     "jon.stark@got.com",
		Requestor: jonSnow,
	}
	err := tenants.SaveVerificationKey("th3-s3cr3t", 15*time.Minute, request)
	Expect(err).IsNil()

	//Find and check values
	result, err := tenants.FindVerificationByKey(models.EmailVerificationKindChangeEmail, "th3-s3cr3t")
	Expect(err).IsNil()
	Expect(result.CreatedAt).TemporarilySimilar(time.Now(), 1*time.Second)
	Expect(result.VerifiedAt).IsNil()
	Expect(result.Email).Equals("jon.stark@got.com")
	Expect(result.Name).Equals("")
	Expect(result.Kind).Equals(models.EmailVerificationKindChangeEmail)
	Expect(result.Key).Equals("th3-s3cr3t")
	Expect(result.UserID).Equals(jonSnow.ID)
	Expect(result.ExpiresAt).TemporarilySimilar(result.CreatedAt.Add(15*time.Minute), 1*time.Second)
}

func TestTenantStorage_FindUnknownVerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	//Find and check values
	result, err := tenants.FindVerificationByKey(models.EmailVerificationKindSignIn, "blahblahblah")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(result).IsNil()
}

func TestTenantStorage_Save_Get_ListOAuthConfig(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	tenant, _ := tenants.GetByDomain("demo")
	tenants.SetCurrentTenant(tenant)

	config, err := tenants.GetOAuthConfigByProvider("_TEST")
	Expect(config).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)

	logo, _ := ioutil.ReadFile(env.Path("./favicon.ico"))

	err = tenants.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		ID: 0,
		Logo: &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content:     logo,
				ContentType: "image/vnd.microsoft.icon",
			},
		},
		Provider:          "_TEST",
		DisplayName:       "My Provider",
		ClientID:          "823187ahjjfdha8fds7yfdashfjkdsa",
		ClientSecret:      "jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij",
		AuthorizeURL:      "http://provider/oauth/authorize",
		TokenURL:          "http://provider/oauth/token",
		Scope:             "profile email",
		ProfileURL:        "http://provider/profile/me",
		JSONUserIDPath:    "user.id",
		JSONUserNamePath:  "user.name",
		JSONUserEmailPath: "user.email",
	})
	Expect(err).IsNil()

	config, err = tenants.GetOAuthConfigByProvider("_TEST")
	Expect(err).IsNil()
	Expect(config.ID).Equals(1)
	Expect(config.Provider).Equals("_TEST")
	Expect(config.DisplayName).Equals("My Provider")
	Expect(config.ClientID).Equals("823187ahjjfdha8fds7yfdashfjkdsa")
	Expect(config.ClientSecret).Equals("jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
	Expect(config.AuthorizeURL).Equals("http://provider/oauth/authorize")
	Expect(config.TokenURL).Equals("http://provider/oauth/token")
	Expect(config.Scope).Equals("profile email")
	Expect(config.Status).Equals(0)
	Expect(config.ProfileURL).Equals("http://provider/profile/me")
	Expect(config.JSONUserIDPath).Equals("user.id")
	Expect(config.JSONUserNamePath).Equals("user.name")
	Expect(config.JSONUserEmailPath).Equals("user.email")

	upload, err := tenants.GetUpload(config.LogoID)
	Expect(err).IsNil()
	Expect(upload.Content).Equals(logo)
	Expect(upload.ContentType).Equals("image/vnd.microsoft.icon")

	err = tenants.SaveOAuthConfig(&models.CreateEditOAuthConfig{
		ID: config.ID,
		Logo: &models.ImageUpload{
			Remove: true,
		},
		Provider:          "_TEST2222", //this has to be ignored
		DisplayName:       "New My Provider",
		ClientID:          "New 823187ahjjfdha8fds7yfdashfjkdsa",
		ClientSecret:      "New jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij",
		AuthorizeURL:      "New http://provider/oauth/authorize",
		TokenURL:          "New http://provider/oauth/token",
		Scope:             "New profile email",
		ProfileURL:        "New http://provider/profile/me",
		JSONUserIDPath:    "New user.id",
		JSONUserNamePath:  "New user.name",
		JSONUserEmailPath: "New user.email",
	})
	Expect(err).IsNil()

	configs, err := tenants.ListOAuthConfig()
	Expect(err).IsNil()
	Expect(configs).HasLen(1)
	Expect(configs[0].ID).Equals(1)
	Expect(configs[0].LogoID).Equals(0)
	Expect(configs[0].Provider).Equals("_TEST")
	Expect(configs[0].DisplayName).Equals("New My Provider")
	Expect(configs[0].ClientID).Equals("New 823187ahjjfdha8fds7yfdashfjkdsa")
	Expect(configs[0].ClientSecret).Equals("New jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
	Expect(configs[0].AuthorizeURL).Equals("New http://provider/oauth/authorize")
	Expect(configs[0].TokenURL).Equals("New http://provider/oauth/token")
	Expect(configs[0].Scope).Equals("New profile email")
	Expect(configs[0].Status).Equals(0)
	Expect(configs[0].ProfileURL).Equals("New http://provider/profile/me")
	Expect(configs[0].JSONUserIDPath).Equals("New user.id")
	Expect(configs[0].JSONUserNamePath).Equals("New user.name")
	Expect(configs[0].JSONUserEmailPath).Equals("New user.email")
}
