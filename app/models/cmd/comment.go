package cmd

import (
	"github.com/getfider/fider/app/models/entities"
)

type AddNewComment struct {
	Post    *entities.Post
	Content string

	Result *entities.Comment
}

type UpdateComment struct {
	CommentID int
	Content   string
}

type DeleteComment struct {
	CommentID int
}
