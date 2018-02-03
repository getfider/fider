package worker_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/worker"

	. "github.com/onsi/gomega"
)

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
