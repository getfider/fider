package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetTagBySlug struct {
	Slug string

	Result *entity.Tag
}

type GetAssignedTags struct {
	Post *entity.Post

	Result []*entity.Tag
}

type GetAllTags struct {
	Result []*entity.Tag
}
