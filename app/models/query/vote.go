package query

import "github.com/getfider/fider/app/models"

type ListPostVotes struct {
	PostID int
	Limit  int

	Result []*models.Vote
}
