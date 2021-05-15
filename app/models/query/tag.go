package query

import (
	"github.com/getfider/fider/app/models/entities"
)

type GetTagBySlug struct {
	Slug string

	Result *entities.Tag
}

type GetAssignedTags struct {
	Post *entities.Post

	Result []*entities.Tag
}

type GetAllTags struct {
	Result []*entities.Tag
}
