package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/web"
)

// SendSampleInvite to current user's email
func SendSampleInvite() web.HandlerFunc {
	return func(c web.Context) error {
		input := &actions.InviteUsers{IsSampleInvite: true}
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		if c.User().Email != "" {
			to := email.NewRecipient(c.User().Name, c.User().Email, email.Params{
				"subject": input.Model.Subject,
				"message": input.Model.Message,
			})
			err := c.Services().Emailer.Send("invite_email", email.Params{}, c.Tenant().Name, to)
			if err != nil {
				return c.Failure(err)
			}
		}

		return c.Ok(web.Map{})
	}
}
