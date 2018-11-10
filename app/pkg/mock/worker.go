package mock

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/log/noop"
	"github.com/getfider/fider/app/pkg/worker"
)

// Worker is fake wrapper for background worker
type Worker struct {
	tenant   *models.Tenant
	user     *models.User
	services *app.Services
	baseURL  string
}

func createWorker(services *app.Services) *Worker {
	return &Worker{
		services: services,
	}
}

// OnTenant set current context tenant
func (w *Worker) OnTenant(tenant *models.Tenant) *Worker {
	w.tenant = tenant
	return w
}

// AsUser set current context user
func (w *Worker) AsUser(user *models.User) *Worker {
	w.user = user
	return w
}

// WithBaseURL set current context baseURL
func (w *Worker) WithBaseURL(baseURL string) *Worker {
	w.baseURL = baseURL
	return w
}

// Execute given task with current context
func (w *Worker) Execute(task worker.Task) error {
	context := worker.NewContext("0", task, nil, noop.NewLogger())
	context.SetServices(w.services)
	context.SetUser(w.user)
	context.SetTenant(w.tenant)
	context.SetBaseURL(w.baseURL)
	return task.Job(context)
}

// NewNoopTask returns a worker task that does nothing
func NewNoopTask() worker.Task {
	return worker.Task{
		Name: "Noop Task",
		Job: func(c *worker.Context) error {
			return nil
		},
	}
}
