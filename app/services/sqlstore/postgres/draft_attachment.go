package postgres

// func setDraftAttachments(ctx context.Context, c *cmd.SetDraftAttachments) error {
// 	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
// 		// First delete all existing attachments for this draft post
// 		_, err := trx.Execute(`
// 			DELETE FROM draft_attachments
// 			WHERE draft_post_id = $1
// 		`, c.DraftPost.ID)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to delete existing draft attachments")
// 		}

// 		// Then insert all new attachments
// 		for _, attachment := range c.Attachments {
// 			if attachment.Remove {
// 				continue
// 			}
// 			_, err := trx.Execute(`
// 				INSERT INTO draft_attachments (draft_post_id, attachment_bkey)
// 				VALUES ($1, $2)
// 			`, c.DraftPost.ID, attachment.BlobKey)
// 			if err != nil {
// 				return errors.Wrap(err, "failed to insert draft attachment")
// 			}
// 		}

// 		return nil
// 	})
// }
