package jobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app/jobs"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestTrialReminderJob_Schedule_IsCorrect(t *testing.T) {
	RegisterT(t)

	job := &jobs.TrialReminderJobHandler{}
	Expect(job.Schedule()).Equals("0 0 7 * * *")
}

func TestTrialReminderJob_DoNothingOnFirstRun(t *testing.T) {
	RegisterT(t)

	job := &jobs.TrialReminderJobHandler{}
	err := job.Run(jobs.Context{
		Context: context.Background(),
	})
	Expect(err).IsNil()
}

func TestTrialReminderJob_ShouldSendEmailsToTrialingContacts(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTrialingTenantContacts) error {
		q.Contacts = []*entity.User{
			{Name: "user1", Email: "user1@gmail.com", Tenant: &entity.Tenant{Subdomain: "demo1"}},
			{Name: "user2", Email: "user2@gmail.com", Tenant: &entity.Tenant{Subdomain: "demo2"}},
		}
		return nil
	})

	bus.AddListener(func(ctx context.Context, c *cmd.SendMail) error {
		Expect(c.TemplateName).Equals("trial_1day")
		Expect(c.To).Equals([]dto.Recipient{
			{Name: "user1", Address: "user1@gmail.com", Props: dto.Props{"name": "user1", "url": "https://demo1.test.fider.io"}},
			{Name: "user2", Address: "user2@gmail.com", Props: dto.Props{"name": "user2", "url": "https://demo2.test.fider.io"}},
		})
		return nil
	})

	lastRun := time.Now()
	job := &jobs.TrialReminderJobHandler{
		Days:         1,
		TemplateName: "trial_1day",
	}
	err := job.Run(jobs.Context{
		Context:           context.Background(),
		LastSuccessfulRun: &lastRun,
	})
	Expect(err).IsNil()
}
