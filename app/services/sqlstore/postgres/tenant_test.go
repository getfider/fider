package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestTenantStorage_Add_Activate(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	createTenant := &cmd.CreateTenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
		Status:    enum.TenantPending,
	}

	err := bus.Dispatch(ctx, createTenant)
	Expect(err).IsNil()
	Expect(createTenant.Result).IsNotNil()

	getByDomain := &query.GetTenantByDomain{Domain: "mydomain"}
	err = bus.Dispatch(ctx, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.Name).Equals("My Domain Inc.")
	Expect(getByDomain.Result.Subdomain).Equals("mydomain")
	Expect(getByDomain.Result.Status).Equals(enum.TenantPending)
	Expect(getByDomain.Result.IsPrivate).IsFalse()

	err = bus.Dispatch(ctx, &cmd.ActivateTenant{TenantID: createTenant.Result.ID})
	Expect(err).IsNil()

	err = bus.Dispatch(ctx, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.Name).Equals("My Domain Inc.")
	Expect(getByDomain.Result.Subdomain).Equals("mydomain")
	Expect(getByDomain.Result.Status).Equals(enum.TenantActive)
	Expect(getByDomain.Result.IsPrivate).IsFalse()
}

func TestTenantStorage_SingleTenant_Add(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	env.Config.HostMode = "single"

	createTenant := &cmd.CreateTenant{
		Name:      "My Domain Inc.",
		Subdomain: "mydomain",
		Status:    enum.TenantPending,
	}
	err := bus.Dispatch(ctx, createTenant)
	Expect(err).IsNil()
	Expect(createTenant.Result).IsNotNil()
}

func TestTenantStorage_First(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	getFirst := &query.GetFirstTenant{}
	err := bus.Dispatch(ctx, getFirst)
	Expect(err).IsNil()
	Expect(getFirst.Result.ID).Equals(1)
	Expect(getFirst.Result.Name).Equals("Demonstration")
}

func TestTenantStorage_Empty_First(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	_, _ = trx.Execute("TRUNCATE tenants CASCADE")

	getFirst := &query.GetFirstTenant{}
	err := bus.Dispatch(ctx, getFirst)

	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(getFirst.Result).IsNil()
}

func TestTenantStorage_UpdatePrivacy(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	setPrivate := &cmd.UpdateTenantPrivacySettings{
		IsPrivate: true,
	}
	getByDomain := &query.GetTenantByDomain{
		Domain: "demo",
	}
	err := bus.Dispatch(demoTenantCtx, setPrivate, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.IsPrivate).IsTrue()
}

func TestTenantStorage_GetByDomain_NotFound(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	getByDomain := &query.GetTenantByDomain{
		Domain: "unknown",
	}
	err := bus.Dispatch(ctx, getByDomain)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(getByDomain.Result).IsNil()
}

func TestTenantStorage_GetByDomain_CNAME(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	err := bus.Dispatch(demoTenantCtx, &cmd.UpdateTenantSettings{
		Title: "My Domain Inc.",
		CNAME: "feedback.mycompany.com",
		Logo:  &dto.ImageUpload{},
	})
	Expect(err).IsNil()

	getByDomain := &query.GetTenantByDomain{Domain: "feedback.mycompany.com"}
	err = bus.Dispatch(demoTenantCtx, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.ID).NotEquals(0)
	Expect(getByDomain.Result.Name).Equals("My Domain Inc.")
	Expect(getByDomain.Result.Subdomain).Equals("demo")
	Expect(getByDomain.Result.CNAME).Equals("feedback.mycompany.com")
	Expect(getByDomain.Result.Status).Equals(enum.TenantActive)
}

func TestTenantStorage_IsSubdomainAvailable_ExistingDomain(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	isAvailable := &query.IsSubdomainAvailable{Subdomain: "demo"}
	err := bus.Dispatch(ctx, isAvailable)
	Expect(err).IsNil()
	Expect(isAvailable.Result).IsFalse()
}

func TestTenantStorage_IsSubdomainAvailable_NewDomain(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	isAvailable := &query.IsSubdomainAvailable{Subdomain: "thisisanewdomain"}
	err := bus.Dispatch(ctx, isAvailable)
	Expect(err).IsNil()
	Expect(isAvailable.Result).IsTrue()
}

