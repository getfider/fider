package mock

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/worker"
)

// Worker is fake wrapper for background worker
type Worker struct {
	context *worker.Context
}

func createWorker(services *app.Services) *Worker {
	return &Worker{
		context: worker.NewContext(services, log.NewNoopLogger()),
	}
}

// OnTenant set current context tenant
func (w *Worker) OnTenant(tenant *models.Tenant) *Worker {
	w.context.SetTenant(tenant)
	return w
}

// AsUser set current context user
func (w *Worker) AsUser(user *models.User) *Worker {
	w.context.SetUser(user)
	return w
}

// Execute given task with current context
func (w *Worker) Execute(task worker.Task) error {
	return task.Job(w.context)
}
