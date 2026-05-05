package jobs_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/jobs"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestEmailSuppressionJob_Schedule_IsCorrect(t *testing.T) {
	RegisterT(t)

	job := &jobs.EmailSuppressionJobHandler{}
	Expect(job.Schedule()).Equals("0 5 * * * *")
}

func TestEmailSuppressionJob_ShouldSuppressRecentFailures(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.FetchRecentSuppressions) error {
		q.EmailAddresses = []string{
			"test1@gmail.com", "test2@gmail.com",
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SuppressEmail) error {
		Expect(c.EmailAddresses).Equals([]string{"test1@gmail.com", "test2@gmail.com"})
		return nil
	})

	job := &jobs.EmailSuppressionJobHandler{}
	err := job.Run(jobs.Context{
		Context: context.Background(),
	})
	Expect(err).IsNil()
}
