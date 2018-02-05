package worker

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/log"
)

//Context holds references to services available for jobs
type Context struct {
	workerID string
	taskName string
	services *app.Services
	logger   log.Logger
	user     *models.User
	tenant   *models.Tenant
}

//NewContext creates a new context
func NewContext(workerID, taskName string, logger log.Logger) *Context {
	return &Context{
		workerID: workerID,
		taskName: taskName,
		logger:   logger,
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
