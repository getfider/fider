package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

func setPostResponse(ctx context.Context, c *cmd.SetPostResponse) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if c.Status == models.PostDuplicate {
			return errors.New("Use MarkAsDuplicate to change an post status to Duplicate")
		}

		respondedAt := time.Now()
		if c.Post.Status == c.Status && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		_, err := trx.Execute(`
		UPDATE posts 
		SET response = $3, original_id = NULL, response_date = $4, response_user_id = $5, status = $6 
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Text, respondedAt, user.ID, c.Status)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = c.Status
		c.Post.Response = &models.PostResponse{
			Text:        c.Text,
			RespondedAt: respondedAt,
			User:        user,
		}
		return nil
	})
}

func markPostAsDuplicate(ctx context.Context, c *cmd.MarkPostAsDuplicate) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		respondedAt := time.Now()
		if c.Post.Status == models.PostDuplicate && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		var users []*dbUser
		err := trx.Select(&users, "SELECT user_id AS id FROM post_votes WHERE post_id = $1 AND tenant_id = $2", c.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get votes of post with id '%d'", c.Post.ID)
		}

		for _, u := range users {
			err := bus.Dispatch(ctx, &cmd.AddVote{Post: c.Original, User: u.toModel(ctx)})
			if err != nil {
				return err
			}
		}

		_, err = trx.Execute(`
		UPDATE posts 
		SET response = '', original_id = $3, response_date = $4, response_user_id = $5, status = $6 
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Original.ID, respondedAt, user.ID, models.PostDuplicate)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = models.PostDuplicate
		c.Post.Response = &models.PostResponse{
			RespondedAt: respondedAt,
			User:        user,
			Original: &models.OriginalPost{
				Number: c.Original.Number,
				Title:  c.Original.Title,
				Slug:   c.Original.Slug,
				Status: c.Original.Status,
			},
		}
		return nil
	})
}
