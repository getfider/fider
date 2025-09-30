package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func ApprovePost(ctx context.Context, c *cmd.ApprovePost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			UPDATE posts SET is_approved = true
			WHERE id = $1 AND tenant_id = $2`, c.PostID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to approve post")
		}
		return nil
	})
}

func DeclinePost(ctx context.Context, c *cmd.DeclinePost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Use existing query to get the post
		getPost := &query.GetPostByID{PostID: c.PostID}
		if err := bus.Dispatch(ctx, getPost); err != nil {
			return err
		}

		// Use SetPostResponse to properly delete the post
		setResponse := &cmd.SetPostResponse{
			Post:   getPost.Result,
			Text:   "Post declined during moderation",
			Status: enum.PostDeleted,
		}

		return bus.Dispatch(ctx, setResponse)
	})
}

func ApproveComment(ctx context.Context, c *cmd.ApproveComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			UPDATE comments SET is_approved = true
			WHERE id = $1 AND tenant_id = $2`, c.CommentID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to approve comment")
		}
		return nil
	})
}

func DeclineComment(ctx context.Context, c *cmd.DeclineComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Use existing delete command
		deleteComment := &cmd.DeleteComment{CommentID: c.CommentID}
		return bus.Dispatch(ctx, deleteComment)
	})
}

func BulkApproveItems(ctx context.Context, c *cmd.BulkApproveItems) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if len(c.PostIDs) > 0 {
			postIDsStr := ""
			for i, id := range c.PostIDs {
				if i > 0 {
					postIDsStr += ","
				}
				postIDsStr += fmt.Sprintf("%d", id)
			}
			_, err := trx.Execute(fmt.Sprintf(`
				UPDATE posts SET is_approved = true
				WHERE id IN (%s) AND tenant_id = $1`, postIDsStr), tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed to bulk approve posts")
			}
		}

		if len(c.CommentIDs) > 0 {
			commentIDsStr := ""
			for i, id := range c.CommentIDs {
				if i > 0 {
					commentIDsStr += ","
				}
				commentIDsStr += fmt.Sprintf("%d", id)
			}
			_, err := trx.Execute(fmt.Sprintf(`
				UPDATE comments SET is_approved = true
				WHERE id IN (%s) AND tenant_id = $1`, commentIDsStr), tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed to bulk approve comments")
			}
		}

		return nil
	})
}

func BulkDeclineItems(ctx context.Context, c *cmd.BulkDeclineItems) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if len(c.PostIDs) > 0 {
			// Use existing commands to properly delete each post
			for _, postID := range c.PostIDs {
				// Use existing query to get the post
				getPost := &query.GetPostByID{PostID: postID}
				if err := bus.Dispatch(ctx, getPost); err != nil {
					return errors.Wrap(err, "failed to get post for bulk decline")
				}

				// Use SetPostResponse to properly delete the post
				setResponse := &cmd.SetPostResponse{
					Post:   getPost.Result,
					Text:   "Post declined during bulk moderation",
					Status: enum.PostDeleted,
				}

				if err := bus.Dispatch(ctx, setResponse); err != nil {
					return errors.Wrap(err, "failed to bulk decline post")
				}
			}
		}

		if len(c.CommentIDs) > 0 {
			// Use existing delete command for each comment
			for _, commentID := range c.CommentIDs {
				deleteComment := &cmd.DeleteComment{CommentID: commentID}
				if err := bus.Dispatch(ctx, deleteComment); err != nil {
					return errors.Wrap(err, "failed to bulk decline comment")
				}
			}
		}

		return nil
	})
}

type dbModerationPost struct {
	ID          int              `db:"id"`
	Number      int              `db:"number"`
	Title       string           `db:"title"`
	Slug        string           `db:"slug"`
	Description string           `db:"description"`
	CreatedAt   time.Time        `db:"created_at"`
	User        *dbEntities.User `db:"user"`
}

type dbModerationComment struct {
	ID         int              `db:"id"`
	PostID     int              `db:"post_id"`
	PostNumber int              `db:"post_number"`
	PostSlug   string           `db:"post_slug"`
	Content    string           `db:"content"`
	CreatedAt  time.Time        `db:"created_at"`
	User       *dbEntities.User `db:"user"`
	PostTitle  string           `db:"post_title"`
}

