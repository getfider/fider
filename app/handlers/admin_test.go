package handlers_test

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/services/blob/fs"

	"github.com/getfider/fider/app/handlers"
)

func TestUpdateSettingsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.LogoBlobKey = "logos/hello-world.png"

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.UpdateSettings(),
			`{ "title": "GoT", "invitation": "Join us!", "welcomeMessage": "Welcome to GoT Feedback Forum" }`,
		)

	tenant, _ := services.Tenants.GetByDomain("demo")
	Expect(code).Equals(http.StatusOK)
	Expect(tenant.Name).Equals("GoT")
	Expect(tenant.Invitation).Equals("Join us!")
	Expect(tenant.WelcomeMessage).Equals("Welcome to GoT Feedback Forum")
	Expect(tenant.LogoBlobKey).Equals("logos/hello-world.png")
}

func TestUpdateSettingsHandler_NewLogo(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})

	logoBytes, _ := ioutil.ReadFile(env.Etc("logo.png"))
	logoB64 := base64.StdEncoding.EncodeToString(logoBytes)

	server, services := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.UpdateSettings(), `{ 
				"title": "GoT", 
				"invitation": "Join us!", 
				"welcomeMessage": "Welcome to GoT Feedback Forum",
				"logo": {
					"upload": {
						"fileName": "picture.png",
						"contentType": "image/png",
						"content": "`+logoB64+`"
					}
				}
			}`)

	tenant, _ := services.Tenants.GetByDomain("demo")
	Expect(code).Equals(http.StatusOK)
	Expect(tenant.Name).Equals("GoT")
	Expect(tenant.Invitation).Equals("Join us!")
	Expect(tenant.WelcomeMessage).Equals("Welcome to GoT Feedback Forum")
	Expect(tenant.LogoBlobKey).ContainsSubstring("logos/")
	Expect(tenant.LogoBlobKey).ContainsSubstring("picture.png")
}

func TestUpdateSettingsHandler_RemoveLogo(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	mock.DemoTenant.LogoBlobKey = "logos/hello-world.png"

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.UpdateSettings(), `{ 
				"title": "GoT", 
				"invitation": "Join us!", 
				"welcomeMessage": "Welcome to GoT Feedback Forum",
				"logo": {
					"remove": true
				}
			}`)

	tenant, _ := services.Tenants.GetByDomain("demo")
	Expect(code).Equals(http.StatusOK)
	Expect(tenant.Name).Equals("GoT")
	Expect(tenant.Invitation).Equals("Join us!")
	Expect(tenant.WelcomeMessage).Equals("Welcome to GoT Feedback Forum")
	Expect(tenant.LogoBlobKey).Equals("")
}

func TestUpdatePrivacyHandler(t *testing.T) {
	RegisterT(t)

	var updateCmd *cmd.UpdateTenantPrivacySettings
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateTenantPrivacySettings) error {
		updateCmd = c
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.UpdatePrivacy(),
			`{ "isPrivate": true }`,
		)

	Expect(code).Equals(http.StatusOK)
	Expect(updateCmd.Settings.IsPrivate).IsTrue()
}

func TestManageMembersHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetAllUsers) error {
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		Execute(
			handlers.ManageMembers(),
		)

	Expect(code).Equals(http.StatusOK)
}
