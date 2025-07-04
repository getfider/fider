package postgres

import (
	"context"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

func setDraftTags(ctx context.Context, c *cmd.SetDraftTags) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// First delete all existing tags for this draft post
		_, err := trx.Execute(`
			DELETE FROM draft_post_tags
			WHERE post_id = $1
		`, c.DraftPost.ID)
		if err != nil {
			return errors.Wrap(err, "failed to delete existing draft post tags")
		}

		// Then insert all new tags
		for _, tag := range c.Tags {
			_, err := trx.Execute(`
				INSERT INTO draft_post_tags (tag_id, post_id)
				VALUES ($1, $2)
			`, tag.ID, c.DraftPost.ID)
			if err != nil {
				return errors.Wrap(err, "failed to insert draft post tag")
			}
		}

		return nil
	})
}

func getDraftTags(ctx context.Context, q *query.GetDraftTags) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Tag, 0)

		postID := q.DraftPost.ID

		tags, err := queryTags(trx, `
			SELECT t.id, t.name, t.slug, t.color, t.is_public
			FROM tags t
			INNER JOIN draft_post_tags dpt
			ON dpt.tag_id = t.id
			WHERE dpt.post_id = $1 AND t.tenant_id = $2
			ORDER BY t.name
		`, postID, tenant.ID)

		if err != nil {
			return errors.Wrap(err, "failed to get tags for draft post")
		}

		q.Result = tags
		return nil
	})
}
