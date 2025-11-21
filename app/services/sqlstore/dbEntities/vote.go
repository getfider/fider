package dbEntities

import (
	"context"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type Vote struct {
	User *struct {
		ID            int    `db:"id"`
		Name          string `db:"name"`
		Email         string `db:"email"`
		AvatarType    int64  `db:"avatar_type"`
		AvatarBlobKey string `db:"avatar_bkey"`
	} `db:"user"`
	CreatedAt time.Time `db:"created_at"`
}

func (v *Vote) ToModel(ctx context.Context) *entity.Vote {
	vote := &entity.Vote{
		CreatedAt: v.CreatedAt,
		User: &entity.VoteUser{
			ID:        v.User.ID,
			Name:      v.User.Name,
			Email:     v.User.Email,
			AvatarURL: buildAvatarURL(ctx, enum.AvatarType(v.User.AvatarType), v.User.ID, v.User.Name, v.User.AvatarBlobKey),
		},
	}
	return vote
}
