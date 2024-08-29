package tasks

import (
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/worker"
)

func CreateUserListCompany(tenant entity.Tenant, user entity.User) worker.Task {
	return describe("Create UserList Company", func(c *worker.Context) error {
		log.Debugf(c, "Sending new tenant @{Tenant} to userlist with user email @{User}", dto.Props{
			"Tenant": tenant.Name,
			"User":   user.Email,
		})
		if err := bus.Dispatch(c, &cmd.CreateUserListCompany{
			Name:          tenant.Name,
			TenantId:      tenant.ID,
			SignedUpAt:    time.Now().Format(time.RFC3339),
			BillingStatus: enum.BillingTrial.String(),
			Subdomain:     tenant.Subdomain,
			UserId:        user.ID,
			UserEmail:     user.Email,
			UserName:      user.Name,
		}); err != nil {
			return c.Failure(err)
		}
		return nil
	})
}
