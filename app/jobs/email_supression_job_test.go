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

func TestEmailSupressionJob_Schedule_IsCorrect(t *testing.T) {
	RegisterT(t)

	job := &jobs.EmailSupressionJobHandler{}
	Expect(job.Schedule()).Equals("0 5 * * * *")
}

func TestEmailSupressionJob_ShouldSupressRecentFailures(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.FetchRecentSupressions) error {
		q.EmailAddresses = []string{
			"test1@gmail.com", "test2@gmail.com",
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SupressEmail) error {
		Expect(c.EmailAddresses).Equals([]string{"test1@gmail.com", "test2@gmail.com"})
		return nil
	})

	job := &jobs.EmailSupressionJobHandler{}
	err := job.Run(jobs.Context{
		Context: context.Background(),
	})
	Expect(err).IsNil()
}
