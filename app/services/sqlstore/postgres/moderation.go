package postgres

import (
	"context"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/errors"
)

// Stub functions - commercial license required for content moderation

func approvePost(ctx context.Context, c *cmd.ApprovePost) error {
	return errors.New("Content moderation requires commercial license")
}

func declinePost(ctx context.Context, c *cmd.DeclinePost) error {
	return errors.New("Content moderation requires commercial license")
}

func approveComment(ctx context.Context, c *cmd.ApproveComment) error {
	return errors.New("Content moderation requires commercial license")
}

func declineComment(ctx context.Context, c *cmd.DeclineComment) error {
	return errors.New("Content moderation requires commercial license")
}

func bulkApproveItems(ctx context.Context, c *cmd.BulkApproveItems) error {
	return errors.New("Content moderation requires commercial license")
}

func bulkDeclineItems(ctx context.Context, c *cmd.BulkDeclineItems) error {
	return errors.New("Content moderation requires commercial license")
}

func getModerationItems(ctx context.Context, q *query.GetModerationItems) error {
	q.Result = make([]*query.ModerationItem, 0)
	return nil
}

func getModerationCount(ctx context.Context, q *query.GetModerationCount) error {
	q.Result = 0
	return nil
}

func trustUser(ctx context.Context, c *cmd.TrustUser) error {
	return errors.New("Content moderation requires commercial license")
}