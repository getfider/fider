package worker

import (
	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

//Context holds references to services available for jobs
type Context struct {
	workerID string
	taskName string
	services *app.Services
	logger   log.Logger
	db       *dbx.Database
	baseURL  string
	user     *models.User
	tenant   *models.Tenant
}

//NewContext creates a new context
func NewContext(workerID, taskName string, db *dbx.Database, logger log.Logger) *Context {
	return &Context{
		workerID: workerID,
		taskName: taskName,
		db:       db,
		logger:   logger.New(),
	}
}

//SetBaseURL on context
func (c *Context) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
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

//SetServices on current context
func (c *Context) SetServices(services *app.Services) {
	c.services = services
}

//WorkerID executing current context
func (c *Context) WorkerID() string {
	return c.workerID
}

//TaskName from current context
func (c *Context) TaskName() string {
	return c.taskName
}

//BaseURL from current context
func (c *Context) BaseURL() string {
	return c.baseURL
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

//Database from current context
func (c *Context) Database() *dbx.Database {
	return c.db
}

//Failure logs details of error
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)

	tenant := "undefined"
	if c.Tenant() != nil {
		tenant = fmt.Sprintf("%s (%d)", c.Tenant().Name, c.Tenant().ID)
	}

	user := "not signed in"
	if c.User() != nil {
		user = fmt.Sprintf("%s (%d)", c.User().Name, c.User().ID)
	}

	message := fmt.Sprintf("Task: %s\nTenant: %s\nUser: %s\n%s", c.taskName, tenant, user, err.Error())
	c.logger.Errorf(log.Red(message))
	return err
}
