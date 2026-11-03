package query

import (
	"strconv"

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
	Query            string
	View             string
	Limit            string
	Statuses         []enum.PostStatus
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

// MaxSearchPostsLimit is the largest number of posts that can be requested over HTTP.
// There is no offset parameter: the home page and roadmap "show more" links work by
// increasing the limit, so this is deliberately generous.
const MaxSearchPostsLimit = 1000

// SetLimitFromString sets Limit from an untrusted value such as a query string parameter.
// Invalid or non-positive values fall back to the default limit. Unless unlimited is true,
// larger values are clamped to MaxSearchPostsLimit and "all" means MaxSearchPostsLimit
// rather than no limit at all.
func (q *SearchPosts) SetLimitFromString(limit string, unlimited bool) {
	if limit == "all" {
		if unlimited {
			q.Limit = "all"
		} else {
			q.Limit = strconv.Itoa(MaxSearchPostsLimit)
		}
		return
	}

	n, err := strconv.Atoi(limit)
	switch {
	case err != nil || n <= 0:
		q.Limit = ""
	case n > MaxSearchPostsLimit && !unlimited:
		q.Limit = strconv.Itoa(MaxSearchPostsLimit)
	default:
		q.Limit = strconv.Itoa(n)
	}
}

func (q *SearchPosts) SetStatusesFromStrings(statuses []string) {
	for _, v := range statuses {
		var postStatus enum.PostStatus
		if err := postStatus.UnmarshalText([]byte(v)); err == nil {
			q.Statuses = append(q.Statuses, postStatus)
		}
	}
}
