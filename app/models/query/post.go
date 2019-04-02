package query

import "github.com/getfider/fider/app/models"

type PostIsReferenced struct {
	PostID int

	Result bool
}

type CountPostPerStatus struct {
	Result map[models.PostStatus]int
}
