package tasks_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/worker"
	"github.com/getfider/fider/app/services/email/emailmock"
	"github.com/getfider/fider/app/tasks"
)

var pendingAuthor = &entity.User{ID: 99, Name: "New Visitor", Email: "new@got.com", Role: enum.RoleVisitor}

// moderationRun executes the task with an admin (JonSnow) and a visitor (AryaStark) subscribed to everything,
// and returns who got web notifications, who got emails, and whether a webhook fired
func moderationRun(task worker.Task) (notified []*entity.User, emailed []dto.Recipient, webhook bool) {
	bus.Init(emailmock.Service{})
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		notified = append(notified, c.User)
		return nil
	})
	bus.AddHandler(func(ctx context.Context, c *cmd.AddMentionNotification) error { return nil })
	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error { return nil })
	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{mock.JonSnow, mock.AryaStark}
		return nil
	})
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error {
		webhook = true
		return nil
	})

	err := mock.NewWorker().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow). // creator or approving admin; recipients must not depend on it
		WithBaseURL("http://domain.com").
		Execute(task)
	Expect(err).IsNil()

	for _, m := range emailmock.MessageHistory {
		emailed = append(emailed, m.To...)
	}
	return notified, emailed, webhook
}

func pendingPost() *entity.Post {
	return &entity.Post{ID: 1, Number: 1, Title: "Pending idea", Slug: "pending-idea", Description: "secret @[Arya Stark]", User: pendingAuthor}
}

func TestNotifyAboutNewPost_PendingOnlyNotifiesStaff(t *testing.T) {
	RegisterT(t)

	notified, emailed, webhook := moderationRun(tasks.NotifyAboutNewPost(pendingPost()))

	Expect(notified).Equals([]*entity.User{mock.JonSnow})
	Expect(emailed).HasLen(1)
	Expect(emailed[0].Address).Equals(mock.JonSnow.Email)
	Expect(webhook).IsFalse()
}

func TestNotifyAboutApprovedPost_NotifiesNonStaff(t *testing.T) {
	RegisterT(t)

	post := pendingPost()
	post.IsApproved = true
	notified, emailed, webhook := moderationRun(tasks.NotifyAboutApprovedPost(post))

	// new post + mention for Arya, nothing for the admin who already had it
	Expect(notified).Equals([]*entity.User{mock.AryaStark, mock.AryaStark})
	Expect(emailed).HasLen(2)
	Expect(emailed[0].Address).Equals(mock.AryaStark.Email)
	Expect(emailed[1].Address).Equals(mock.AryaStark.Email)
	Expect(emailmock.MessageHistory[0].From.Name).Equals(pendingAuthor.Name)
	Expect(webhook).IsTrue()
}

func TestNotifyAboutNewComment_PendingOnlyNotifiesStaff(t *testing.T) {
	RegisterT(t)

	post := pendingPost()
	post.IsApproved = true
	comment := &entity.Comment{ID: 1, Content: "hey @[Arya Stark]", User: pendingAuthor}
	notified, emailed, webhook := moderationRun(tasks.NotifyAboutNewComment(comment, post))

	Expect(notified).Equals([]*entity.User{mock.JonSnow})
	Expect(emailed).HasLen(1)
	Expect(emailed[0].Address).Equals(mock.JonSnow.Email)
	Expect(webhook).IsFalse()
}
