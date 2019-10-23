package query

import "github.com/getfider/fider/app/models"

type GetCommentByID struct {
	CommentID int

	Result *models.Comment
}

type GetCommentsByPost struct {
	Post *models.Post

	Result []*models.Comment
}
