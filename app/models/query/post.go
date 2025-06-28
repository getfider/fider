package query

import (
	"github.com/getfider/fider/app/models/entity"
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

	Result *entity.Post
}

type GetPostBySlug struct {
	Slug string

	Result *entity.Post
}

type GetPostByNumber struct {
	Number int

	Result *entity.Post
}

type SearchPosts struct {
	Query       string
	View        string
	Limit       string
	Statuses    []enum.PostStatus
	Tags        []string
	MyVotesOnly bool

	Result []*entity.Post
}

type FindSimilarPosts struct {
	Query string

	Result []*entity.Post
}

type GetAllPosts struct {
	Result []*entity.Post
}

func (q *SearchPosts) SetStatusesFromStrings(statuses []string) {
	for _, v := range statuses {
		var postStatus enum.PostStatus
		if err := postStatus.UnmarshalText([]byte(v)); err == nil {
			q.Statuses = append(q.Statuses, postStatus)
		}
	}
}
