package entity

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

// Notification is the system generated notification entity
type Notification struct {
	ID            int             `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Link          string          `json:"link" db:"link"`
	Read          bool            `json:"read" db:"read"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
	AuthorName    string          `json:"authorName" db:"name"`
	AuthorID      int             `json:"-" db:"author_id"`
	AvatarBlobKey string          `json:"-" db:"avatar_bkey"`
	AvatarType    enum.AvatarType `json:"-" db:"avatar_type"`
	AvatarURL     string          `json:"avatarURL,omitempty"`
}
