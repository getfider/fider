package jobs

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
)

type PurgeExpiredNotificationsJobHandler struct {
}

func (e PurgeExpiredNotificationsJobHandler) Schedule() string {
	return "0 0 * * * *" // every hour at minute 0
}

func (e PurgeExpiredNotificationsJobHandler) Run(ctx Context) error {
	log.Debug(ctx, "deleting notifications older than 1 year")

	c := &cmd.PurgeExpiredNotifications{}
	err := bus.Dispatch(ctx, c)
	if err != nil {
		return err
	}

	log.Debugf(ctx, "@{RowsDeleted} notifications were deleted", dto.Props{
		"RowsDeleted": c.NumOfDeletedNotifications,
	})

	return nil
}
