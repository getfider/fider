package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

type AddNewComment struct {
	Post    *entity.Post
	Content string

	Result *entity.Comment
}

type UpdateComment struct {
	CommentID int
	Content   string
}

type DeleteComment struct {
	CommentID int
}
