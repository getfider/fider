package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbDraftPost struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Attachments []string  `db:"attachments"`
}

func (i *dbDraftPost) toModel() *entity.DraftPost {
	return &entity.DraftPost{
		ID:          i.ID,
		Code:        i.Code,
		Title:       i.Title,
		Description: i.Description,
		CreatedAt:   i.CreatedAt,
	}
}

func getDraftPostByCode(ctx context.Context, q *query.GetDraftPostByCode) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		draftPost := dbDraftPost{}
		err := trx.Get(&draftPost, `
			SELECT id, code, title, description, created_at
			FROM draft_posts
			WHERE code = $1
		`, q.Code)
		if err != nil {
			return errors.Wrap(err, "failed to get draft post with code '%s'", q.Code)
		}
		q.Result = draftPost.toModel()
		return nil
	})
}

func getDraftAttachments(ctx context.Context, q *query.GetDraftAttachments) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]string, 0)

		postID := q.DraftPost.ID

		type entry struct {
			BlobKey string `db:"attachment_bkey"`
		}

		entries := []*entry{}
		err := trx.Select(&entries, `
			SELECT attachment_bkey
			FROM draft_attachments
			WHERE draft_post_id = $1
		`, postID)
		if err != nil {
			return errors.Wrap(err, "failed to get attachments")
		}

		q.Result = make([]string, len(entries))
		for i, entry := range entries {
			q.Result[i] = entry.BlobKey
		}

		return nil
	})
}
