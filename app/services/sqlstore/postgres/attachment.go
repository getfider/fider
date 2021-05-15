package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/services/blob"
)

func setAttachments(ctx context.Context, c *cmd.SetAttachments) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		postID := c.Post.ID
		var commentID sql.NullInt64
		if c.Comment != nil {
			err := commentID.Scan(c.Comment.ID)
			if err != nil {
				return errors.Wrap(err, "failed scan comment id")
			}
		}

		for _, attachment := range c.Attachments {
			if attachment.Remove {
				if _, err := trx.Execute(
					"DELETE FROM attachments WHERE tenant_id = $1 AND post_id = $2 AND (comment_id = $3 OR ($3 IS NULL AND comment_id IS NULL)) AND attachment_bkey = $4",
					tenant.ID, postID, commentID, attachment.BlobKey,
				); err != nil {
					return errors.Wrap(err, "failed to delete attachment")
				}
			} else {
				if _, err := trx.Execute(
					"INSERT INTO attachments (tenant_id, post_id, comment_id, user_id, attachment_bkey) VALUES ($1, $2, $3, $4, $5)",
					tenant.ID, postID, commentID, user.ID, attachment.BlobKey,
				); err != nil {
					return errors.Wrap(err, "failed to insert attachment")
				}
			}
		}

		return nil
	})
}

func getAttachments(ctx context.Context, q *query.GetAttachments) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]string, 0)

		postID := q.Post.ID
		var commentID sql.NullInt64
		if q.Comment != nil {
			err := commentID.Scan(q.Comment.ID)
			if err != nil {
				return errors.Wrap(err, "failed scan comment id")
			}
		}

		type entry struct {
			BlobKey string `db:"attachment_bkey"`
		}

		entries := []*entry{}
		err := trx.Select(&entries, `
			SELECT attachment_bkey
			FROM attachments
			WHERE tenant_id = $1 AND post_id = $2 AND (comment_id = $3 OR ($3 IS NULL AND comment_id IS NULL))
		`, tenant.ID, postID, commentID)
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

func uploadImage(ctx context.Context, c *cmd.UploadImage) error {
	if c.Image.Upload != nil && len(c.Image.Upload.Content) > 0 {
		bkey := fmt.Sprintf("%s/%s-%s", c.Folder, rand.String(64), blob.SanitizeFileName(c.Image.Upload.FileName))
		err := bus.Dispatch(ctx, &cmd.StoreBlob{
			Key:         bkey,
			Content:     c.Image.Upload.Content,
			ContentType: c.Image.Upload.ContentType,
		})
		if err != nil {
			return errors.Wrap(err, "failed to upload new blob")
		}
		c.Image.BlobKey = bkey
	}
	return nil
}

func uploadImages(ctx context.Context, c *cmd.UploadImages) error {
	for _, img := range c.Images {
		if err := bus.Dispatch(ctx, &cmd.UploadImage{
			Image:  img,
			Folder: c.Folder,
		}); err != nil {
			return err
		}
	}
	return nil
}
