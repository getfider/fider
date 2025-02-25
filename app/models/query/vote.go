package query

import "github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"

type ListPostVotes struct {
	PostID       int
	Limit        int
	IncludeEmail bool

	Result []*entity.Vote
}
