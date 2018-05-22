package middlewares

import (
	"fmt"
	"time"

	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/worker"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/postgres"
)

// Noop does nothing
func Noop() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			return next(c)
		}
	}
}

//WorkerSetup current context with some services
func WorkerSetup(logger log.Logger) worker.MiddlewareFunc {
	db := dbx.NewWithLogger(logger)
	emailer := app.NewEmailer(logger)
	return func(next worker.Job) worker.Job {
		return func(c *worker.Context) (err error) {
			start := time.Now()
			c.Logger().Debugf("Task '%s' started on worker '%s'", log.Magenta(c.TaskName()), log.Magenta(c.WorkerID()))

			trx, err := db.Begin()
			if err != nil {
				err = c.Failure(err)
				c.Logger().Debugf("Task '%s' finished in %s", log.Magenta(c.TaskName()), log.Magenta(time.Since(start).String()))
				return err
			}

			c.SetServices(&app.Services{
				Tenants:       postgres.NewTenantStorage(trx),
				Users:         postgres.NewUserStorage(trx),
				Ideas:         postgres.NewIdeaStorage(trx),
				Tags:          postgres.NewTagStorage(trx),
				Notifications: postgres.NewNotificationStorage(trx),
				Emailer:       emailer,
			})

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					c.Failure(err)
					trx.Rollback()
					c.Logger().Debugf("Task '%s' panicked in %s (rolled back)", log.Magenta(c.TaskName()), log.Magenta(time.Since(start).String()))
				}
			}()

			//Execute the chain
			if err = next(c); err != nil {
				trx.Rollback()
				c.Logger().Debugf("Task '%s' finished in %s (rolled back)", log.Magenta(c.TaskName()), log.Magenta(time.Since(start).String()))
				return err
			}

			//No errors, so try to commit it
			if err = trx.Commit(); err != nil {
				c.Logger().Errorf("Failed to commit request with: %s", err.Error())
				c.Logger().Debugf("Task '%s' finished in %s (rolled back)", log.Magenta(c.TaskName()), log.Magenta(time.Since(start).String()))
				return err
			}

			//Still no errors, everything is fine!
			c.Logger().Debugf("Task '%s' finished in %s (committed)", log.Magenta(c.TaskName()), log.Magenta(time.Since(start).String()))
			return nil
		}
	}
}

//WebSetup current context with some services
func WebSetup(logger log.Logger) web.MiddlewareFunc {
	db := dbx.NewWithLogger(logger)
	db.Migrate()
	emailer := app.NewEmailer(logger)
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			path := log.Magenta(c.Request.Method + " " + c.Request.URL.String())

			start := time.Now()
			c.Logger().Debugf("%s started", path)

			trx, err := db.Begin()
			if err != nil {
				err = c.Failure(err)
				c.Logger().Debugf("%s finished in %s", path, log.Magenta(time.Since(start).String()))
				return err
			}

			c.SetActiveTransaction(trx)
			c.SetServices(&app.Services{
				Tenants:       postgres.NewTenantStorage(trx),
				OAuth:         &oauth.HTTPService{},
				Users:         postgres.NewUserStorage(trx),
				Ideas:         postgres.NewIdeaStorage(trx),
				Tags:          postgres.NewTagStorage(trx),
				Notifications: postgres.NewNotificationStorage(trx),
				Emailer:       emailer,
			})

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					c.Failure(err)
					c.Rollback()
					c.Logger().Debugf("%s panicked in %s (rolled back)", path, log.Magenta(time.Since(start).String()))
				}
			}()

			//Execute the chain
			if err := next(c); err != nil {
				c.Rollback()
				c.Logger().Debugf("%s finished in %s (rolled back)", path, log.Magenta(time.Since(start).String()))
				return err
			}

			//No errors, so try to commit it
			if err := c.Commit(); err != nil {
				c.Logger().Errorf("Failed to commit request with: %s", err.Error())
				c.Logger().Debugf("%s finished in %s (rolled back)", path, log.Magenta(time.Since(start).String()))
				return err
			}

			//Still no errors, everything is fine!
			c.Logger().Debugf("%s finished in %s (committed)", path, log.Magenta(time.Since(start).String()))
			return nil
		}
	}
}
