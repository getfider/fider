package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)


func addNewComment(ctx context.Context, c *cmd.AddNewComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		isApproved := !tenant.IsModerationEnabled || !user.RequiresModeration()
		var id int
		if err := trx.Get(&id, `
			INSERT INTO comments (tenant_id, post_id, content, user_id, created_at, is_approved) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id
		`, tenant.ID, c.Post.ID, c.Content, user.ID, time.Now(), isApproved); err != nil {
			return errors.Wrap(err, "failed add new comment")
		}

		q := &query.GetCommentByID{CommentID: id}
		if err := getCommentByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result

		return nil
	})
}

func toggleCommentReaction(ctx context.Context, c *cmd.ToggleCommentReaction) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var added bool
		err := trx.Scalar(&added, `
			WITH toggle_reaction AS (
				INSERT INTO reactions (comment_id, user_id, emoji, created_on)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (comment_id, user_id, emoji) DO NOTHING
				RETURNING true AS added
			),
			delete_existing AS (
				DELETE FROM reactions
				WHERE comment_id = $1 AND user_id = $2 AND emoji = $3
				AND NOT EXISTS (SELECT 1 FROM toggle_reaction)
				RETURNING false AS added
			)
			SELECT COALESCE(
				(SELECT added FROM toggle_reaction),
				(SELECT added FROM delete_existing),
				false
			)
		`, c.Comment.ID, user.ID, c.Emoji, time.Now())

		if err != nil {
			return errors.Wrap(err, "failed to toggle reaction")
		}

		c.Result = added
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

		comment := dbEntities.Comment{}
		err := trx.Get(&comment,
			`SELECT c.id, 
							c.content, 
							c.created_at, 
							c.edited_at, 
							c.is_approved,
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

		q.Result = comment.ToModel(ctx)
		return nil
	})
}

func getCommentsByPost(ctx context.Context, q *query.GetCommentsByPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Comment, 0)

		comments := []*dbEntities.Comment{}
		userId := 0
		if user != nil {
			userId = user.ID
		}
		
		// Build approval filter based on user permissions
		approvalFilter := ""
		if user != nil && user.IsCollaborator() {
			// Admins and collaborators can see all comments
			approvalFilter = ""
		} else if user != nil {
			// Regular users can see approved comments + their own unapproved comments
			approvalFilter = fmt.Sprintf(" AND (c.is_approved = true OR c.user_id = %d)", user.ID)
		} else {
			// Anonymous users can only see approved comments
			approvalFilter = " AND c.is_approved = true"
		}
		
		query := fmt.Sprintf(`
			WITH agg_attachments AS ( 
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
			),
			agg_reactions AS (
				SELECT 
					comment_id,
					json_agg(json_build_object(
						'emoji', emoji,
						'count', count,
						'includesMe', CASE WHEN $3 = ANY(user_ids) THEN true ELSE false END
					) ORDER BY count DESC) as reaction_counts
				FROM (
					SELECT 
						comment_id, 
						emoji, 
						COUNT(*) as count,
						array_agg(user_id) as user_ids
					FROM reactions
					WHERE comment_id IN (SELECT id FROM comments WHERE post_id = $1)
					GROUP BY comment_id, emoji
				) r
				GROUP BY comment_id
			)
			SELECT c.id, 
					c.content, 
					c.created_at, 
					c.edited_at, 
					c.is_approved,
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
					at.attachment_bkeys,
					ar.reaction_counts
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
			LEFT JOIN agg_reactions ar
			ON ar.comment_id = c.id
			WHERE p.id = $1
			AND p.tenant_id = $2
			AND c.deleted_at IS NULL%s
			ORDER BY c.created_at DESC`, approvalFilter)
		
		err := trx.Select(&comments, query, q.Post.ID, tenant.ID, userId)
		if err != nil {
			return errors.Wrap(err, "failed get comments of post with id '%d'", q.Post.ID)
		}

		q.Result = make([]*entity.Comment, len(comments))
		for i, comment := range comments {
			q.Result[i] = comment.ToModel(ctx)
		}
		return nil
	})
}
