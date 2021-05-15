package query

import (
	"github.com/getfider/fider/app/models/entities"
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

	Result *entities.Post
}

type GetPostBySlug struct {
	Slug string

	Result *entities.Post
}

type GetPostByNumber struct {
	Number int

	Result *entities.Post
}

type SearchPosts struct {
	Query string
	View  string
	Limit string
	Tags  []string

	Result []*entities.Post
}

type GetAllPosts struct {
	Result []*entities.Post
}
