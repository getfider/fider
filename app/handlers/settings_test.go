package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSettingsHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCurrentUserSettings) error {
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		AsUser(mock.JonSnow).
		Execute(handlers.UserSettings())

	Expect(code).Equals(http.StatusOK)
}

func TestUpdateUserSettingsHandler_EmptyInput(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		AsUser(mock.JonSnow).
		ExecutePost(handlers.UpdateUserSettings(), `{ }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdateUserSettingsHandler_ValidName(t *testing.T) {
	RegisterT(t)

	var newName string
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateCurrentUser) error {
		newName = c.Name
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateCurrentUserSettings) error {
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.UpdateUserSettings(), `{ "name": "Jon Stark", "avatarType": "gravatar" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newName).Equals("Jon Stark")
}

func TestUpdateUserSettingsHandler_NewSettings(t *testing.T) {
	RegisterT(t)

	var updateCmd *cmd.UpdateCurrentUser
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateCurrentUser) error {
		updateCmd = c
		return nil
	})

	var updateSettingsCmd *cmd.UpdateCurrentUserSettings
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateCurrentUserSettings) error {
		updateSettingsCmd = c
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(handlers.UpdateUserSettings(), `{
			"name": "Jon Stark",
			"avatarType": "gravatar",
			"settings": {
				"event_notification_new_post": "1",
				"event_notification_new_comment": "2",
				"event_notification_change_status": "3"
			}
		}`)

	Expect(code).Equals(http.StatusOK)
	Expect(updateCmd.Name).Equals("Jon Stark")
	Expect(updateCmd.AvatarType).Equals(models.AvatarTypeGravatar)

	Expect(updateSettingsCmd.Settings[models.NotificationEventNewPost.UserSettingsKeyName]).Equals("1")
	Expect(updateSettingsCmd.Settings[models.NotificationEventNewComment.UserSettingsKeyName]).Equals("2")
	Expect(updateSettingsCmd.Settings[models.NotificationEventChangeStatus.UserSettingsKeyName]).Equals("3")
}

func TestChangeRoleHandler_Valid(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == mock.AryaStark.ID {
			q.Result = mock.AryaStark
			return nil
		}
		return app.ErrNotFound
	})

	var changeRole *cmd.ChangeUserRole
	bus.AddHandler(func(ctx context.Context, c *cmd.ChangeUserRole) error {
		changeRole = c
		return nil
	})

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("role", models.RoleAdministrator).
		ExecutePost(handlers.ChangeUserRole(), fmt.Sprintf(`{ "userID": %d }`, mock.AryaStark.ID))

	Expect(code).Equals(http.StatusOK)
	Expect(changeRole.UserID).Equals(mock.AryaStark.ID)
	Expect(changeRole.Role).Equals(models.RoleAdministrator)
}

func TestChangeUserEmailHandler_Valid(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	for _, email := range []string{
		"jon.another@got.com",
		"another.snow@got.com",
	} {
		server, _ := mock.NewServer()
		code, _ := server.
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(handlers.ChangeUserEmail(), fmt.Sprintf(`{ "email": "%s" }`, email))

		Expect(code).Equals(http.StatusOK)
	}
}

func TestChangeUserEmailHandler_Invalid(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == mock.JonSnow.Email {
			q.Result = mock.JonSnow
			return nil
		}

		if q.Email == mock.AryaStark.Email {
			q.Result = mock.AryaStark
			return nil
		}

		return app.ErrNotFound
	})

	for _, email := range []string{
		"",
		"jon.snow@got.com",
		"jon.snow",
		"arya.stark@got.com",
	} {
		server, _ := mock.NewServer()
		code, _ := server.
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(handlers.ChangeUserEmail(), fmt.Sprintf(`{ "email": "%s" }`, email))

		Expect(code).Equals(http.StatusBadRequest)
	}
}

func TestVerifyChangeEmailKeyHandler_Success(t *testing.T) {
	RegisterT(t)

	var changeEmailCmd *cmd.ChangeUserEmail
	bus.AddHandler(func(ctx context.Context, c *cmd.ChangeUserEmail) error {
		changeEmailCmd = c
		return nil
	})

	server, services := mock.NewServer()
	services.Tenants.SaveVerificationKey("th3-s3cr3t", 24*time.Hour, &models.ChangeUserEmail{
		Requestor: mock.JonSnow,
		Email:     "jon.stark@got.com",
	})
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/change-email/verify?k=th3-s3cr3t").
		Execute(handlers.VerifyChangeEmailKey())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(changeEmailCmd.UserID).Equals(mock.JonSnow.ID)
	Expect(changeEmailCmd.Email).Equals("jon.stark@got.com")

	result, err := services.Tenants.FindVerificationByKey(models.EmailVerificationKindChangeEmail, "th3-s3cr3t")
	Expect(err).IsNil()
	Expect(result.VerifiedAt).IsNotNil()
}

func TestVerifyChangeEmailKeyHandler_DifferentUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	request := &models.ChangeUserEmail{
		Requestor: mock.JonSnow,
		Email:     "jon.stark@got.com",
	}
	services.Tenants.SaveVerificationKey("th3-s3cr3t", 24*time.Hour, request)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithURL("/change-email/verify?k=th3-s3cr3t").
		Execute(handlers.VerifyChangeEmailKey())

	Expect(code).Equals(http.StatusTemporaryRedirect)
}

func TestDeleteUserHandler(t *testing.T) {
	RegisterT(t)

	var deleteCmd *cmd.DeleteCurrentUser
	bus.AddHandler(func(ctx context.Context, c *cmd.DeleteCurrentUser) error {
		deleteCmd = c
		return nil
	})

	server, _ := mock.NewServer()
	code, response := server.
		AsUser(mock.JonSnow).
		Execute(handlers.DeleteUser())

	Expect(code).Equals(http.StatusOK)
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=; Path=/; Expires=")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring("Max-Age=0; HttpOnly")

	Expect(deleteCmd).IsNotNil()
}
