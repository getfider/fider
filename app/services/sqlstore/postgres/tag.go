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
	"github.com/gosimple/slug"
)

type dbTag struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Slug     string `db:"slug"`
	Color    string `db:"color"`
	IsPublic bool   `db:"is_public"`
}

func (t *dbTag) toModel() *entity.Tag {
	return &entity.Tag{
		ID:       t.ID,
		Name:     t.Name,
		Slug:     t.Slug,
		Color:    t.Color,
		IsPublic: t.IsPublic,
	}
}

func getTagBySlug(ctx context.Context, q *query.GetTagBySlug) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		tag, err := queryTagBySlug(trx, tenant, q.Slug)
		q.Result = tag
		return err
	})
}

func getAssignedTags(ctx context.Context, q *query.GetAssignedTags) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Tag, 0)

		tags, err := queryTags(trx, `
			SELECT t.id, t.name, t.slug, t.color, t.is_public 
			FROM tags t
			INNER JOIN post_tags pt
			ON pt.tag_id = t.id
			AND pt.tenant_id = t.tenant_id
			WHERE pt.post_id = $1 AND t.tenant_id = $2
			ORDER BY t.name
		`, q.Post.ID, tenant.ID)

		if err != nil {
			return errors.Wrap(err, "failed get assigned tags")
		}

		q.Result = tags
		return nil
	})
}

func getAllTags(ctx context.Context, q *query.GetAllTags) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Tag, 0)

		condition := `AND t.is_public = true`
		if user != nil && user.IsCollaborator() {
			condition = ``
		}

		query := fmt.Sprintf(`
			SELECT t.id, t.name, t.slug, t.color, t.is_public 
			FROM tags t
			WHERE t.tenant_id = $1 %s
			ORDER BY t.name
		`, condition)
		tags, err := queryTags(trx, query, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed get all tags")
		}

		q.Result = tags
		return nil
	})
}

func addNewTag(ctx context.Context, c *cmd.AddNewTag) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		c.Result = nil
		newSlug := slug.Make(c.Name)

		_, err := trx.Execute(`
			INSERT INTO tags (name, slug, color, is_public, created_at, tenant_id) 
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
		`, c.Name, newSlug, c.Color, c.IsPublic, time.Now(), tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to add new tag")
		}

		tag, err := queryTagBySlug(trx, tenant, newSlug)
		c.Result = tag
		return err
	})
}

func updateTag(ctx context.Context, c *cmd.UpdateTag) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		c.Result = nil
		newSlug := slug.Make(c.Name)

		_, err := trx.Execute(`UPDATE tags SET name = $1, slug = $2, color = $3, is_public = $4
													 WHERE id = $5 AND tenant_id = $6`, c.Name, newSlug, c.Color, c.IsPublic, c.TagID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to update tag")
		}

		tag, err := queryTagBySlug(trx, tenant, newSlug)
		c.Result = tag
		return err
	})
}

func deleteTag(ctx context.Context, c *cmd.DeleteTag) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`DELETE FROM post_tags WHERE tag_id = $1 AND tenant_id = $2`, c.Tag.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to remove tag with id '%d' from all posts", c.Tag.ID)
		}

		_, err = trx.Execute(`DELETE FROM tags WHERE id = $1 AND tenant_id = $2`, c.Tag.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to delete tag with id '%d'", c.Tag.ID)
		}
		return nil
	})
}

func assignTag(ctx context.Context, c *cmd.AssignTag) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		alreadyAssigned, err := trx.Exists("SELECT 1 FROM post_tags WHERE post_id = $1 AND tag_id = $2 AND tenant_id = $3", c.Post.ID, c.Tag.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to check if tag is already assigned")
		}

		if alreadyAssigned {
			return nil
		}

		_, err = trx.Execute(
			`INSERT INTO post_tags (tag_id, post_id, created_at, created_by_id, tenant_id) VALUES ($1, $2, $3, $4, $5)`,
			c.Tag.ID, c.Post.ID, time.Now(), user.ID, tenant.ID,
		)

		if err != nil {
			return errors.Wrap(err, "failed to assign tag to post")
		}
		return nil
	})
}

func unassignTag(ctx context.Context, c *cmd.UnassignTag) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(
			`DELETE FROM post_tags WHERE tag_id = $1 AND post_id = $2 AND tenant_id = $3`,
			c.Tag.ID, c.Post.ID, tenant.ID,
		)

		if err != nil {
			return errors.Wrap(err, "failed to unassign tag from post")
		}
		return nil
	})
}

func queryTagBySlug(trx *dbx.Trx, tenant *entity.Tenant, slug string) (*entity.Tag, error) {
	tag := dbTag{}

	err := trx.Get(&tag, "SELECT id, name, slug, color, is_public FROM tags WHERE tenant_id = $1 AND slug = $2", tenant.ID, slug)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tag with slug '%s'", slug)
	}

	return tag.toModel(), nil
}

func queryTags(trx *dbx.Trx, query string, args ...interface{}) ([]*entity.Tag, error) {
	tags := []*dbTag{}
	err := trx.Select(&tags, query, args...)
	if err != nil {
		return nil, err
	}

	var result = make([]*entity.Tag, len(tags))
	for i, tag := range tags {
		result[i] = tag.toModel()
	}
	return result, nil
}
