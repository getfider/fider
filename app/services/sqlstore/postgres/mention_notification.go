package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

func AddMentionNotification(ctx context.Context, c *cmd.AddMentionNotification) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		now := time.Now()

		query := `
			INSERT INTO mention_notifications (tenant_id, user_id, comment_id, created_on) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (tenant_id, user_id, comment_id) DO NOTHING
		`
		_, err := trx.Execute(query,
			tenant.ID,
			c.UserID,
			c.CommentID,
			now)

		if err != nil {
			return errors.Wrap(err, "failed to insert notification log")
		}

		return nil
	})
}

func getMentionsNotifications(ctx context.Context, q *query.GetMentionNotifications) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.MentionNotification, 0)

		params := []interface{}{
			tenant.ID,
			q.CommentID,
		}

		query := `
			SELECT id, tenant_id, user_id, comment_id, created_on
			FROM mention_notifications
			WHERE tenant_id = $1 
			AND comment_id = $2
		`

		query += " ORDER BY created_on DESC"

		err := trx.Select(&q.Result, query, params...)
		if err != nil {
			return errors.Wrap(err, "failed to get mention notification logs")
		}

		return nil
	})
}
