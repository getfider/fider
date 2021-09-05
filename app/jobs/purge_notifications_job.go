package jobs

import (
	"context"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
)

type PurgeExpiredNotificationsJobHandler struct {
}

func (e PurgeExpiredNotificationsJobHandler) Schedule() string {
	return "0 0 * * * *" // every hour
}

func (e PurgeExpiredNotificationsJobHandler) Run(ctx context.Context) {
	bus.Publish(ctx, &cmd.PurgeExpiredNotifications{})
}