func TestTenantStorage_UpdateSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	err := bus.Dispatch(demoTenantCtx, &cmd.UpdateTenantSettings{
		Logo: &dto.ImageUpload{
			BlobKey: "some-logo-key.png",
		},
		Title:          "New Demonstration",
		Invitation:     "Leave us your suggestion",
		WelcomeMessage: "Welcome!",
		CNAME:          "demo.company.com",
	})
	Expect(err).IsNil()

	getByDomain := &query.GetTenantByDomain{Domain: "demo"}
	err = bus.Dispatch(demoTenantCtx, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.ID).Equals(1)
	Expect(getByDomain.Result.Name).Equals("New Demonstration")
	Expect(getByDomain.Result.Invitation).Equals("Leave us your suggestion")
	Expect(getByDomain.Result.WelcomeMessage).Equals("Welcome!")
	Expect(getByDomain.Result.CNAME).Equals("demo.company.com")
	Expect(getByDomain.Result.LogoBlobKey).Equals("some-logo-key.png")
}

func TestTenantStorage_AdvancedSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	err := bus.Dispatch(demoTenantCtx, &cmd.UpdateTenantAdvancedSettings{
		CustomCSS: ".primary { color: red; }",
	})
	Expect(err).IsNil()

	getByDomain := &query.GetTenantByDomain{Domain: "demo"}
	err = bus.Dispatch(demoTenantCtx, getByDomain)
	Expect(err).IsNil()
	Expect(getByDomain.Result.CustomCSS).Equals(".primary { color: red; }")
}

func TestTenantStorage_SaveFindSet_VerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	//Save new Key
	err := bus.Dispatch(demoTenantCtx, &cmd.SaveVerificationKey{
		Key:      "s3cr3tk3y",
		Duration: 15 * time.Minute,
		Request: &actions.CreateTenant{
			Email: "jon.snow@got.com",
			Name:  "Jon Snow",
		},
	})
	Expect(err).IsNil()

	//Find and check values
	getKey := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignUp, Key: "s3cr3tk3y"}

	err = bus.Dispatch(demoTenantCtx, getKey)
	Expect(err).IsNil()
	Expect(getKey.Result.CreatedAt).TemporarilySimilar(time.Now(), 1*time.Second)
	Expect(getKey.Result.VerifiedAt).IsNil()
	Expect(getKey.Result.Email).Equals("jon.snow@got.com")
	Expect(getKey.Result.Name).Equals("Jon Snow")
	Expect(getKey.Result.Kind).Equals(enum.EmailVerificationKindSignUp)
	Expect(getKey.Result.Key).Equals("s3cr3tk3y")
	Expect(getKey.Result.UserID).Equals(0)
	Expect(getKey.Result.ExpiresAt).TemporarilySimilar(getKey.Result.CreatedAt.Add(15*time.Minute), 1*time.Second)

	//Set as verified check values
	err = bus.Dispatch(demoTenantCtx, &cmd.SetKeyAsVerified{Key: "s3cr3tk3y"})
	Expect(err).IsNil()

	//Find and check that VerifiedAt is now set
	err = bus.Dispatch(demoTenantCtx, getKey)

	Expect(err).IsNil()
	Expect(time.Now().After(getKey.Result.CreatedAt)).IsTrue()
	Expect(getKey.Result.VerifiedAt.After(getKey.Result.CreatedAt)).IsTrue()
	Expect(getKey.Result.Email).Equals("jon.snow@got.com")
	Expect(getKey.Result.Name).Equals("Jon Snow")
	Expect(getKey.Result.Key).Equals("s3cr3tk3y")
	Expect(getKey.Result.UserID).Equals(0)
	Expect(getKey.Result.ExpiresAt).TemporarilySimilar(getKey.Result.CreatedAt.Add(15*time.Minute), 1*time.Second)

	//Wrong kind should not find it
	getKeyWithWrongKind := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignIn, Key: "s3cr3tk3y"}
	err = bus.Dispatch(demoTenantCtx, getKeyWithWrongKind)

	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(getKeyWithWrongKind.Result).IsNil()
}

