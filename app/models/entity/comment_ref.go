package entity

import (
	"time"
)

type CommentRef struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UserID    int        `json:"userId"`
	PostID    int        `json:"postId"`
	EditedAt  *time.Time `json:"editedAt,omitempty"`
}
