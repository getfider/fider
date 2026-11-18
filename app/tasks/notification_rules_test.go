package tasks_test

import (
	"context"
	"slices"
	"testing"

	"github.com/getfider/fider/app"

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

var sansaStark = &entity.User{ID: 3, Name: "Sansa Stark", Email: "sansa@got.com", Role: enum.RoleVisitor}

var (
	chWeb   = enum.NotificationChannelWeb
	chEmail = enum.NotificationChannelEmail
)

// subscribers maps an event's settings key and a channel to the users who have it enabled
type subscribers map[string]map[enum.NotificationChannel][]*entity.User

func everywhere(users ...*entity.User) map[enum.NotificationChannel][]*entity.User {
	return map[enum.NotificationChannel][]*entity.User{chWeb: users, chEmail: users}
}

type webNotification struct {
	user  string
	title string
}

type sentEmail struct {
	text string
	to   []string
}

type rulesResult struct {
	web    []webNotification
	emails []sentEmail
	logs   []int
}

// rulesRun runs a task where the parent of any reply is written by parentAuthor
func rulesRun(task worker.Task, subs subscribers, parentAuthor *entity.User) rulesResult {
	var result rulesResult
	bus.Init(emailmock.Service{})
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		result.web = append(result.web, webNotification{c.User.Name, c.Title})
		return nil
	})
	bus.AddHandler(func(ctx context.Context, c *cmd.AddMentionNotification) error {
		result.logs = append(result.logs, c.UserID)
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error {
		q.Result = []*entity.MentionNotification{}
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = slices.Clone(subs[q.Event.UserSettingsKeyName][q.Channel]) // the task filters it in place
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = &entity.Comment{ID: q.CommentID, User: parentAuthor}
		return nil
	})
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error { return nil })

	err := mock.NewWorker().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)
	Expect(err).IsNil()

	for _, m := range emailmock.MessageHistory {
		e := sentEmail{text: m.TemplateName}
		if s, ok := m.Props["messageLocaleString"].(string); ok {
			e.text = s
		}
		for _, to := range m.To {
			e.to = append(e.to, to.Name)
		}
		result.emails = append(result.emails, e)
	}
	return result
}

func rulesPost() *entity.Post {
	return &entity.Post{ID: 1, Number: 1, Title: "Dark mode", Slug: "dark-mode", User: mock.JonSnow, IsApproved: true}
}

func reply(content string, author *entity.User) *entity.Comment {
	parentID := 10
	return &entity.Comment{ID: 11, ParentID: &parentID, Content: content, User: author, IsApproved: true}
}

const (
	mentionKey    = "event_notification_mention"
	newCommentKey = "event_notification_new_comment"
	newPostKey    = "event_notification_new_post"
)

func TestNotificationRules_SelfMentionDoesNotNotify(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{ID: 1, Content: "Note to self @[Arya Stark]", User: mock.AryaStark, IsApproved: true}
	result := rulesRun(tasks.NotifyAboutNewComment(comment, rulesPost()), subscribers{mentionKey: everywhere(mock.AryaStark)}, nil)

	Expect(result.web).HasLen(0)
	Expect(result.emails).HasLen(0)
	Expect(result.logs).HasLen(0)
}

func TestNotificationRules_SelfMentionInPostDoesNotNotify(t *testing.T) {
	RegisterT(t)

	post := rulesPost()
	post.User = mock.AryaStark
	post.Description = "Me again @[Arya Stark]"
	result := rulesRun(tasks.NotifyAboutNewPost(post), subscribers{mentionKey: everywhere(mock.AryaStark)}, nil)

	Expect(result.web).HasLen(0)
	Expect(result.emails).HasLen(0)
}

func TestNotificationRules_RepeatedMentionNotifiesOnce(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{ID: 1, Content: "@[Jon Snow] look, @[Jon Snow]!", User: mock.AryaStark, IsApproved: true}
	result := rulesRun(tasks.NotifyAboutNewComment(comment, rulesPost()), subscribers{mentionKey: everywhere(mock.JonSnow)}, nil)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Arya Stark** mentioned you in **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_mention.text", []string{"Jon Snow"}}})
	Expect(result.logs).Equals([]int{mock.JonSnow.ID})
}

func TestNotificationRules_ReplyNotifiesParentAuthor(t *testing.T) {
	RegisterT(t)

	result := rulesRun(tasks.NotifyAboutNewComment(reply("Good point", mock.AryaStark), rulesPost()), subscribers{mentionKey: everywhere(mock.JonSnow)}, mock.JonSnow)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Arya Stark** replied to your comment on **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_reply.text", []string{"Jon Snow"}}})
	Expect(result.logs).Equals([]int{mock.JonSnow.ID})
}

func TestNotificationRules_ReplyInOwnThreadDoesNotNotify(t *testing.T) {
	RegisterT(t)

	result := rulesRun(tasks.NotifyAboutNewComment(reply("Adding to my point", mock.AryaStark), rulesPost()), subscribers{mentionKey: everywhere(mock.AryaStark)}, mock.AryaStark)

	Expect(result.web).HasLen(0)
	Expect(result.emails).HasLen(0)
}