func TestTenantStorage_SaveFindSet_ChangeEmailVerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	//Save new Key
	err := bus.Dispatch(demoTenantCtx, &cmd.SaveVerificationKey{
		Key:      "th3-s3cr3t",
		Duration: 15 * time.Minute,
		Request: &actions.ChangeUserEmail{
			Email:     "jon.stark@got.com",
			Requestor: jonSnow,
		},
	})
	Expect(err).IsNil()

	//Find and check values
	getKey := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindChangeEmail, Key: "th3-s3cr3t"}
	err = bus.Dispatch(demoTenantCtx, getKey)
	Expect(err).IsNil()
	Expect(getKey.Result.CreatedAt).TemporarilySimilar(time.Now(), 1*time.Second)
	Expect(getKey.Result.VerifiedAt).IsNil()
	Expect(getKey.Result.Email).Equals("jon.stark@got.com")
	Expect(getKey.Result.Name).Equals("")
	Expect(getKey.Result.Kind).Equals(enum.EmailVerificationKindChangeEmail)
	Expect(getKey.Result.Key).Equals("th3-s3cr3t")
	Expect(getKey.Result.UserID).Equals(jonSnow.ID)
	Expect(getKey.Result.ExpiresAt).TemporarilySimilar(getKey.Result.CreatedAt.Add(15*time.Minute), 1*time.Second)
}

func TestTenantStorage_FindUnknownVerificationKey(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	//Find and check values
	getKey := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignIn, Key: "blahblahblah"}
	err := bus.Dispatch(demoTenantCtx, getKey)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(getKey.Result).IsNil()
}

func TestTenantStorage_Save_Get_ListOAuthConfig(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	getConfig := &query.GetCustomOAuthConfigByProvider{Provider: "_TEST"}
	err := bus.Dispatch(demoTenantCtx, getConfig)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(getConfig.Result).IsNil()

	err = bus.Dispatch(demoTenantCtx, &cmd.SaveCustomOAuthConfig{
		Logo: &dto.ImageUpload{
			BlobKey: "uploads/my-logo-key.png",
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

	err = bus.Dispatch(demoTenantCtx, getConfig)
	Expect(err).IsNil()
	Expect(getConfig.Result.ID).Equals(1)
	Expect(getConfig.Result.LogoBlobKey).Equals("uploads/my-logo-key.png")
	Expect(getConfig.Result.Provider).Equals("_TEST")
	Expect(getConfig.Result.DisplayName).Equals("My Provider")
	Expect(getConfig.Result.ClientID).Equals("823187ahjjfdha8fds7yfdashfjkdsa")
	Expect(getConfig.Result.ClientSecret).Equals("jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
	Expect(getConfig.Result.AuthorizeURL).Equals("http://provider/oauth/authorize")
	Expect(getConfig.Result.TokenURL).Equals("http://provider/oauth/token")
	Expect(getConfig.Result.Scope).Equals("profile email")
	Expect(getConfig.Result.Status).Equals(0)
	Expect(getConfig.Result.ProfileURL).Equals("http://provider/profile/me")
	Expect(getConfig.Result.JSONUserIDPath).Equals("user.id")
	Expect(getConfig.Result.JSONUserNamePath).Equals("user.name")
	Expect(getConfig.Result.JSONUserEmailPath).Equals("user.email")

	err = bus.Dispatch(demoTenantCtx, &cmd.SaveCustomOAuthConfig{
		ID: getConfig.Result.ID,
		Logo: &dto.ImageUpload{
			BlobKey: "",
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

	customConfigs := &query.ListCustomOAuthConfig{}
	err = bus.Dispatch(demoTenantCtx, customConfigs)
	Expect(err).IsNil()

	Expect(customConfigs.Result).HasLen(1)
	Expect(customConfigs.Result[0].ID).Equals(1)
	Expect(customConfigs.Result[0].LogoBlobKey).Equals("")
	Expect(customConfigs.Result[0].Provider).Equals("_TEST")
	Expect(customConfigs.Result[0].DisplayName).Equals("New My Provider")
	Expect(customConfigs.Result[0].ClientID).Equals("New 823187ahjjfdha8fds7yfdashfjkdsa")
	Expect(customConfigs.Result[0].ClientSecret).Equals("New jijads78d76cn347768x3t4668q275@ˆ&Tnycasdgsacuyhij")
	Expect(customConfigs.Result[0].AuthorizeURL).Equals("New http://provider/oauth/authorize")
	Expect(customConfigs.Result[0].TokenURL).Equals("New http://provider/oauth/token")
	Expect(customConfigs.Result[0].Scope).Equals("New profile email")
	Expect(customConfigs.Result[0].Status).Equals(0)
	Expect(customConfigs.Result[0].ProfileURL).Equals("New http://provider/profile/me")
	Expect(customConfigs.Result[0].JSONUserIDPath).Equals("New user.id")
	Expect(customConfigs.Result[0].JSONUserNamePath).Equals("New user.name")
	Expect(customConfigs.Result[0].JSONUserEmailPath).Equals("New user.email")
}
