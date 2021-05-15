package entity

import "time"

//Comment represents an user comment on an post
type Comment struct {
	ID          int        `json:"id"`
	Content     string     `json:"content"`
	CreatedAt   time.Time  `json:"createdAt"`
	User        *User      `json:"user"`
	Attachments []string   `json:"attachments,omitempty"`
	EditedAt    *time.Time `json:"editedAt,omitempty"`
	EditedBy    *User      `json:"editedBy,omitempty"`
}
