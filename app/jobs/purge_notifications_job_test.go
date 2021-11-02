package jobs_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/jobs"
	"github.com/getfider/fider/app/models/cmd"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestPurgeExpiredNotificationsJob_Schedule_IsCorrect(t *testing.T) {
	RegisterT(t)

	job := &jobs.PurgeExpiredNotificationsJobHandler{}
	Expect(job.Schedule()).Equals("0 0 * * * *")
}

func TestPurgeExpiredNotificationsJob_ShouldJustDispatchCommand(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.PurgeExpiredNotifications) error {
		return nil
	})

	job := &jobs.PurgeExpiredNotificationsJobHandler{}
	err := job.Run(jobs.Context{
		Context: context.Background(),
	})
	Expect(err).IsNil()
}
