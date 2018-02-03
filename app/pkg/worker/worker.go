package worker

import (
	"strconv"
	"time"

	"github.com/getfider/fider/app/pkg/log"
)

//MiddlewareFunc is worker middleware
type MiddlewareFunc func(Job) Job

//Job is what's going to be run on background
type Job func(c *Context) error

//Task represents the Name and Job to be run on background
type Task struct {
	Name string
	Job  Job
}

//Worker is a process that runs tasks
type Worker interface {
	Run(id int)
	Enqueue(task Task)
	Logger() log.Logger
	Use(middleware MiddlewareFunc)
}

//BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	logger     log.Logger
	queue      chan Task
	middleware MiddlewareFunc
}

//New creates a new BackgroundWorker
func New() *BackgroundWorker {
	return &BackgroundWorker{
		logger: log.NewConsoleLogger("BGW"),
		queue:  make(chan Task),
		middleware: func(next Job) Job {
			return next
		},
	}
}

//Run initializes the worker loop
func (w *BackgroundWorker) Run(id int) {
	for task := range w.queue {

		c := &Context{
			logger: w.logger,
		}

		start := time.Now()
		c.Logger().Infof("Task '%s' started on worker '%s'", log.Magenta(task.Name), log.Magenta(strconv.Itoa(id)))
		if err := w.middleware(task.Job)(c); err != nil {
			w.logger.Error(err)
		}
		c.Logger().Infof("Task '%s' finished in %s", log.Magenta(task.Name), log.Magenta(time.Since(start).String()))

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
func (w *BackgroundWorker) Use(middleware MiddlewareFunc) {
	w.middleware = middleware
}
