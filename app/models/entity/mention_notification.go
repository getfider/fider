package entity

import (
	"time"
)

// NotificationLog represents a record of a notification that was sent to a user
type MentionNotification struct {
	ID        int       `json:"id" db:"id"`
	TenantID  int       `json:"-" db:"tenant_id"`
	UserID    int       `json:"userId" db:"user_id"`
	CommentID int       `json:"commentId,omitempty" db:"comment_id"`
	PostID    int       `json:"postId,omitempty" db:"post_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_on"`
}
