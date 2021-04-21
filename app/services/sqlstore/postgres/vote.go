package postgres

import (
	"context"
	"strconv"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbVote struct {
	User *struct {
		ID            int    `db:"id"`
		Name          string `db:"name"`
		Email         string `db:"email"`
		AvatarType    int64  `db:"avatar_type"`
		AvatarBlobKey string `db:"avatar_bkey"`
	} `db:"user"`
	CreatedAt time.Time `db:"created_at"`
}

func (v *dbVote) toModel(ctx context.Context) *models.Vote {
	vote := &models.Vote{
		CreatedAt: v.CreatedAt,
		User: &models.VoteUser{
			ID:        v.User.ID,
			Name:      v.User.Name,
			Email:     v.User.Email,
			AvatarURL: buildAvatarURL(ctx, enum.AvatarType(v.User.AvatarType), v.User.ID, v.User.Name, v.User.AvatarBlobKey),
		},
	}
	return vote
}

func addVote(ctx context.Context, c *cmd.AddVote) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if !c.Post.CanBeVoted() {
			return nil
		}

		_, err := trx.Execute(
			`INSERT INTO post_votes (tenant_id, user_id, post_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
			tenant.ID, c.User.ID, c.Post.ID, time.Now(),
		)

		if err != nil {
			return errors.Wrap(err, "failed add vote to post")
		}

		return internalAddSubscriber(trx, c.Post, tenant, c.User, false)
	})
}

func removeVote(ctx context.Context, c *cmd.RemoveVote) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		if !c.Post.CanBeVoted() {
			return nil
		}

		_, err := trx.Execute(`DELETE FROM post_votes WHERE user_id = $1 AND post_id = $2 AND tenant_id = $3`, c.User.ID, c.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to remove vote from post")
		}

		return err
	})
}

func listPostVotes(ctx context.Context, q *query.ListPostVotes) error {
	return using(ctx, func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error {
		q.Result = make([]*models.Vote, 0)
		sqlLimit := "ALL"
		if q.Limit > 0 {
			sqlLimit = strconv.Itoa(q.Limit)
		}

		emailColumn := "''"
		if q.IncludeEmail {
			emailColumn = "u.email"
		}

		votes := []*dbVote{}
		err := trx.Select(&votes, `
		SELECT 
			pv.created_at, 
			u.id AS user_id,
			u.name AS user_name,
			`+emailColumn+` AS user_email,
			u.avatar_type AS user_avatar_type,
			u.avatar_bkey AS user_avatar_bkey
		FROM post_votes pv
		INNER JOIN users u
		ON u.id = pv.user_id
		AND u.tenant_id = pv.tenant_id 
		WHERE pv.post_id = $1  
		AND pv.tenant_id = $2
		ORDER BY pv.created_at
		LIMIT `+sqlLimit, q.PostID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get votes of post")
		}

		q.Result = make([]*models.Vote, len(votes))
		for i, vote := range votes {
			q.Result[i] = vote.toModel(ctx)
		}

		return nil
	})
}
