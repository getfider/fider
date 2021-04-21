package query

import "github.com/getfider/fider/app/models"

type ListPostVotes struct {
	PostID       int
	Limit        int
	IncludeEmail bool

	Result []*models.Vote
}
