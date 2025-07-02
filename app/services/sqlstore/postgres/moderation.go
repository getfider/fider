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
)

func approvePost(ctx context.Context, c *cmd.ApprovePost) error {
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

func declinePost(ctx context.Context, c *cmd.DeclinePost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			DELETE FROM posts 
			WHERE id = $1 AND tenant_id = $2`, c.PostID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to decline post")
		}
		return nil
	})
}

func approveComment(ctx context.Context, c *cmd.ApproveComment) error {
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

func declineComment(ctx context.Context, c *cmd.DeclineComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			DELETE FROM comments 
			WHERE id = $1 AND tenant_id = $2`, c.CommentID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to decline comment")
		}
		return nil
	})
}

func bulkApproveItems(ctx context.Context, c *cmd.BulkApproveItems) error {
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

func bulkDeclineItems(ctx context.Context, c *cmd.BulkDeclineItems) error {
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
				DELETE FROM posts 
				WHERE id IN (%s) AND tenant_id = $1`, postIDsStr), tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed to bulk decline posts")
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
				DELETE FROM comments 
				WHERE id IN (%s) AND tenant_id = $1`, commentIDsStr), tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed to bulk decline comments")
			}
		}

		return nil
	})
}

type dbModerationPost struct {
	ID          int       `db:"id"`
	Number      int       `db:"number"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UserID      int       `db:"user_id"`
	UserName    string    `db:"user_name"`
	UserEmail   string    `db:"user_email"`
}

type dbModerationComment struct {
	ID         int       `db:"id"`
	PostID     int       `db:"post_id"`
	PostNumber int       `db:"post_number"`
	PostSlug   string    `db:"post_slug"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	UserID     int       `db:"user_id"`
	UserName   string    `db:"user_name"`
	UserEmail  string    `db:"user_email"`
	PostTitle  string    `db:"post_title"`
}

func getModerationItems(ctx context.Context, q *query.GetModerationItems) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.ModerationItem, 0)

		// Get unmoderated posts
		var posts []*dbModerationPost

		err := trx.Select(&posts, `
			SELECT p.id, p.number, p.title, p.slug, p.description, p.created_at,
				   u.id as user_id, u.name as user_name, u.email as user_email
			FROM posts p
			INNER JOIN users u ON u.id = p.user_id AND u.tenant_id = p.tenant_id
			WHERE p.tenant_id = $1 AND p.is_approved = false
			ORDER BY p.created_at DESC`, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get unmoderated posts")
		}

		for _, post := range posts {
			q.Result = append(q.Result, &query.ModerationItem{
				Type:      "post",
				ID:        post.ID,
				PostNumber: post.Number,
				PostSlug:  post.Slug,
				Title:     post.Title,
				Content:   post.Description,
				CreatedAt: post.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
				User: &entity.User{
					ID:    post.UserID,
					Name:  post.UserName,
					Email: post.UserEmail,
				},
			})
		}

		// Get unmoderated comments
		var comments []*dbModerationComment

		err = trx.Select(&comments, `
			SELECT c.id, c.post_id, p.number as post_number, p.slug as post_slug, c.content, c.created_at,
				   u.id as user_id, u.name as user_name, u.email as user_email,
				   p.title as post_title
			FROM comments c
			INNER JOIN users u ON u.id = c.user_id AND u.tenant_id = c.tenant_id
			INNER JOIN posts p ON p.id = c.post_id AND p.tenant_id = c.tenant_id
			WHERE c.tenant_id = $1 AND c.is_approved = false
			ORDER BY c.created_at DESC`, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get unmoderated comments")
		}

		for _, comment := range comments {
			q.Result = append(q.Result, &query.ModerationItem{
				Type:       "comment",
				ID:         comment.ID,
				PostID:     comment.PostID,
				PostNumber: comment.PostNumber,
				PostSlug:   comment.PostSlug,
				Content:    comment.Content,
				CreatedAt:  comment.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
				PostTitle:  comment.PostTitle,
				User: &entity.User{
					ID:    comment.UserID,
					Name:  comment.UserName,
					Email: comment.UserEmail,
				},
			})
		}

		return nil
	})
}