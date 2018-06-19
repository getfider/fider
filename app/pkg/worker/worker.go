package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/log/console"
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
	sync.RWMutex
}

var maxQueueSize = 100

//New creates a new BackgroundWorker
func New() *BackgroundWorker {
	return &BackgroundWorker{
		logger: console.NewLogger("BGW"),
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
		w.Lock()
		w.len = w.len - 1
		w.Unlock()
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
	w.Lock()
	w.len = w.len + 1
	w.Unlock()
	w.queue <- task
}

//Logger from current worker
func (w *BackgroundWorker) Logger() log.Logger {
	return w.logger
}

//Length from current queue length
func (w *BackgroundWorker) Length() int64 {
	w.RLock()
	defer w.RUnlock()
	return w.len
}

//Use this to inject worker dependencies
func (w *BackgroundWorker) Use(middleware MiddlewareFunc) {
	w.middleware = middleware
}
