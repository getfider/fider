package entity

import "time"

// Reaction represents a user's emoji reaction to a comment
type Reaction struct {
	ID        int       `json:"id"`
	Emoji     string    `json:"emoji"`
	Comment   *Comment  `json:"-"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}
