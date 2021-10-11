package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/lib/pq"
)

func purgeExpiredNotifications(ctx context.Context, c *cmd.PurgeExpiredNotifications) error {
	log.Debug(ctx, "deleting notifications more than 1 year old")

	trx, err := dbx.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to open transaction")
	}

	count, err := trx.Execute("DELETE FROM notifications WHERE CREATED_AT <= NOW() - INTERVAL '365 days'")
	if err != nil {
		return errors.Wrap(err, "failed to delete expired notifications")
	}

	log.Debugf(ctx, "@{RowsDeleted} notifications were deleted", dto.Props{
		"RowsDeleted": count,
	})

	if err = trx.Commit(); err != nil {
		return errors.Wrap(err, "failed commit transaction")
	}

	c.NumOfDeletedNotifications = int(count)
	return nil
}

func markAllNotificationsAsRead(ctx context.Context, c *cmd.MarkAllNotificationsAsRead) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if user == nil {
			return nil
		}
		_, err := trx.Execute(`
			UPDATE notifications SET read = true, updated_at = $1
			WHERE tenant_id = $2 AND user_id = $3 AND read = false
		`, time.Now(), tenant.ID, user.ID)
		if err != nil {
			return errors.Wrap(err, "failed to mark all notifications as read")
		}
		return nil
	})
}

func countUnreadNotifications(ctx context.Context, q *query.CountUnreadNotifications) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = 0

		if user != nil {
			err := trx.Scalar(&q.Result, "SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = false", tenant.ID, user.ID)
			if err != nil {
				return errors.Wrap(err, "failed count total unread notifications")
			}
		}
		return nil
	})
}

func markNotificationAsRead(ctx context.Context, c *cmd.MarkNotificationAsRead) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if user == nil {
			return nil
		}

		_, err := trx.Execute(`
			UPDATE notifications SET read = true, updated_at = $1
			WHERE id = $2 AND tenant_id = $3 AND user_id = $4 AND read = false
		`, time.Now(), c.ID, tenant.ID, user.ID)
		if err != nil {
			return errors.Wrap(err, "failed to mark notification as read")
		}
		return nil
	})
}

func getNotificationByID(ctx context.Context, q *query.GetNotificationByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil
		notification := &entity.Notification{}

		err := trx.Get(notification, `
			SELECT id, title, link, read, created_at 
			FROM notifications
			WHERE id = $1 AND tenant_id = $2 AND user_id = $3
		`, q.ID, tenant.ID, user.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get notifications with id '%d'", q.ID)
		}

		q.Result = notification
		return nil
	})
}

func getActiveNotifications(ctx context.Context, q *query.GetActiveNotifications) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = []*entity.Notification{}
		err := trx.Select(&q.Result, `
			SELECT id, title, link, read, created_at 
			FROM notifications 
			WHERE tenant_id = $1 AND user_id = $2
			AND (read = false OR updated_at > CURRENT_DATE - INTERVAL '30 days')
		`, tenant.ID, user.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get active notifications")
		}
		return nil
	})
}

func addNewNotification(ctx context.Context, c *cmd.AddNewNotification) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		c.Result = nil
		if user.ID == c.User.ID {
			return nil
		}

		now := time.Now()
		notification := &entity.Notification{
			Title:     c.Title,
			Link:      c.Link,
			CreatedAt: now,
			Read:      false,
		}
		err := trx.Get(&notification.ID, `
			INSERT INTO notifications (tenant_id, user_id, title, link, read, post_id, author_id, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
			RETURNING id
		`, tenant.ID, c.User.ID, c.Title, c.Link, false, c.PostID, user.ID, now)
		if err != nil {
			return errors.Wrap(err, "failed to insert notification")
		}

		c.Result = notification
		return nil
	})
}

func addSubscriber(ctx context.Context, c *cmd.AddSubscriber) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		return internalAddSubscriber(trx, c.Post, tenant, c.User, true)
	})
}

