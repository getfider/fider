package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/rand"
)

//MiddlewareFunc is worker middleware
type MiddlewareFunc func(Job) Job

//Job is what's going to be run on background
type Job func(c *Context) error

//Task represents the Name and Job to be run on background
type Task struct {
	OriginSessionID string
	Name            string
	Job             Job
}

//Worker is a process that runs tasks
type Worker interface {
	Run(id string)
	Enqueue(task Task)
	Use(middleware MiddlewareFunc)
	Length() int64
	Shutdown(ctx context.Context) error
}

//BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	innerCtx   context.Context
	db         *dbx.Database
	queue      chan Task
	len        int64
	middleware MiddlewareFunc
	sync.RWMutex
}

var maxQueueSize = 100

//New creates a new BackgroundWorker
func New(db *dbx.Database) *BackgroundWorker {
	ctx := context.Background()
	ctx = context.WithValue(ctx, app.DatabaseCtxKey, db)
	ctx = log.SetProperty(ctx, log.PropertyKeyContextID, rand.String(32))
	ctx = log.SetProperty(ctx, log.PropertyKeyTag, "BGW")

	return &BackgroundWorker{
		innerCtx: ctx,
		db:       db,
		queue:    make(chan Task, maxQueueSize),
		middleware: func(next Job) Job {
			return next
		},
	}
}

//Run initializes the worker loop
func (w *BackgroundWorker) Run(workerID string) {
	log.Infof(w.innerCtx, "Starting worker @{WorkerID:magenta}.", dto.Props{
		"WorkerID": workerID,
	})
	for task := range w.queue {
		c := NewContext(w.innerCtx, workerID, task)

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

			log.Infof(w.innerCtx, "Waiting for work queue: @{Count}", dto.Props{
				"Count": count,
			})

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
