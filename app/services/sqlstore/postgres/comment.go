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

type dbComment struct {
	ID          int          `db:"id"`
	Content     string       `db:"content"`
	CreatedAt   time.Time    `db:"created_at"`
	User        *dbUser      `db:"user"`
	Attachments []string     `db:"attachment_bkeys"`
	EditedAt    dbx.NullTime `db:"edited_at"`
	EditedBy    *dbUser      `db:"edited_by"`
}

func (c *dbComment) toModel(ctx context.Context) *entity.Comment {
	comment := &entity.Comment{
		ID:          c.ID,
		Content:     c.Content,
		CreatedAt:   c.CreatedAt,
		User:        c.User.toModel(ctx),
		Attachments: c.Attachments,
	}
	if c.EditedAt.Valid {
		comment.EditedBy = c.EditedBy.toModel(ctx)
		comment.EditedAt = &c.EditedAt.Time
	}
	return comment
}

func addNewComment(ctx context.Context, c *cmd.AddNewComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var id int
		if err := trx.Get(&id, `
			INSERT INTO comments (tenant_id, post_id, content, user_id, created_at) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id
		`, tenant.ID, c.Post.ID, c.Content, user.ID, time.Now()); err != nil {
			return errors.Wrap(err, "failed add new comment")
		}

		if err := internalAddSubscriber(trx, c.Post, tenant, user, false); err != nil {
			return err
		}

		q := &query.GetCommentByID{CommentID: id}
		if err := getCommentByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result

		return nil
	})
}

func updateComment(ctx context.Context, c *cmd.UpdateComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			UPDATE comments SET content = $1, edited_at = $2, edited_by_id = $3 
			WHERE id = $4 AND tenant_id = $5`, c.Content, time.Now(), user.ID, c.CommentID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update comment")
		}
		return nil
	})
}

func deleteComment(ctx context.Context, c *cmd.DeleteComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if _, err := trx.Execute(
			"UPDATE comments SET deleted_at = $1, deleted_by_id = $2 WHERE id = $3 AND tenant_id = $4",
			time.Now(), user.ID, c.CommentID, tenant.ID,
		); err != nil {
			return errors.Wrap(err, "failed delete comment")
		}
		return nil
	})
}

func getCommentByID(ctx context.Context, q *query.GetCommentByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil

		comment := dbComment{}
		err := trx.Get(&comment,
			`SELECT c.id, 
							c.content, 
							c.created_at, 
							c.edited_at, 
							u.id AS user_id, 
							u.name AS user_name,
							u.email AS user_email,
							u.role AS user_role, 
							u.status AS user_status,
							u.avatar_type AS user_avatar_type,
							u.avatar_bkey AS user_avatar_bkey, 
							e.id AS edited_by_id, 
							e.name AS edited_by_name,
							e.email AS edited_by_email,
							e.role AS edited_by_role,
							e.status AS edited_by_status,
							e.avatar_type AS edited_by_avatar_type,
							e.avatar_bkey AS edited_by_avatar_bkey
			FROM comments c
			INNER JOIN users u
			ON u.id = c.user_id
			AND u.tenant_id = c.tenant_id
			LEFT JOIN users e
			ON e.id = c.edited_by_id
			AND e.tenant_id = c.tenant_id
			WHERE c.id = $1
			AND c.tenant_id = $2
			AND c.deleted_at IS NULL`, q.CommentID, tenant.ID)

		if err != nil {
			return err
		}

		q.Result = comment.toModel(ctx)
		return nil
	})
}

func getCommentsByPost(ctx context.Context, q *query.GetCommentsByPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Comment, 0)

		comments := []*dbComment{}
		err := trx.Select(&comments,
			`WITH agg_attachments AS ( 
					SELECT 
							c.id as comment_id, 
							ARRAY_REMOVE(ARRAY_AGG(at.attachment_bkey), NULL) as attachment_bkeys
					FROM attachments at
					INNER JOIN comments c
					ON at.tenant_id = c.tenant_id
					AND at.post_id = c.post_id
					AND at.comment_id = c.id
					WHERE at.post_id = $1
					AND at.tenant_id = $2
					AND at.comment_id IS NOT NULL
					GROUP BY c.id 
			)
			SELECT c.id, 
					c.content, 
					c.created_at, 
					c.edited_at, 
					u.id AS user_id, 
					u.name AS user_name,
					u.email AS user_email,
					u.role AS user_role, 
					u.status AS user_status, 
					u.avatar_type AS user_avatar_type, 
					u.avatar_bkey AS user_avatar_bkey, 
					e.id AS edited_by_id, 
					e.name AS edited_by_name,
					e.email AS edited_by_email,
					e.role AS edited_by_role,
					e.status AS edited_by_status,
					e.avatar_type AS edited_by_avatar_type, 
					e.avatar_bkey AS edited_by_avatar_bkey,
					at.attachment_bkeys
			FROM comments c
			INNER JOIN posts p
			ON p.id = c.post_id
			AND p.tenant_id = c.tenant_id
			INNER JOIN users u
			ON u.id = c.user_id
			AND u.tenant_id = c.tenant_id
			LEFT JOIN users e
			ON e.id = c.edited_by_id
			AND e.tenant_id = c.tenant_id
			LEFT JOIN agg_attachments at
			ON at.comment_id = c.id
			WHERE p.id = $1
			AND p.tenant_id = $2
			AND c.deleted_at IS NULL
			ORDER BY c.created_at ASC`, q.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed get comments of post with id '%d'", q.Post.ID)
		}

		q.Result = make([]*entity.Comment, len(comments))
		for i, comment := range comments {
			q.Result[i] = comment.toModel(ctx)
		}
		return nil
	})
}
