package middlewares

import (
	"fmt"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/email/noop"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/web/util"
	"github.com/getfider/fider/app/pkg/worker"
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
func WorkerSetup() worker.MiddlewareFunc {
	return func(next worker.Job) worker.Job {
		return func(c *worker.Context) (err error) {
			start := time.Now()
			c.Logger().Debugf("Task '@{TaskName:magenta}' started on worker '@{WorkerID:magenta}'", log.Props{
				"TaskName": c.TaskName(),
				"WorkerID": c.WorkerID(),
			})

			trx, err := c.Database().Begin()
			if err != nil {
				err = c.Failure(err)
				c.Logger().Debugf("Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms", log.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			trx.SetLogger(c.Logger())

			//In case it panics somewhere
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					c.Failure(err)
					trx.Rollback()
					c.Logger().Debugf("Task '@{TaskName:magenta}' panicked in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
						"TaskName":  c.TaskName(),
						"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
					})
				}
			}()

			c.SetServices(&app.Services{
				Tenants:       postgres.NewTenantStorage(trx),
				Users:         postgres.NewUserStorage(trx),
				Posts:         postgres.NewPostStorage(trx),
				Tags:          postgres.NewTagStorage(trx),
				Notifications: postgres.NewNotificationStorage(trx),
				Emailer:       newEmailer(c.Logger()),
			})

			//Execute the chain
			if err = next(c); err != nil {
				trx.Rollback()
				c.Logger().Debugf("Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//No errors, so try to commit it
			if err = trx.Commit(); err != nil {
				c.Logger().Errorf("Failed to commit request with: @{Error}", log.Props{
					"Error": err.Error(),
				})
				c.Logger().Debugf("Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
					"TaskName":  c.TaskName(),
					"ElapsedMs": time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//Still no errors, everything is fine!
			c.Logger().Debugf("Task '@{TaskName:magenta}' finished in @{ElapsedMs:magenta}ms (committed)", log.Props{
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
		return func(c web.Context) error {
			start := time.Now()
			c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} started for @{ClientIP:magenta}", log.Props{
				"HttpMethod": c.Request.Method,
				"URL":        c.Request.URL.String(),
				"ClientIP":   c.Request.ClientIP,
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
					c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} panicked in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
						"HttpMethod": c.Request.Method,
						"URL":        c.Request.URL.String(),
						"ElapsedMs":  time.Since(start).Nanoseconds() / int64(time.Millisecond),
					})
				}
			}()

			trx, err := c.Engine().Database().Begin()
			if err != nil {
				err = c.Failure(err)
				c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms", log.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			trx.SetLogger(c.Logger())

			oauthBaseURL := webutil.GetOAuthBaseURL(c)
			tenantStorage := postgres.NewTenantStorage(trx)

			c.SetActiveTransaction(trx)
			c.SetServices(&app.Services{
				Tenants:       tenantStorage,
				OAuth:         web.NewOAuthService(oauthBaseURL, tenantStorage),
				Users:         postgres.NewUserStorage(trx),
				Posts:         postgres.NewPostStorage(trx),
				Tags:          postgres.NewTagStorage(trx),
				Notifications: postgres.NewNotificationStorage(trx),
				Emailer:       newEmailer(c.Logger()),
			})

			//Execute the chain
			if err := next(c); err != nil {
				c.Rollback()
				c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//No errors, so try to commit it
			if err := c.Commit(); err != nil {
				c.Logger().Errorf("Failed to commit request with: @{Error}", log.Props{
					"Error": err.Error(),
				})
				c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (rolled back)", log.Props{
					"HttpMethod": c.Request.Method,
					"URL":        c.Request.URL.String(),
					"ElapsedMs":  time.Since(start).Nanoseconds() / int64(time.Millisecond),
				})
				return err
			}

			//Still no errors, everything is fine!
			c.Logger().Infof("@{HttpMethod:magenta} @{URL:magenta} finished in @{ElapsedMs:magenta}ms (committed)", log.Props{
				"HttpMethod": c.Request.Method,
				"URL":        c.Request.URL.String(),
				"ElapsedMs":  time.Since(start).Nanoseconds() / int64(time.Millisecond),
			})
			return nil
		}
	}
}

func newEmailer(logger log.Logger) email.Sender {
	if env.IsTest() {
		return noop.NewSender()
	}
	if env.IsDefined("EMAIL_MAILGUN_API") {
		return mailgun.NewSender(logger, web.NewHTTPClient(), env.MustGet("EMAIL_MAILGUN_DOMAIN"), env.MustGet("EMAIL_MAILGUN_API"))
	}
	return smtp.NewSender(
		logger,
		env.MustGet("EMAIL_SMTP_HOST"),
		env.MustGet("EMAIL_SMTP_PORT"),
		env.GetEnvOrDefault("EMAIL_SMTP_USERNAME", ""),
		env.GetEnvOrDefault("EMAIL_SMTP_PASSWORD", ""),
	)
}
