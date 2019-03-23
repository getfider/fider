package worker

import (
	"context"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
)

//Context holds references to services available for jobs
type Context struct {
	innerCtx context.Context
	workerID string
	taskName string
	services *app.Services
	user     *models.User
	tenant   *models.Tenant
}

//NewContext creates a new context
func NewContext(ctx context.Context, workerID string, task Task) *Context {
	contextID := rand.String(32)

	ctx = log.SetProperty(ctx, log.PropertyKeyContextID, contextID)
	if task.OriginSessionID != "" {
		ctx = log.SetProperty(ctx, log.PropertyKeySessionID, task.OriginSessionID)
	}

	return &Context{
		innerCtx: ctx,
		workerID: workerID,
		taskName: task.Name,
	}
}

//Database returns current database
func (c *Context) Database() *dbx.Database {
	return c.innerCtx.Value(app.DatabaseCtxKey).(*dbx.Database)
}

//SetUser on context
func (c *Context) SetUser(user *models.User) {
	c.user = user
	if user != nil {
		c.innerCtx = context.WithValue(c.innerCtx, app.UserCtxKey, user)
		c.innerCtx = log.SetProperty(c.innerCtx, log.PropertyKeyUserID, user.ID)
	}
	if c.services != nil {
		c.services.SetCurrentUser(user)
	}
}

//SetTenant on context
func (c *Context) SetTenant(tenant *models.Tenant) {
	c.tenant = tenant
	if tenant != nil {
		c.innerCtx = context.WithValue(c.innerCtx, app.TenantCtxKey, tenant)
		c.innerCtx = log.SetProperty(c.innerCtx, log.PropertyKeyTenantID, tenant.ID)
	}
	if c.services != nil {
		c.services.SetCurrentTenant(tenant)
	}
}

//Set value on current context based on given key
func (c *Context) Set(key string, value interface{}) {
	c.innerCtx = context.WithValue(c.innerCtx, key, value)
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

//Failure logs details of error
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	log.Error(c.innerCtx, err)
	return err
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.innerCtx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.innerCtx.Done()
}

func (c *Context) Err() error {
	return c.innerCtx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.innerCtx.Value(key)
}
