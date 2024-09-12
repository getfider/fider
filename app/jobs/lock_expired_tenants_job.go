package jobs

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

type LockExpiredTenantsJobHandler struct {
}

func (e LockExpiredTenantsJobHandler) Schedule() string {
	return "0 0 0 * * *" // every day at 0:00
}

func (e LockExpiredTenantsJobHandler) Run(ctx Context) error {
	c := &cmd.LockExpiredTenants{}
	err := bus.Dispatch(ctx, c)
	if err != nil {
		return err
	}

	log.Debugf(ctx, "@{Count} tenants marked as locked", dto.Props{
		"Count": c.NumOfTenantsLocked,
	})

	// Handle userlist
	if env.Config.UserList.Enabled && c.NumOfTenantsLocked > 0 {
		for _, tenant := range c.TenantsLocked {
			err := bus.Dispatch(ctx, &cmd.UserListUpdateCompany{TenantId: tenant, BillingStatus: enum.BillingCancelled})
			if err != nil {
				return err
			}
		}

	}

	return nil
}
