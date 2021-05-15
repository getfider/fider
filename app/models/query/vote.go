package query

import "github.com/getfider/fider/app/models/entity"

type ListPostVotes struct {
	PostID       int
	Limit        int
	IncludeEmail bool

	Result []*entity.Vote
}
