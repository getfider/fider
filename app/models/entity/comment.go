package entity

import (
	"time"
)

type ReactionCounts struct {
	Emoji      string `json:"emoji"`
	Count      int    `json:"count"`
	IncludesMe bool   `json:"includesMe"`
}

// Comment represents an user comment on an post
type Comment struct {
	ID             int              `json:"id"`
	Content        string           `json:"content"`
	CreatedAt      time.Time        `json:"createdAt"`
	User           *User            `json:"user"`
	Attachments    []string         `json:"attachments,omitempty"`
	EditedAt       *time.Time       `json:"editedAt,omitempty"`
	EditedBy       *User            `json:"editedBy,omitempty"`
	ReactionCounts []ReactionCounts `json:"reactionCounts,omitempty"`
	IsApproved     bool             `json:"isApproved"`
}
