package tasks

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/worker"
)

// SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(email, verificationKey string) worker.Task {
	return describe("Send sign in email", func(c *worker.Context) error {
		to := dto.NewRecipient("", email, dto.Props{
			"siteName": c.Tenant().Name,
			"link":     link(web.BaseURL(c), "/signin/verify?k=%s", verificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: c.Tenant().Name},
			To:           []dto.Recipient{to},
			TemplateName: "signin_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}
