package worker

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/storage/postgres"
)

type Job func(c *Context) error

type Task struct {
	Name string
	Job  Job
}

type Context struct {
	services *app.Services
	logger   log.Logger
	user     *models.User
	tenant   *models.Tenant
}

func (c *Context) SetUser(user *models.User) {
	c.user = user
	c.services.SetCurrentUser(user)
}

func (c *Context) SetTenant(tenant *models.Tenant) {
	c.tenant = tenant
	c.services.SetCurrentTenant(tenant)
}

func (c *Context) User() *models.User {
	return c.user
}

func (c *Context) Tenant() *models.Tenant {
	return c.tenant
}

func (c *Context) Services() *app.Services {
	return c.services
}

func (c *Context) Logger() log.Logger {
	return c.logger
}

type Worker interface {
	Run(id int)
	Enqueue(task Task)
	Logger() log.Logger
	Use(db *dbx.Database, emailer email.Sender)
}

type BackgroundWorker struct {
	db      *dbx.Database
	logger  log.Logger
	emailer email.Sender
	queue   chan Task
}

func New() *BackgroundWorker {
	return &BackgroundWorker{
		logger: log.NewConsoleLogger("BGW"),
		queue:  make(chan Task),
	}
}

func (w *BackgroundWorker) Run(id int) {
	for task := range w.queue {
		trx, err := w.db.Begin()
		if err != nil {
			w.logger.Error(err)
		} else {
			c := &Context{
				services: &app.Services{
					Tenants: postgres.NewTenantStorage(trx),
					Users:   postgres.NewUserStorage(trx),
					Ideas:   postgres.NewIdeaStorage(trx),
					Tags:    postgres.NewTagStorage(trx),
					Emailer: w.emailer,
				},
				logger: w.logger,
			}

			start := time.Now()
			w.logger.Infof("Task '%s' started on worker '%d'", log.Magenta(task.Name), log.Magenta(string(id)))
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("%v\n%v", r, string(debug.Stack()))

					c.Logger().Error(err)
					w.logger.Infof("Task '%s' finished in %s", log.Magenta(task.Name), log.Magenta(time.Since(start).String()))
					if trx != nil {
						trx.Rollback()
					}
				}
			}()
			if err := task.Job(c); err != nil {
				panic(err)
			}
			w.logger.Infof("Task '%s' finished in %s", log.Magenta(task.Name), log.Magenta(time.Since(start).String()))
			trx.Commit()
		}
	}
}

func (w *BackgroundWorker) Enqueue(task Task) {
	w.queue <- task
}

func (w *BackgroundWorker) Logger() log.Logger {
	return w.logger
}

func (w *BackgroundWorker) Use(db *dbx.Database, emailer email.Sender) {
	w.db = db
	w.emailer = emailer
}
