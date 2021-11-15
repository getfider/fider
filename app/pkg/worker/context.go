package worker

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
)

//Context holds references to services available for jobs
type Context struct {
	context.Context
	workerID string
	taskName string
}

//NewContext creates a new context
func NewContext(ctx context.Context, workerID string, task Task) *Context {
	ctx = log.WithProperty(ctx, log.PropertyKeyContextID, rand.String(32))

	if task.OriginContext != nil {
		ctx = context.WithValue(ctx, app.RequestCtxKey, task.OriginContext.Value(app.RequestCtxKey))
		ctx = context.WithValue(ctx, app.TenantCtxKey, task.OriginContext.Value(app.TenantCtxKey))
		ctx = context.WithValue(ctx, app.LocaleCtxKey, task.OriginContext.Value(app.LocaleCtxKey))
		ctx = context.WithValue(ctx, app.UserCtxKey, task.OriginContext.Value(app.UserCtxKey))

		ctx = log.WithProperties(ctx, dto.Props{
			log.PropertyKeySessionID: log.GetProperty(task.OriginContext, log.PropertyKeySessionID),
			log.PropertyKeyUserID:    log.GetProperty(task.OriginContext, log.PropertyKeyUserID),
			log.PropertyKeyTenantID:  log.GetProperty(task.OriginContext, log.PropertyKeyTenantID),
		})
	}

	return &Context{
		Context:  ctx,
		workerID: workerID,
		taskName: task.Name,
	}
}

//WorkerID executing current context
func (c *Context) WorkerID() string {
	return c.workerID
}

//TaskName from current context
func (c *Context) TaskName() string {
	return c.taskName
}

// Set saves data in the context.
func (c *Context) Set(key interface{}, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

//User from current context
func (c *Context) User() *entity.User {
	user, ok := c.Value(app.UserCtxKey).(*entity.User)
	if ok {
		return user
	}
	return nil
}

//Tenant from current context
func (c *Context) Tenant() *entity.Tenant {
	tenant, ok := c.Value(app.TenantCtxKey).(*entity.Tenant)
	if ok {
		return tenant
	}
	return nil
}

//Failure logs details of error
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	log.Error(c, err)
	return err
}
