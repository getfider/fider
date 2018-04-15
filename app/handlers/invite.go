package handlers

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/markdown"
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
			input.Model.Message = strings.Replace(input.Model.Message, app.InvitePlaceholder, "*the link to accept invitation will be here*", -1)
			to := email.NewRecipient(c.User().Name, c.User().Email, email.Params{
				"subject": input.Model.Subject,
				"message": markdown.Parse(input.Model.Message),
			})
			err := c.Services().Emailer.Send("invite_email", email.Params{}, c.Tenant().Name, to)
			if err != nil {
				return c.Failure(err)
			}
		}

		return c.Ok(web.Map{})
	}
}
