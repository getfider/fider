package dbEntities

import (
	"context"
	"encoding/json"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/dbx"
)

type comment struct {
	ID             int            `db:"id"`
	Content        string         `db:"content"`
	CreatedAt      time.Time      `db:"created_at"`
	User           *User          `db:"user"`
	Attachments    []string       `db:"attachment_bkeys"`
	EditedAt       dbx.NullTime   `db:"edited_at"`
	EditedBy       *User          `db:"edited_by"`
	ReactionCounts dbx.NullString `db:"reaction_counts"`
	IsApproved     bool           `db:"is_approved"`
}

func (c *comment) toModel(ctx context.Context) *entity.Comment {
	comment := &entity.Comment{
		ID:          c.ID,
		Content:     c.Content,
		CreatedAt:   c.CreatedAt,
		User:        c.User.ToModel(ctx),
		Attachments: c.Attachments,
		IsApproved:  c.IsApproved,
	}
	if c.EditedAt.Valid {
		comment.EditedBy = c.EditedBy.ToModel(ctx)
		comment.EditedAt = &c.EditedAt.Time
	}

	if c.ReactionCounts.Valid {
		_ = json.Unmarshal([]byte(c.ReactionCounts.String), &comment.ReactionCounts)
	}
	return comment
}
