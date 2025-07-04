package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

type SetDraftTags struct {
	DraftPost *entity.DraftPost
	Tags      []*entity.Tag
}
