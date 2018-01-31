package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// ChangeUserEmail register the intent of changing user e-mail
func ChangeUserEmail() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.ChangeUserEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.SaveVerificationKey(input.Model.VerificationKey, 24*time.Hour, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}
		err = c.Services().Emailer.Send(c.Tenant().Name, input.Model.Email, "Confirm your new e-mail", "change_emailaddress_email", web.Map{
			"name":            c.User().Name,
			"oldEmail":        previous,
			"newEmail":        input.Model.Email,
			"baseUrl":         c.BaseURL(),
			"verificationKey": input.Model.VerificationKey,
		})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// VerifyChangeEmailKey checks if key is correct and update user's email
func VerifyChangeEmailKey() web.HandlerFunc {
	return func(c web.Context) error {
		result, err := validateKey(models.EmailVerificationKindChangeEmail, c)
		if err != nil {
			return err
		}

		if result.UserID != c.User().ID {
			return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}

		err = c.Services().Users.ChangeEmail(result.UserID, result.Email)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.SetKeyAsVerified(result.Key)
		if err != nil {
			return c.Failure(err)
		}
		return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL()+"/settings")
	}
}

// UpdateUserSettings updates current user settings
func UpdateUserSettings() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateUserSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.Update(c.User().ID, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ChangeUserRole changes given user role
func ChangeUserRole() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.ChangeUserRole)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.ChangeRole(input.Model.UserID, input.Model.Role)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
