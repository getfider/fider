package cmd

import (
	"github.com/getfider/fider/app/models/entities"
)

type AddNewTag struct {
	Name     string
	Color    string
	IsPublic bool

	Result *entities.Tag
}

type UpdateTag struct {
	TagID    int
	Name     string
	Color    string
	IsPublic bool

	Result *entities.Tag
}

type DeleteTag struct {
	Tag *entities.Tag
}

type AssignTag struct {
	Tag  *entities.Tag
	Post *entities.Post
}

type UnassignTag struct {
	Tag  *entities.Tag
	Post *entities.Post
}