func GetModerationItems(ctx context.Context, q *query.GetModerationItems) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.ModerationItem, 0)

		// Get unmoderated posts
		var posts []*dbModerationPost

		err := trx.Select(&posts, `
			SELECT p.id, p.number, p.title, p.slug, p.description, p.created_at,
				u.id AS user_id,
				u.name AS user_name,
				u.email AS user_email,
				u.role AS user_role,
				u.status AS user_status,
				u.avatar_type AS user_avatar_type,
				u.avatar_bkey AS user_avatar_bkey
			FROM posts p
			INNER JOIN users u ON u.id = p.user_id AND u.tenant_id = p.tenant_id
			WHERE p.tenant_id = $1 AND p.is_approved = false and p.status <> $2
			ORDER BY p.created_at DESC`, tenant.ID, enum.PostDeleted)
		if err != nil {
			return errors.Wrap(err, "failed to get unmoderated posts")
		}

		for _, post := range posts {
			userWithEmail := &entity.UserWithEmail{
				User: post.User.ToModel(ctx),
			}

			q.Result = append(q.Result, &query.ModerationItem{
				Type:       "post",
				ID:         post.ID,
				PostNumber: post.Number,
				PostSlug:   post.Slug,
				Title:      post.Title,
				Content:    post.Description,
				CreatedAt:  post.CreatedAt,
				User:       userWithEmail,
			})
		}

		// Get unmoderated comments
		var comments []*dbModerationComment

		err = trx.Select(&comments, `
			SELECT c.id, c.post_id, p.number as post_number, p.slug as post_slug, c.content, c.created_at,
					u.id AS user_id,
					u.name AS user_name,
					u.email AS user_email,
					u.role AS user_role,
					u.status AS user_status,
					u.avatar_type AS user_avatar_type,
					u.avatar_bkey AS user_avatar_bkey,
				   p.title as post_title
			FROM comments c
			INNER JOIN users u ON u.id = c.user_id AND u.tenant_id = c.tenant_id
			INNER JOIN posts p ON p.id = c.post_id AND p.tenant_id = c.tenant_id
			WHERE c.tenant_id = $1 AND c.is_approved = false and p.status <> $2
			AND c.deleted_at IS NULL
			ORDER BY c.created_at DESC`, tenant.ID, enum.PostDeleted)
		if err != nil {
			return errors.Wrap(err, "failed to get unmoderated comments")
		}

		for _, comment := range comments {

			userWithEmail := &entity.UserWithEmail{
				User: comment.User.ToModel(ctx),
			}

			q.Result = append(q.Result, &query.ModerationItem{
				Type:       "comment",
				ID:         comment.ID,
				PostID:     comment.PostID,
				PostNumber: comment.PostNumber,
				PostSlug:   comment.PostSlug,
				Content:    comment.Content,
				CreatedAt:  comment.CreatedAt,
				PostTitle:  comment.PostTitle,
				User:       userWithEmail,
			})
		}

		return nil
	})
}

func GetModerationCount(ctx context.Context, q *query.GetModerationCount) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var count int

		err := trx.Get(&count, `
			SELECT
				(SELECT COUNT(*) FROM posts WHERE tenant_id = $1 AND is_approved = false and status <> $2) +
				(SELECT COUNT(*) FROM comments c JOIN posts p on c.post_id = p.id WHERE p.tenant_id = $1 AND c.is_approved = false AND p.status <> $2 AND c.deleted_at IS NULL)
		`, tenant.ID, enum.PostDeleted)

		if err != nil {
			return errors.Wrap(err, "failed to get moderation count")
		}

		q.Result = count
		return nil
	})
}

func VerifyUser(ctx context.Context, c *cmd.VerifyUser) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Verify and unblock the user (can't be both blocked and verified)
		_, err := trx.Execute(`
			UPDATE users SET is_verified = true, status = $1
			WHERE id = $2 AND tenant_id = $3`, enum.UserActive, c.UserID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to verify user")
		}

		return nil
	})
}

func using(ctx context.Context, handler func(*dbx.Trx, *entity.Tenant, *entity.User) error) error {
	trx, _ := ctx.Value(app.TransactionCtxKey).(*dbx.Trx)
	tenant, _ := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	user, _ := ctx.Value(app.UserCtxKey).(*entity.User)
	return handler(trx, tenant, user)
}
