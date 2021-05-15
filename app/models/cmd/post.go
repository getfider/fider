package cmd

import (
	"github.com/getfider/fider/app/models/entities"
	"github.com/getfider/fider/app/models/enum"
)

type AddNewPost struct {
	Title       string
	Description string

	Result *entities.Post
}

type UpdatePost struct {
	Post        *entities.Post
	Title       string
	Description string

	Result *entities.Post
}

type SetPostResponse struct {
	Post   *entities.Post
	Text   string
	Status enum.PostStatus
}
