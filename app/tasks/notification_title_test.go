package tasks_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/cmd"
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

func TestNotificationTitles_EscapeMarkdownInNamesAndTitles(t *testing.T) {
	RegisterT(t)

	name := "[Your account needs verifying](https://evil.example) <b>x</b> **bold**"
	escapedName := `\[Your account needs verifying\]\(https\:\/\/evil\.example\) <b\>x<\/b\> \*\*bold\*\*`
	title := "_Urgent_ www.evil.example ![i](x.png)"
	escapedTitle := `\_Urgent\_ www\.evil\.example \!\[i\]\(x\.png\)`

	newPost := func() *entity.Post {
		return &entity.Post{
			ID:          1,
			Number:      1,
			Title:       title,
			Slug:        "urgent",
			Description: "Hello @[Jon Snow]",
			User:        mock.AryaStark,
			Status:      enum.PostDuplicate,
			Response: &entity.PostResponse{
				RespondedAt: time.Now(),
				User:        mock.JonSnow,
				Original:    &entity.OriginalPost{Number: 2, Title: "Original", Slug: "original"},
			},
		}
	}

	testCases := []struct {
		task     func(post *entity.Post) worker.Task
		expected []string
	}{
		{
			task: func(post *entity.Post) worker.Task { return tasks.NotifyAboutNewPost(post) },
			expected: []string{
				"New post: **" + escapedTitle + "**",
				"**" + escapedName + "** mentioned you in **" + escapedTitle + "**",
			},
		},
		{
			task: func(post *entity.Post) worker.Task { return tasks.NotifyAboutUpdatedPost(post) },
			expected: []string{
				"**" + escapedName + "** mentioned you in **" + escapedTitle + "**",
			},
		},
		{
			task: func(post *entity.Post) worker.Task {
				return tasks.NotifyAboutNewComment(&entity.Comment{Content: "Hello @[Jon Snow]"}, post)
			},
			expected: []string{
				"**" + escapedName + "** left a comment on **" + escapedTitle + "**",
				"**" + escapedName + "** mentioned you in **" + escapedTitle + "**",
			},
		},
		{
			task: func(post *entity.Post) worker.Task {
				return tasks.NotifyAboutUpdatedComment(post, &entity.Comment{Content: "Hello @[Jon Snow]"})
			},
			expected: []string{
				"**" + escapedName + "** mentioned you in **" + escapedTitle + "**",
			},
		},
		{
			task: func(post *entity.Post) worker.Task { return tasks.NotifyAboutStatusChange(post, enum.PostOpen) },
			expected: []string{
				"**" + escapedName + "** changed status of **" + escapedTitle + "** to **duplicate**",
			},
		},
		{
			task: func(post *entity.Post) worker.Task { return tasks.NotifyAboutDeletedPost(post, true) },
			expected: []string{
				"**" + escapedName + "** deleted **" + escapedTitle + "**",
			},
		},
	}

	for _, testCase := range testCases {
		bus.Init(emailmock.Service{})

		titles := make([]string, 0)
		bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
			titles = append(titles, c.Title)
			return nil
		})
		bus.AddHandler(func(ctx context.Context, c *cmd.AddMentionNotification) error {
			return nil
		})
		bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
			q.Result = []*entity.MentionNotification{}
			return nil
		})
		bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
			q.Result = []*entity.User{mock.JonSnow}
			return nil
		})
		bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error {
			return nil
		})

		author := *mock.AryaStark
		author.Name = name

		err := mock.NewWorker().
			OnTenant(mock.DemoTenant).
			AsUser(&author).
			WithBaseURL("http://domain.com").
			Execute(testCase.task(newPost()))

		Expect(err).IsNil()
		Expect(titles).Equals(testCase.expected)
	}
}
