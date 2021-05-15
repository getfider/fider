package query

import (
	"github.com/getfider/fider/app/models/entities"
)

type GetCommentByID struct {
	CommentID int

	Result *entities.Comment
}

type GetCommentsByPost struct {
	Post *entities.Post

	Result []*entities.Comment
}
