package tasks

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

//SendChangeEmailConfirmation is used to send the change email confirmation email to requestor
func SendChangeEmailConfirmation(action *actions.ChangeUserEmail) worker.Task {
	return describe("Send change email confirmation", func(c *worker.Context) error {

		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}

		to := dto.NewRecipient(action.Requestor.Name, action.Email, dto.Props{
			"name":     c.User().Name,
			"oldEmail": previous,
			"newEmail": action.Email,
			"link":     link(web.BaseURL(c), "/change-email/verify?k=%s", action.VerificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: c.Tenant().Name},
			To:           []dto.Recipient{to},
			TemplateName: "change_emailaddress_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}
