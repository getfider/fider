package models

import (
	"time"

	"github.com/getfider/fider/app/models/enum"
)

//Post represents an post on a tenant board
type Post struct {
	ID            int             `json:"id"`
	Number        int             `json:"number"`
	Title         string          `json:"title"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"createdAt"`
	User          *User           `json:"user"`
	HasVoted      bool            `json:"hasVoted"`
	VotesCount    int             `json:"votesCount"`
	CommentsCount int             `json:"commentsCount"`
	Status        enum.PostStatus `json:"status"`
	Response      *PostResponse   `json:"response,omitempty"`
	Tags          []string        `json:"tags"`
}

//VoteUser represents a user that voted on a post
type VoteUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatarURL,omitempty"`
}

//Vote represents a vote given by a user on a post
type Vote struct {
	User      *VoteUser `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

// CanBeVoted returns true if this post can have its vote changed
func (i *Post) CanBeVoted() bool {
	return i.Status != enum.PostCompleted && i.Status != enum.PostDeclined && i.Status != enum.PostDuplicate
}

// NewPost represents a new post
type NewPost struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Attachments []*ImageUpload `json:"attachments"`
}

// UpdatePost represents a request to edit an existing post
type UpdatePost struct {
	Number      int            `route:"number"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Attachments []*ImageUpload `json:"attachments"`
}

// DeletePost represents a request to delete an existing post
type DeletePost struct {
	Number int    `route:"number"`
	Text   string `json:"text"`
}

// NewComment represents a new comment
type NewComment struct {
	Number      int            `route:"number"`
	Content     string         `json:"content"`
	Attachments []*ImageUpload `json:"attachments"`
}

// EditComment represents a request to edit existing comment
type EditComment struct {
	PostNumber  int            `route:"number"`
	ID          int            `route:"id"`
	Content     string         `json:"content"`
	Attachments []*ImageUpload `json:"attachments"`
}

// DeleteComment represents a request to delete an existing comment
type DeleteComment struct {
	PostNumber int `route:"number"`
	CommentID  int `route:"id"`
}

// SetResponse represents the action to update an post response
type SetResponse struct {
	Number         int             `route:"number"`
	Status         enum.PostStatus `json:"status"`
	Text           string          `json:"text"`
	OriginalNumber int             `json:"originalNumber"`
}

//PostResponse is a staff response to a given post
type PostResponse struct {
	Text        string        `json:"text"`
	RespondedAt time.Time     `json:"respondedAt"`
	User        *User         `json:"user"`
	Original    *OriginalPost `json:"original"`
}

//OriginalPost holds details of the original post of a duplicate
type OriginalPost struct {
	Number int             `json:"number"`
	Title  string          `json:"title"`
	Slug   string          `json:"slug"`
	Status enum.PostStatus `json:"status"`
}

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

//Tag represents a simple tag
type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Color    string `json:"color"`
	IsPublic bool   `json:"isPublic"`
}

//CreateEditTag is used to create a new tag or edit existing
type CreateEditTag struct {
	Slug     string `route:"slug"`
	Name     string `json:"name"`
	Color    string `json:"color" format:"upper"`
	IsPublic bool   `json:"isPublic"`
}

// DeleteTag is used to delete an existing tag
type DeleteTag struct {
	Slug string `route:"slug"`
}

// AssignUnassignTag is used to assign or remove a tag to/from an post
type AssignUnassignTag struct {
	Slug   string `route:"slug"`
	Number int    `route:"number"`
}
