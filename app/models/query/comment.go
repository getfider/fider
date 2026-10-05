package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetCommentByID struct {
	CommentID int
	// PostID, when set, restricts the lookup to comments on that post
	PostID int

	Result *entity.Comment
}

type GetCommentsByPost struct {
	Post *entity.Post

	Result []*entity.Comment
}
