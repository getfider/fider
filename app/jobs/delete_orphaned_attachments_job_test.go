package jobs_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/jobs"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestDeleteOrphanedAttachmentsJob_Schedule_IsCorrect(t *testing.T) {
	RegisterT(t)

	job := &jobs.DeleteOrphanedAttachmentsJobHandler{}
	Expect(job.Schedule()).Equals("0 0 0 * * 0")
}

func TestDeleteOrphanedAttachmentsJob_ShouldProcessBlobs(t *testing.T) {
	RegisterT(t)

	// Mock the ListBlobs query to return some test blob keys
	bus.AddHandler(func(ctx context.Context, q *query.ListBlobs) error {
		q.Result = []string{"attachments/test1.png", "attachments/test2.png"}
		return nil
	})

	// Mock the IsAttachmentReferenced query
	bus.AddHandler(func(ctx context.Context, q *query.IsAttachmentReferenced) error {
		// test1.png is referenced, test2.png is not
		q.Result = q.BlobKey == "attachments/test1.png"
		return nil
	})

	// Track which blobs were deleted
	deletedBlobs := make([]string, 0)
	bus.AddHandler(func(ctx context.Context, c *cmd.DeleteBlob) error {
		deletedBlobs = append(deletedBlobs, c.Key)
		return nil
	})

	job := &jobs.DeleteOrphanedAttachmentsJobHandler{}
	err := job.Run(jobs.Context{
		Context: context.Background(),
	})

	Expect(err).IsNil()
	Expect(deletedBlobs).HasLen(1)
	Expect(deletedBlobs[0]).Equals("attachments/test2.png")
}
