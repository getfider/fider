package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type PostIsReferenced struct {
	PostID int

	Result bool
}

// CountPostPerStatus keyed by tenant status slug.
type CountPostPerStatus struct {
	Result map[string]int
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
	Query            string
	View             string
	Limit            string
	Statuses         []string
	Tags             []string
	MyVotesOnly      bool
	NoTagsOnly       bool
	MyPostsOnly      bool
	ModerationFilter string // "pending", "approved", or empty (all)

	Result []*entity.Post
}

type FindSimilarPosts struct {
	Query string

	Result []*entity.Post
}

type GetAllPosts struct {
	Result []*entity.Post
}

// SetStatusesFromStrings accepts the raw query-param strings the SearchPosts
// handler receives. Slugs are stored verbatim — tenant-defined custom slugs
// flow through alongside the built-ins.
func (q *SearchPosts) SetStatusesFromStrings(statuses []string) {
	for _, v := range statuses {
		if v == "" {
			continue
		}
		q.Statuses = append(q.Statuses, v)
	}
}
