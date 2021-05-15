package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetCommentByID struct {
	CommentID int

	Result *entity.Comment
}

type GetCommentsByPost struct {
	Post *entity.Post

	Result []*entity.Comment
}
