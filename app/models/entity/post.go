package entity

import (
	"fmt"
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

// CanBeVoted returns true if this post can have its vote changed
func (i *Post) CanBeVoted() bool {
	return i.Status != enum.PostCompleted && i.Status != enum.PostDeclined && i.Status != enum.PostDuplicate
}

func (i *Post) Url(baseURL string) string {
	return fmt.Sprintf("%s/posts/%d/%s", baseURL, i.Number, i.Slug)
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

func (i *OriginalPost) Url(baseURL string) string {
	return fmt.Sprintf("%s/posts/%d/%s", baseURL, i.Number, i.Slug)
}
