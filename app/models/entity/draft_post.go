package entity

import (
	"time"
)

// DraftPost represents an anonymous draft post
type DraftPost struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdOn"`
}
