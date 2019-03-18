package worker

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
)

//Context holds references to services available for jobs
type Context struct {
	innerCtx      context.Context
	workerID      string
	taskName      string
	services      *app.Services
	baseURL       string
	logoURL       string
	assetsBaseURL string
	user          *models.User
	tenant        *models.Tenant
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

//SetBaseURL on context
func (c *Context) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

//SetLogoURL on context
func (c *Context) SetLogoURL(logoURL string) {
	c.logoURL = logoURL
}

//SetAssetsBaseURL on context
func (c *Context) SetAssetsBaseURL(assetsBaseURL string) {
	c.assetsBaseURL = assetsBaseURL
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

//SetServices on current context
func (c *Context) SetServices(services *app.Services) {
	c.services = services
}

//WorkerID executing current context
func (c *Context) WorkerID() string {
	return c.workerID
}

func (c *Context) Dispatch(m bus.Msg) error {
	return bus.Dispatch(c.innerCtx, m)
}

func (c *Context) Publish(evt bus.Event) {
	bus.Publish(c.innerCtx, evt)
}

//TaskName from current context
func (c *Context) TaskName() string {
	return c.taskName
}

//BaseURL from current context
func (c Context) BaseURL() string {
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

//Failure logs details of error
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	log.Error(c.innerCtx, err)
	return err
}

// LogoURL return the full URL to the tenant-specific logo URL
func (c Context) LogoURL() string {
	return c.logoURL
}

// TenantAssetsURL return the full URL to a tenant-specific static asset
func (c Context) TenantAssetsURL(path string, a ...interface{}) string {
	path = fmt.Sprintf(path, a...)
	return c.assetsBaseURL + path
}

func (ctx *Context) Debug(message string) {
	log.Debug(ctx.innerCtx, message)
}

func (ctx *Context) Debugf(message string, props log.Props) {
	log.Debugf(ctx.innerCtx, message, props)
}

func (ctx *Context) Info(message string) {
	log.Info(ctx.innerCtx, message)
}

func (ctx *Context) Infof(message string, props log.Props) {
	log.Infof(ctx.innerCtx, message, props)
}

func (ctx *Context) Warn(message string) {
	log.Warn(ctx.innerCtx, message)
}

func (ctx *Context) Warnf(message string, props log.Props) {
	log.Warnf(ctx.innerCtx, message, props)
}

func (ctx *Context) Error(err error) {
	log.Error(ctx.innerCtx, err)
}

func (ctx *Context) Errorf(message string, props log.Props) {
	log.Errorf(ctx.innerCtx, message, props)
}
