package middlewares

import (
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
			log.Debugf(c, "Task '@{TaskName:magenta}' started on worker '@{WorkerID:magenta}'", dto.Props{
				"TaskName": c.TaskName(),
				"WorkerID": c.WorkerID(),
			})

			trx, err := dbx.BeginTx(c)
			if err != nil {
				err = c.Failure(err)
				log.Debugf(c, "Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms", dto.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					_ = c.Failure(errors.Panicked(r))
					trx.MustRollback()
					log.Debugf(c, "Task '@{TaskName:magenta}' panicked in @{ElapsedMs:magenta}ms (rolled back)", dto.Props{
						"TaskName":  c.TaskName(),
						"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
					})
				}
			}()

			c.Set(app.TransactionCtxKey, trx)

			//Execute the chain
			if err = next(c); err != nil {
				trx.MustRollback()
				log.Debugf(c, "Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (rolled back)", dto.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//No errors, so try to commit it
			if err = trx.Commit(); err != nil {
				log.Errorf(c, "Failed to commit request with: @{Error}", dto.Props{
					"Error": err.Error(),
				})
				log.Debugf(c, "Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (rolled back)", dto.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//Still no errors, everything is fine!
			log.Debugf(c, "Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (committed)", dto.Props{
				"TaskName":  c.TaskName(),
				"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
			return nil
		}
	}
}

//WebSetup current context with some services
func WebSetup() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} started for @{ClientIP:magenta}", dto.Props{
				"HttpMethod": c.Request.Method,
				"URL":        c.Request.URL.String(),
				"ClientIP":   c.Request.ClientIP,
			})

			trx, err := dbx.BeginTx(c)
			if err != nil {
				err = c.Failure(err)
				log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms", dto.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(c.Request.StartTime).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			c.Set(app.TransactionCtxKey, trx)

			//Execute the chain
			if err := next(c); err != nil {
				c.Rollback()
				log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (rolled back)", dto.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(c.Request.StartTime).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//No errors, so try to commit it
			if err := c.Commit(); err != nil {
				log.Errorf(c, "Failed to commit request with: @{Error}", dto.Props{
					"Error": err.Error(),
				})
				log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (rolled back)", dto.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(c.Request.StartTime).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//Still no errors, everything is fine!
			log.Infof(c, "@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (committed)", dto.Props{
				"HttpMethod": c.Request.Method,
				"URL":        c.Request.URL.String(),
				"ElapsedMs":  time.Since(c.Request.StartTime).Nanoseconds() / int64(time.Millisecond),
			})
			return nil
		}
	}
}