func removeSubscriber(ctx context.Context, c *cmd.RemoveSubscriber) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			INSERT INTO post_subscribers (tenant_id, user_id, post_id, created_at, updated_at, status)
			VALUES ($1, $2, $3, $4, $4, $5) ON CONFLICT (user_id, post_id)
			DO UPDATE SET status = 0, updated_at = $4`,
			tenant.ID, c.User.ID, c.Post.ID, time.Now(), enum.SubscriberInactive,
		)
		if err != nil {
			return errors.Wrap(err, "failed remove post subscriber")
		}
		return nil
	})
}

func getActiveSubscribers(ctx context.Context, q *query.GetActiveSubscribers) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.User, 0)

		var (
			users []*dbUser
			err   error
		)

		// When searching for email subscrivers, skip users with email supressed
		supressionCondition := ""
		if q.Channel == enum.NotificationChannelEmail {
			supressionCondition = "AND u.email_supressed_at IS NULL"
		}

		// If the event doesn't require a subscription, notify everyone
		if len(q.Event.RequiresSubscriptionUserRoles) == 0 {
			err = trx.Select(&users, fmt.Sprintf(`
				SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role, u.status
				FROM users u
				LEFT JOIN user_settings set
				ON set.user_id = u.id
				AND set.tenant_id = u.tenant_id
				AND set.key = $1
				WHERE u.tenant_id = $2
				AND u.status = $5
				%s
				AND (
					(set.value IS NULL AND u.role = ANY($3))
					OR CAST(set.value AS integer) & $4 > 0
				)
				ORDER by u.id`, supressionCondition),
				q.Event.UserSettingsKeyName,
				tenant.ID,
				pq.Array(q.Event.DefaultEnabledUserRoles),
				q.Channel,
				enum.UserActive,
			)
		} else {
			// If the event requires a subscription, notify only those who subscribed
			err = trx.Select(&users, fmt.Sprintf(`
				SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role, u.status
				FROM users u
				LEFT JOIN post_subscribers sub
				ON sub.user_id = u.id
				AND sub.post_id = (SELECT id FROM posts p WHERE p.tenant_id = $4 and p.number = $1 LIMIT 1)
				AND sub.tenant_id = u.tenant_id
				LEFT JOIN user_settings set
				ON set.user_id = u.id
				AND set.key = $3
				AND set.tenant_id = u.tenant_id
				WHERE u.tenant_id = $4
				AND u.status = $8
				%s
				AND ( sub.status = $2 OR (sub.status IS NULL AND NOT u.role = ANY($7)) )
				AND (
					(set.value IS NULL AND u.role = ANY($5))
					OR CAST(set.value AS integer) & $6 > 0
				)
				ORDER by u.id`, supressionCondition),
				q.Number,
				enum.SubscriberActive,
				q.Event.UserSettingsKeyName,
				tenant.ID,
				pq.Array(q.Event.DefaultEnabledUserRoles),
				q.Channel,
				pq.Array(q.Event.RequiresSubscriptionUserRoles),
				enum.UserActive,
			)
		}

		if err != nil {
			return errors.Wrap(err, "failed to get post number '%d' subscribers", q.Number)
		}

		q.Result = make([]*entity.User, len(users))
		for i, user := range users {
			q.Result[i] = user.toModel(ctx)
		}
		return nil
	})
}

func internalAddSubscriber(trx *dbx.Trx, post *entity.Post, tenant *entity.Tenant, user *entity.User, force bool) error {
	conflict := " DO NOTHING"
	if force {
		conflict = "(user_id, post_id) DO UPDATE SET status = $5, updated_at = $4"
	}

	_, err := trx.Execute(fmt.Sprintf(`
	INSERT INTO post_subscribers (tenant_id, user_id, post_id, created_at, updated_at, status)
	VALUES ($1, $2, $3, $4, $4, $5)  ON CONFLICT %s`, conflict),
		tenant.ID, user.ID, post.ID, time.Now(), enum.SubscriberActive,
	)
	if err != nil {
		return errors.Wrap(err, "failed insert post subscriber")
	}
	return nil
}

func supressEmail(ctx context.Context, c *cmd.SupressEmail) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		cmd := "UPDATE users SET email_supressed_at = $1 WHERE email = ANY($2) AND email_supressed_at IS NULL"
		rowsCount, err := trx.Execute(cmd, time.Now(), pq.Array(c.EmailAddresses))
		if err != nil {
			return errors.Wrap(err, "failed to update supress email: %s", strings.Join(c.EmailAddresses, ","))
		}
		c.NumOfSupressedEmailAddresses = int(rowsCount)
		return nil
	})
}
