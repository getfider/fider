package tasks

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

// SignUpEmailData contains the data needed to send a signup email
type SignUpEmailData interface {
	GetEmail() string
	GetName() string
	GetVerificationKey() string
}

// SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(data SignUpEmailData, baseURL string) worker.Task {
	return describe("Send sign up email", func(c *worker.Context) error {
		to := dto.NewRecipient(data.GetName(), data.GetEmail(), dto.Props{
			"link": link(baseURL, "/signup/verify?k=%s", data.GetVerificationKey()),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: "Fider"},
			To:           []dto.Recipient{to},
			TemplateName: "signup_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}
