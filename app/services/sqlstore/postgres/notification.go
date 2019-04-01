package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

func markAllNotificationsAsRead(ctx context.Context, c *cmd.MarkAllNotificationsAsRead) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
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
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
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
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
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
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		q.Result = nil
		notification := &models.Notification{}

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
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		q.Result = []*models.Notification{}
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
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		c.Result = nil
		if user.ID == c.User.ID {
			return nil
		}

		now := time.Now()
		notification := &models.Notification{
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
