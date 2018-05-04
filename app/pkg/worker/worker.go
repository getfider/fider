package worker

import (
	"context"
	"errors"
	"sync/atomic"
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
	Run(id string)
	Enqueue(task Task)
	Logger() log.Logger
	Use(middleware MiddlewareFunc)
	Length() int64
	Shutdown(ctx context.Context) error
}

//BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	logger     log.Logger
	queue      chan Task
	len        int64
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
		atomic.AddInt64(&w.len, -1)
	}
}

//Shutdown current worker
func (w *BackgroundWorker) Shutdown(ctx context.Context) error {
	if w.Length() > 0 {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			count := w.Length()
			if count == 0 {
				return nil
			}
			w.logger.Infof("waiting for work queue: %d", count)
			select {
			case <-ctx.Done():
				return errors.New("timeout waiting for worker queue")
			case <-ticker.C:
			}
		}
	}
	return nil
}

//Enqueue a task on current worker
func (w *BackgroundWorker) Enqueue(task Task) {
	atomic.AddInt64(&w.len, 1)
	w.queue <- task
}

//Logger from current worker
func (w *BackgroundWorker) Logger() log.Logger {
	return w.logger
}

//Length from current queue length
func (w *BackgroundWorker) Length() int64 {
	return w.len
}

//Use this to inject worker dependencies
func (w *BackgroundWorker) Use(middleware MiddlewareFunc) {
	w.middleware = middleware
}
