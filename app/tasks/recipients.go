package tasks

import (
	"context"
	"slices"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
)

// recipients is who gets notified about one comment or post on one channel.
// Each user is in at most one group.
type recipients struct {
	reply   []*entity.User
	mention []*entity.User
	general []*entity.User // new comment or new post
}

// recipientRules describes one comment or post being notified about
type recipientRules struct {
	authorID int
	// replyTo is the author of the comment being replied to, if any
	replyTo  *entity.User
	mentions []string
	// alreadyNotified have had a reply or mention notification for this content before
	alreadyNotified []*entity.MentionNotification
	// general is the new comment/new post event, or nil when only mentions apply (edits)
	general *enum.NotificationEvent
	aud     audience
	// skipUserID is never notified, e.g. the moderator who just approved the content
	skipUserID int
}

// pickRecipients gives each user at most one notification on a channel.
// A reply wins over a mention, which wins over the general event. Replies use the mention setting.
func pickRecipients(ctx context.Context, post *entity.Post, channel enum.NotificationChannel, r recipientRules) (recipients, error) {
	var result recipients

	done := map[int]bool{r.authorID: true, r.skipUserID: true}
	for _, n := range r.alreadyNotified {
		done[n.UserID] = true
	}

	if r.replyTo != nil || len(r.mentions) > 0 {
		users, err := getSubscribers(ctx, post, channel, enum.NotificationEventMention, r.aud)
		if err != nil {
			return result, err
		}

		if r.replyTo != nil && !done[r.replyTo.ID] {
			if i := slices.IndexFunc(users, func(u *entity.User) bool { return u.ID == r.replyTo.ID }); i >= 0 {
				result.reply = append(result.reply, users[i])
				done[r.replyTo.ID] = true
			}
		}

		// Iterating users rather than mentions notifies someone once, however often they're mentioned
		for _, u := range users {
			if !done[u.ID] && slices.Contains(r.mentions, u.Name) {
				result.mention = append(result.mention, u)
				done[u.ID] = true
			}
		}
	}

	if r.general != nil {
		users, err := getSubscribers(ctx, post, channel, *r.general, r.aud)
		if err != nil {
			return result, err
		}
		for _, u := range users {
			if !done[u.ID] {
				result.general = append(result.general, u)
				done[u.ID] = true
			}
		}
	}

	return result, nil
}

func addWebNotifications(ctx context.Context, users []*entity.User, title, link string, postID int) error {
	for _, u := range users {
		if err := bus.Dispatch(ctx, &cmd.AddNewNotification{User: u, Title: title, Link: link, PostID: postID}); err != nil {
			return err
		}
	}
	return nil
}

// logMentionNotifications records who had a reply or mention notification, so an edit doesn't notify them again
func logMentionNotifications(ctx context.Context, web, email recipients, commentID, postID int) error {
	logged := map[int]bool{}
	for _, group := range [][]*entity.User{web.reply, web.mention, email.reply, email.mention} {
		for _, u := range group {
			if logged[u.ID] {
				continue
			}
			logged[u.ID] = true
			if err := bus.Dispatch(ctx, &cmd.AddMentionNotification{UserID: u.ID, CommentID: commentID, PostID: postID}); err != nil {
				return err
			}
		}
	}
	return nil
}

func toRecipients(users []*entity.User) []dto.Recipient {
	to := make([]dto.Recipient, 0, len(users))
	for _, u := range users {
		to = append(to, dto.NewRecipient(u.Name, u.Email, dto.Props{}))
	}
	return to
}
