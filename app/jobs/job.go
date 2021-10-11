package jobs

import (
	"context"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
)

type Handler interface {
	Run(ctx Context) error
	Schedule() string
}

type fiderJob struct {
	Name    string
	Handler Handler
}

func NewJob(ctx context.Context, name string, handler Handler) (string, fiderJob) {
	schedule := handler.Schedule()
	log.Debugf(ctx, "Job '@{JobName}' scheduled to run '@{Schedule}'", dto.Props{
		"JobName":  name,
		"Schedule": schedule,
	})
	return schedule, fiderJob{Name: name, Handler: handler}
}

func (j fiderJob) Run() {
	ctx, trx, err := newJobContext()
	if err != nil {
		log.Error(ctx, err)
		return
	}

	start := time.Now()
	ctx.LastSuccessfulRun = getLastSuccessfulRun(ctx, j.Name)

	defer func() {
		if r := recover(); r != nil {
			log.Error(ctx, errors.Panicked(r))
			setLastFailedRun(j.Name, start)
			trx.MustRollback()
		}
	}()

	log.Debugf(ctx, "Job '@{JobName}' started", dto.Props{
		"JobName": j.Name,
	})

	locked, unlock := dbx.TryLock(ctx, trx, j.Name)
	if !locked {
		log.Debugf(ctx, "Job '@{JobName}' skipped, could not acquire lock", dto.Props{
			"JobName": j.Name,
		})
		trx.MustCommit()
		return
	}

	defer func() {
		elapsedMs := time.Since(start).Nanoseconds() / int64(time.Millisecond)
		log.Debugf(ctx, "Job '@{JobName}' finished in @{ElapsedMs:magenta}ms", dto.Props{
			"ElapsedMs": elapsedMs,
			"JobName":   j.Name,
		})

		// Jobs should take at least 1sec before unlocking to avoid double execution
		if elapsedMs <= 1000 {
			waitMs := time.Duration(1000 - elapsedMs)
			time.Sleep(waitMs * time.Millisecond)
		}
		unlock()
	}()

	if err := j.Handler.Run(ctx); err != nil {
		log.Error(ctx, err)
		setLastFailedRun(j.Name, start)
		trx.MustRollback()
	} else {
		setLastSuccessfulRun(j.Name, start)
		trx.MustCommit()
	}
}

func newJobContext() (Context, *dbx.Trx, error) {
	ctx := context.Background()
	ctx = log.WithProperties(ctx, dto.Props{
		log.PropertyKeyContextID: rand.String(32),
		log.PropertyKeyTag:       "JOBS",
	})

	trx, err := dbx.BeginTx(ctx)
	if err != nil {
		log.Error(ctx, err)
		return Context{Context: ctx}, nil, err
	}

	ctx = context.WithValue(ctx, app.TransactionCtxKey, trx)
	return Context{Context: ctx}, trx, nil
}
