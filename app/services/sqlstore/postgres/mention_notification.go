package postgres

import (
	"context"
	"database/sql"
	"fmt"
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

		// Create NullInt64 variables for CommentID and PostID
		var commentID sql.NullInt64
		var postID sql.NullInt64

		// Only set Valid=true if the ID is not 0
		if c.CommentID != 0 {
			commentID.Int64 = int64(c.CommentID)
			commentID.Valid = true
		}

		if c.PostID != 0 {
			postID.Int64 = int64(c.PostID)
			postID.Valid = true
		}

		query := `
			INSERT INTO mention_notifications (tenant_id, user_id, comment_id, post_id, created_on)
			SELECT $1, $2, $3, $4, $5
			WHERE NOT EXISTS (
				SELECT 1
				FROM mention_notifications
				WHERE tenant_id = $1
				AND user_id = $2
				AND COALESCE(comment_id, -1) = COALESCE($3, -1)
				AND COALESCE(post_id, -1) = COALESCE($4, -1)
			);
		`
		_, err := trx.Execute(query,
			tenant.ID,
			c.UserID,
			commentID,
			postID,
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
		}

		query := `
			SELECT id, tenant_id, user_id, comment_id, post_id, created_on
			FROM mention_notifications
			WHERE tenant_id = $1
		`

		paramCount := 1

		if q.CommentID > 0 {
			paramCount++
			query += fmt.Sprintf(" AND comment_id = $%d", paramCount)
			params = append(params, q.CommentID)
		}

		if q.PostID > 0 {
			paramCount++
			query += fmt.Sprintf(" AND post_id = $%d", paramCount)
			params = append(params, q.PostID)
		}

		query += " ORDER BY created_on DESC"

		err := trx.Select(&q.Result, query, params...)
		if err != nil {
			return errors.Wrap(err, "failed to get mention notification logs")
		}

		return nil
	})
}
