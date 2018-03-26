package worker

import (
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
	Run(id string)
	Enqueue(task Task)
	Logger() log.Logger
	Use(middleware MiddlewareFunc)
	Length() int
}

//BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	logger     log.Logger
	queue      chan Task
	middleware MiddlewareFunc
}

var maxQueueSize = 100

//New creates a new BackgroundWorker
func New() *BackgroundWorker {
	return &BackgroundWorker{
		logger: log.NewConsoleLogger("BGW"),
		queue:  make(chan Task, maxQueueSize),
		middleware: func(next Job) Job {
			return next
		},
	}
}

//Run initializes the worker loop
func (w *BackgroundWorker) Run(id string) {
	w.logger.Infof("Starting worker %s.", log.Magenta(id))
	for task := range w.queue {

		c := &Context{
			workerID: id,
			taskName: task.Name,
			logger:   w.logger,
		}

		w.middleware(task.Job)(c)
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

//Length from current queue length
func (w *BackgroundWorker) Length() int {
	return len(w.queue)
}

//Use this to inject worker dependencies
func (w *BackgroundWorker) Use(middleware MiddlewareFunc) {
	w.middleware = middleware
}
