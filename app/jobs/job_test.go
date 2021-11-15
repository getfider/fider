package jobs_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/getfider/fider/app/jobs"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

type MockJobHandler struct {
	WaitTime   time.Duration
	ShouldFail bool
}

func (e MockJobHandler) Schedule() string {
	return "0 * * * * *"
}

func (e MockJobHandler) Run(ctx jobs.Context) error {
	if e.WaitTime > 0 {
		time.Sleep(e.WaitTime)
	}

	if e.ShouldFail {
		return errors.New("Failed")
	}
	return nil
}

func TestJob_WhenSuccessful_ShouldUpdateLastSuccessfulRun(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetSystemSettings) error {
		Expect(q.Key).Equals("jobs.Test.last_successful_run")
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetSystemSettings) error {
		Expect(c.Key).Equals("jobs.Test.last_successful_run")
		Expect(c.Value).ContainsSubstring(time.Now().Format("2006-01-02"))
		return nil
	})

	schedule, job := jobs.NewJob(context.Background(), "Test", MockJobHandler{
		ShouldFail: false,
	})
	Expect(schedule).Equals("0 * * * * *")

	job.Run()
}

func TestJob_WhenError_ShouldUpdateLastFailedRun(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetSystemSettings) error {
		Expect(q.Key).Equals("jobs.Test.last_successful_run")
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetSystemSettings) error {
		Expect(c.Key).Equals("jobs.Test.last_failed_run")
		Expect(c.Value).ContainsSubstring(time.Now().Format("2006-01-02"))
		return nil
	})

	schedule, job := jobs.NewJob(context.Background(), "Test", MockJobHandler{
		ShouldFail: true,
	})
	Expect(schedule).Equals("0 * * * * *")

	job.Run()
}

func TestJob_WhenTwoConcurrentRuns_ShouldExecuteOnce(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetSystemSettings) error {
		return nil
	})

	counter := 0
	bus.AddHandler(func(ctx context.Context, c *cmd.SetSystemSettings) error {
		counter++
		return nil
	})

	_, job1 := jobs.NewJob(context.Background(), "Test", MockJobHandler{
		WaitTime: 1 * time.Second,
	})
	_, job2 := jobs.NewJob(context.Background(), "Test", MockJobHandler{
		WaitTime: 1 * time.Second,
	})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		job1.Run()
	}()
	go func() {
		defer wg.Done()
		job2.Run()
	}()

	wg.Wait()

	Expect(counter).Equals(1)
}
