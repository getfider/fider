package query

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/enum"
)

type PostIsReferenced struct {
	PostID int

	Result bool
}

type CountPostPerStatus struct {
	Result map[enum.PostStatus]int
}

type GetPostByID struct {
	PostID int

	Result *models.Post
}

type GetPostBySlug struct {
	Slug string

	Result *models.Post
}

type GetPostByNumber struct {
	Number int

	Result *models.Post
}

type SearchPosts struct {
	Query string
	View  string
	Limit string
	Tags  []string

	Result []*models.Post
}

type GetAllPosts struct {
	Result []*models.Post
}
