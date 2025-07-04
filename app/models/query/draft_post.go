package query

import (
	"github.com/getfider/fider/app/models/entity"
)

// GetDraftPostByCode is used to get a draft post by its code
type GetDraftPostByCode struct {
	Code   string
	Result *entity.DraftPost
}

// GetDraftAttachments is used to get attachments for a draft post
type GetDraftAttachments struct {
	DraftPost *entity.DraftPost
	Result    []string
}
