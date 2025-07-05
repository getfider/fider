package jobs

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
)

type DeleteOrphanedAttachmentsJobHandler struct {
}

func (e DeleteOrphanedAttachmentsJobHandler) Schedule() string {
	// return "0 0 0 * * 0" // Run weekly at midnight on Sunday
	return "*/30 * * * * *" // Run every 30 seconds
}

func (e DeleteOrphanedAttachmentsJobHandler) Run(ctx Context) error {
	log.Debug(ctx, "checking for orphaned attachment blobs")

	// 1. List all blobs with the "attachments/" prefix
	listBlobsQuery := &query.ListBlobs{
		Prefix: "attachments/",
	}

	if err := bus.Dispatch(ctx, listBlobsQuery); err != nil {
		return err
	}

	blobKeys := listBlobsQuery.Result
	log.Debugf(ctx, "found @{BlobCount} attachment blobs to check", dto.Props{
		"BlobCount": len(blobKeys),
	})

	// 2. Check each blob against the attachments table and delete if orphaned
	deletedCount := 0

	for _, blobKey := range blobKeys {
		// Check if this blob is referenced in the attachments table
		isReferenced := &query.IsAttachmentReferenced{
			BlobKey: blobKey,
		}

		if err := bus.Dispatch(ctx, isReferenced); err != nil {
			log.Error(ctx, err)
			log.Errorf(ctx, "failed to check if blob '@{BlobKey}' is referenced", dto.Props{
				"BlobKey": blobKey,
			})
			continue
		}

		// If not referenced, delete it
		if !isReferenced.Result {
			deleteCmd := &cmd.DeleteBlob{
				Key: blobKey,
			}

			if err := bus.Dispatch(ctx, deleteCmd); err != nil {
				log.Error(ctx, err)
				log.Errorf(ctx, "failed to delete orphaned blob '@{BlobKey}'", dto.Props{
					"BlobKey": blobKey,
				})
				continue
			}

			deletedCount++
			log.Debugf(ctx, "deleted orphaned blob '@{BlobKey}'", dto.Props{
				"BlobKey": blobKey,
			})
		}
	}

	log.Debugf(ctx, "@{DeletedCount} orphaned attachment blobs were deleted", dto.Props{
		"DeletedCount": deletedCount,
	})

	return nil
}
