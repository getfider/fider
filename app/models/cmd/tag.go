package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

type AddNewTag struct {
	Name     string
	Color    string
	IsPublic bool

	Result *entity.Tag
}

type UpdateTag struct {
	TagID    int
	Name     string
	Color    string
	IsPublic bool

	Result *entity.Tag
}

type DeleteTag struct {
	Tag *entity.Tag
}

type AssignTag struct {
	Tag  *entity.Tag
	Post *entity.Post
}

type UnassignTag struct {
	Tag  *entity.Tag
	Post *entity.Post
}
