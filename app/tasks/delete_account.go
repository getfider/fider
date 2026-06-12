package tasks

import (
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

// SendDeleteAccountScheduledEmail notifies the account owner that the whole site is scheduled
// for deletion and gives them a link to cancel during the grace window.
func SendDeleteAccountScheduledEmail(owner *entity.User, tenantName string, scheduledAt time.Time, baseURL, cancelKey string) worker.Task {
	return describe("Send delete account scheduled email", func(c *worker.Context) error {
		to := dto.NewRecipient(owner.Name, owner.Email, dto.Props{
			"tenantName":  tenantName,
			"scheduledAt": scheduledAt.UTC().Format("2 Jan 2006 15:04 MST"),
			"cancelLink":  linkWithText("Cancel the scheduled deletion", baseURL, "/admin/danger-zone/cancel?k=%s", cancelKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         dto.Recipient{Name: "Fider"},
			To:           []dto.Recipient{to},
			TemplateName: "delete_account_requested",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}
