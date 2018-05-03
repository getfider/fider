package worker_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/worker"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestBackgroundWorker(t *testing.T) {
	RegisterT(t)

	var finished bool

	w := worker.New()
	w.Enqueue(worker.Task{
		Name: "Do Something",
		Job: func(c *worker.Context) error {
			finished = true
			return nil
		},
	})
	Expect(w.Length()).Equals(1)
	go w.Run("worker-1")
	Expect(func() bool {
		return finished
	}).EventuallyEquals(true)
}
