package worker

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/storage/postgres"
)

//Job is what's going to be run on background
type Job func(c *Context) error

//Task represents the Name and Job to be run on background
type Task struct {
	Name string
	Job  Job
}

//Context holds references to services available for jobs
type Context struct {
	services *app.Services
	logger   log.Logger
	user     *models.User
	tenant   *models.Tenant
}

//NewContext creates a new context
func NewContext(s *app.Services, l log.Logger) *Context {
	return &Context{
		services: s,
		logger:   l,
	}
}

//SetUser on context
func (c *Context) SetUser(user *models.User) {
	c.user = user
	c.services.SetCurrentUser(user)
}

//SetTenant on context
func (c *Context) SetTenant(tenant *models.Tenant) {
	c.tenant = tenant
	c.services.SetCurrentTenant(tenant)
}

//User from current context
func (c *Context) User() *models.User {
	return c.user
}

//Tenant from current context
func (c *Context) Tenant() *models.Tenant {
	return c.tenant
}

//Services from current context
func (c *Context) Services() *app.Services {
	return c.services
}

//Logger from current context
func (c *Context) Logger() log.Logger {
	return c.logger
}

//Worker is a process that runs tasks
type Worker interface {
	Run(id int)
	Enqueue(task Task)
	Logger() log.Logger
	Use(db *dbx.Database, emailer email.Sender)
}

//BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	db      *dbx.Database
	logger  log.Logger
	emailer email.Sender
	queue   chan Task
}

//New creates a new BackgroundWorker
func New() *BackgroundWorker {
	return &BackgroundWorker{
		logger: log.NewConsoleLogger("BGW"),
		queue:  make(chan Task),
	}
}

//Run initializes the worker loop
func (w *BackgroundWorker) Run(id int) {
	for task := range w.queue {
		trx, err := w.db.Begin()
		if err != nil {
			w.logger.Error(err)
		} else {
			c := NewContext(
				&app.Services{
					Tenants: postgres.NewTenantStorage(trx),
					Users:   postgres.NewUserStorage(trx),
					Ideas:   postgres.NewIdeaStorage(trx),
					Tags:    postgres.NewTagStorage(trx),
					Emailer: w.emailer,
				},
				w.logger,
			)

			start := time.Now()
			w.logger.Infof("Task '%s' started on worker '%s'", log.Magenta(task.Name), log.Magenta(strconv.Itoa(id)))
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

//Enqueue a task on current worker
func (w *BackgroundWorker) Enqueue(task Task) {
	w.queue <- task
}

//Logger from current worker
func (w *BackgroundWorker) Logger() log.Logger {
	return w.logger
}

//Use this to inject worker dependencies
func (w *BackgroundWorker) Use(db *dbx.Database, emailer email.Sender) {
	w.db = db
	w.emailer = emailer
}
