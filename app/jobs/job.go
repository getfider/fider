package jobs

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

type Handler interface {
	Run(ctx context.Context)
	Schedule() string
}

type fiderJob struct {
	Name    string
	Handler Handler
}

func NewJob(ctx context.Context, name string, handler Handler) (string, fiderJob) {
	schedule := handler.Schedule()
	if env.IsDevelopment() {
		// job runs every 1 minute during development mode
		schedule = "0 * * * * *"
	}
	log.Debugf(ctx, "Job '@{JobName}' scheduled to run '@{Schedule}'", dto.Props{
		"JobName":  name,
		"Schedule": schedule,
	})
	return schedule, fiderJob{Name: name, Handler: handler}
}

func (j fiderJob) Run() {
	ctx := log.WithProperty(context.Background(), log.PropertyKeyTag, "JOBS")

	start := time.Now()
	log.Debugf(ctx, "Job '@{JobName}' started", dto.Props{
		"JobName": j.Name,
	})

	locked, unlock := dbx.TryLock(ctx, j.Name)
	if !locked {
		log.Debugf(ctx, "Job '@{JobName}' skipped, could not acquire lock", dto.Props{
			"JobName": j.Name,
		})
		return
	}

	defer unlock()
	defer func() {
		elapsedMs := time.Since(start).Nanoseconds() / int64(time.Millisecond)
		log.Debugf(ctx, "Job '@{JobName}' finished in @{ElapsedMs:magenta}ms", dto.Props{
			"ElapsedMs": elapsedMs,
			"JobName":   j.Name,
		})
	}()

	j.Handler.Run(ctx)
}
