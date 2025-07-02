package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type ModerationItem struct {
	Type       string       `json:"type"` // "post" or "comment"
	ID         int          `json:"id"`
	PostID     int          `json:"postId,omitempty"`
	PostNumber int          `json:"postNumber,omitempty"`
	PostSlug   string       `json:"postSlug,omitempty"`
	Title      string       `json:"title,omitempty"`
	Content    string       `json:"content"`
	User       *entity.User `json:"user"`
	CreatedAt  string       `json:"createdAt"`
	PostTitle  string       `json:"postTitle,omitempty"`
}

type GetModerationItems struct {
	Result []*ModerationItem
}