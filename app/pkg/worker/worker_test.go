package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app/pkg/worker"

	. "github.com/onsi/gomega"
)

var dummyTask = worker.Task{
	Name: "Do Something",
	Job: func(c *worker.Context) error {
		return nil
	},
}

func TestBackgroundWorker(t *testing.T) {
	RegisterTestingT(t)

	var finished bool

	w := worker.New()
	w.Enqueue(worker.Task{
		Name: "Do Something",
		Job: func(c *worker.Context) error {
			finished = true
			return nil
		},
	})
	Expect(w.Length()).To(Equal(1))
	go w.Run("worker-1")
	Eventually(func() bool {
		return finished
	}).Should(BeTrue())
}

func TestBackgroundWorker_ShutdownWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	w := worker.New()
	Expect(w.Shutdown(ctx)).To(BeNil())
}

func TestBackgroundWorker_ShutdownWithStuckTasks(t *testing.T) {
	RegisterTestingT(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	w := worker.New()
	w.Enqueue(dummyTask)
	Expect(w.Shutdown(ctx)).NotTo(BeNil())
}

func TestBackgroundWorker_ShutdownWithRunningTasks(t *testing.T) {
	RegisterTestingT(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	w := worker.New()
	w.Enqueue(dummyTask)
	go w.Run("worker-1")
	Expect(w.Shutdown(ctx)).To(BeNil())
}
