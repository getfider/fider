package entity

import "time"

type ReactionCounts struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

// Comment represents an user comment on an post
type Comment struct {
	ID             int            `json:"id"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"createdAt"`
	User           *User          `json:"user"`
	Attachments    []string       `json:"attachments,omitempty"`
	EditedAt       *time.Time     `json:"editedAt,omitempty"`
	EditedBy       *User          `json:"editedBy,omitempty"`
	ReactionCounts map[string]int `json:"reactionCounts,omitempty"`
	Reactions      string         `json:"reactions,omitempty"`
}
