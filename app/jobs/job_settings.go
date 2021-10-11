package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
)

type Context struct {
	context.Context
	LastSuccessfulRun *time.Time
}

func getLastSuccessfulRun(ctx context.Context, jobName string) *time.Time {
	key := fmt.Sprintf("jobs.%s.last_successful_run", jobName)
	get := &query.GetSystemSettings{
		Key: key,
	}
	if err := bus.Dispatch(ctx, get); err != nil {
		log.Error(ctx, err)
	}

	if get.Value != "" {
		t, err := time.Parse(time.RFC3339, get.Value)
		if err != nil {
			log.Error(ctx, err)
			return nil
		}
		return &t
	}
	return nil
}

func setLastSuccessfulRun(jobName string, t time.Time) {
	key := fmt.Sprintf("jobs.%s.last_successful_run", jobName)
	setLastRun(key, t)
}

func setLastFailedRun(jobName string, t time.Time) {
	key := fmt.Sprintf("jobs.%s.last_failed_run", jobName)
	setLastRun(key, t)
}

func setLastRun(key string, t time.Time) {
	ctx, trx, err := newJobContext()
	if err != nil {
		log.Error(ctx, err)
		return
	}
	defer trx.MustCommit()

	if err = bus.Dispatch(ctx, &cmd.SetSystemSettings{
		Key:   key,
		Value: t.Format(time.RFC3339),
	}); err != nil {
		log.Error(ctx, err)
	}
}