func TestNotificationRules_ReplyNeedsMentionSetting(t *testing.T) {
	RegisterT(t)

	// Jon has turned off mentions and replies, but gets new comments
	result := rulesRun(tasks.NotifyAboutNewComment(reply("Good point", mock.AryaStark), rulesPost()), subscribers{newCommentKey: everywhere(mock.JonSnow)}, mock.JonSnow)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Arya Stark** left a comment on **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_comment.text", []string{"Jon Snow"}}})
	Expect(result.logs).HasLen(0)
}

func TestNotificationRules_OneNotificationPerUser_ReplyWins(t *testing.T) {
	RegisterT(t)

	// Jon is replied to, mentioned and subscribed to new comments; Sansa only to new comments
	subs := subscribers{
		mentionKey:    everywhere(mock.JonSnow),
		newCommentKey: everywhere(mock.JonSnow, sansaStark),
	}
	result := rulesRun(tasks.NotifyAboutNewComment(reply("Agreed @[Jon Snow]", mock.AryaStark), rulesPost()), subs, mock.JonSnow)

	Expect(result.web).Equals([]webNotification{
		{"Jon Snow", "**Arya Stark** replied to your comment on **Dark mode**"},
		{"Sansa Stark", "**Arya Stark** left a comment on **Dark mode**"},
	})
	Expect(result.emails).Equals([]sentEmail{
		{"email.new_reply.text", []string{"Jon Snow"}},
		{"email.new_comment.text", []string{"Sansa Stark"}},
	})
	Expect(result.logs).Equals([]int{mock.JonSnow.ID})
}

func TestNotificationRules_OneNotificationPerUser_MentionBeatsNewComment(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{ID: 1, Content: "What do you think @[Jon Snow]?", User: mock.AryaStark, IsApproved: true}
	subs := subscribers{
		mentionKey:    everywhere(mock.JonSnow),
		newCommentKey: everywhere(mock.JonSnow),
	}
	result := rulesRun(tasks.NotifyAboutNewComment(comment, rulesPost()), subs, nil)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Arya Stark** mentioned you in **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_mention.text", []string{"Jon Snow"}}})
}

func TestNotificationRules_DecidedPerChannel(t *testing.T) {
	RegisterT(t)

	// Jon gets mentions and replies on the web only, and new comments by email only
	subs := subscribers{
		mentionKey:    {chWeb: {mock.JonSnow}},
		newCommentKey: {chEmail: {mock.JonSnow}},
	}
	result := rulesRun(tasks.NotifyAboutNewComment(reply("Good point", mock.AryaStark), rulesPost()), subs, mock.JonSnow)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Arya Stark** replied to your comment on **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_comment.text", []string{"Jon Snow"}}})
}

func TestNotificationRules_NewPostMentionBeatsNewPost(t *testing.T) {
	RegisterT(t)

	post := rulesPost()
	post.User = mock.AryaStark
	post.Description = "cc @[Jon Snow]"
	subs := subscribers{
		mentionKey: everywhere(mock.JonSnow),
		newPostKey: everywhere(mock.JonSnow, sansaStark),
	}
	result := rulesRun(tasks.NotifyAboutNewPost(post), subs, nil)

	Expect(result.web).Equals([]webNotification{
		{"Jon Snow", "**Arya Stark** mentioned you in **Dark mode**"},
		{"Sansa Stark", "New post: **Dark mode**"},
	})
	Expect(result.emails).Equals([]sentEmail{
		{"email.new_mention.text", []string{"Jon Snow"}},
		{"new_post", []string{"Sansa Stark"}},
	})
}

// Arya (a visitor) is replied to by a comment that needs moderation; Jon (staff) gets new comments
var pendingReplySubs = subscribers{
	mentionKey:    everywhere(mock.AryaStark),
	newCommentKey: everywhere(mock.JonSnow),
}

func TestNotificationRules_PendingReplyOnlyNotifiesStaff(t *testing.T) {
	RegisterT(t)

	pending := reply("Hmm", sansaStark)
	pending.IsApproved = false
	result := rulesRun(tasks.NotifyAboutNewComment(pending, rulesPost()), pendingReplySubs, mock.AryaStark)

	Expect(result.web).Equals([]webNotification{{"Jon Snow", "**Sansa Stark** left a comment on **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_comment.text", []string{"Jon Snow"}}})
}

func TestNotificationRules_ApprovedReplyNotifiesParentAuthor(t *testing.T) {
	RegisterT(t)

	result := rulesRun(tasks.NotifyAboutApprovedComment(reply("Hmm", sansaStark), rulesPost()), pendingReplySubs, mock.AryaStark)

	Expect(result.web).Equals([]webNotification{{"Arya Stark", "**Sansa Stark** replied to your comment on **Dark mode**"}})
	Expect(result.emails).Equals([]sentEmail{{"email.new_reply.text", []string{"Arya Stark"}}})
}

func TestNotificationRules_DeletedParentSkipsReply(t *testing.T) {
	RegisterT(t)

	bus.Init(emailmock.Service{})
	notified := 0
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		notified++
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetMentionNotifications) error { return nil })
	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{mock.JonSnow}
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooks) error { return nil })

	err := mock.NewWorker().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(tasks.NotifyAboutNewComment(reply("Too late", mock.AryaStark), rulesPost()))

	// Falls back to the new comment notification
	Expect(err).IsNil()
	Expect(notified).Equals(1)
}
