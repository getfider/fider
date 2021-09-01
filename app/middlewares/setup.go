package middlewares

import (
	"context"
	"time"

	"github.com/getfider/fider/app/pkg/dbx"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

//WorkerSetup current context with some services
func WorkerSetup() worker.MiddlewareFunc {
	return func(next worker.Job) worker.Job {
		return func(c *worker.Context) (err error) {
			start := time.Now()
			log.Infof(c, "Task '@{TaskName:magenta}' started on worker '@{WorkerID:magenta}'", dto.Props{
				"TaskName": c.TaskName(),
				"WorkerID": c.WorkerID(),
			})

			logFinish := func(state string, err error) {
				elapsedMs := time.Since(start).Nanoseconds() / int64(time.Millisecond)

				if errors.Cause(err) == context.Canceled {
					log.Infof(c, "Task '@{TaskName:magenta}' was canceled after @{ElapsedMs:magenta}ms", dto.Props{
						"State":     state,
						"TaskName":  c.TaskName(),
						"ElapsedMs": elapsedMs,
					})
					return
				}

				if err != nil {
					log.Error(c, err)
				}
				log.Infof(c, "Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (@{State})", dto.Props{
					"State":     state,
					"TaskName":  c.TaskName(),
					"ElapsedMs": elapsedMs,
				})
			}

			trx, err := dbx.BeginTx(c)
			if err != nil {
				err = c.Failure(err)
				logFinish("begin_error", err)
				return err
			}

			c.Set(app.TransactionCtxKey, trx)

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					err := c.Failure(errors.Panicked(r))
					trx.MustRollback()
					logFinish("panicked", err)
				}
			}()

			//Execute the chain
			if err = next(c); err != nil {
				trx.MustRollback()
				logFinish("next_error", err)
				return err
			}

			//No errors, so try to commit it
			if err = trx.Commit(); err != nil {
				logFinish("commit_error", err)
				return err
			}

			//Still no errors, everything is fine!
			logFinish("committed", nil)
			return nil
		}
	}
}

//WebSetup current context with some services
func WebSetup() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} started", dto.Props{
				"HttpMethod": c.Request.Method,
				"URL":        c.Request.URL.String(),
			})

			logFinish := func(state string, err error) {
				elapsedMs := time.Since(c.Request.StartTime).Nanoseconds() / int64(time.Millisecond)

				if errors.Cause(err) == context.Canceled {
					log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} was canceled after @{ElapsedMs:magenta}ms", dto.Props{
						"HttpMethod": c.Request.Method,
						"URL":        c.Request.URL.String(),
						"ElapsedMs":  elapsedMs,
					})
					return
				}

				if err != nil {
					log.Error(c, err)
				}

				log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} finished with @{StatusCode:magenta} in @{ElapsedMs:magenta}ms (@{State})", dto.Props{
					"State":      state,
					"HttpMethod": c.Request.Method,
					"StatusCode": c.Response.StatusCode,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  elapsedMs,
				})
			}

			trx, err := dbx.BeginTx(c)
			if err != nil {
				err = c.Failure(err)
				logFinish("begin_error", err)
				return err
			}

			c.Set(app.TransactionCtxKey, trx)

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					err := c.Failure(errors.Panicked(r))
					trx.MustRollback()
					logFinish("panicked", err)
				}
			}()

			//Execute the chain
			if err := next(c); err != nil {
				c.Rollback()
				logFinish("next_error", err)
				return err
			}

			//No errors, so try to commit it
			if err := c.Commit(); err != nil {
				logFinish("commit_error", err)
				return err
			}

			//Still no errors, everything is fine!
			logFinish("committed", nil)
			return nil
		}
	}
}
