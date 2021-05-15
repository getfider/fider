package query

import "github.com/getfider/fider/app/models/entities"

type ListPostVotes struct {
	PostID       int
	Limit        int
	IncludeEmail bool

	Result []*entities.Vote
}
