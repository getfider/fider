package query

import (
	"github.com/getfider/fider/app/models/entity"
)

// GetDraftTags is used to get tags for a draft post
type GetDraftTags struct {
	DraftPost *entity.DraftPost
	Result    []*entity.Tag
}
