package jobs

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
)

type LockExpiredTrialTenantsJobHandler struct {
}

func (e LockExpiredTrialTenantsJobHandler) Schedule() string {
	return "0 0 0 * * *" // every day at minute 0
}

func (e LockExpiredTrialTenantsJobHandler) Run(ctx Context) error {
	c := &cmd.LockExpiredTrialTenants{}
	err := bus.Dispatch(ctx, c)
	if err != nil {
		return err
	}

	log.Debugf(ctx, "@{Count} tenants marked as locked", dto.Props{
		"Count": c.NumOfTenantsLocked,
	})

	return nil
}
