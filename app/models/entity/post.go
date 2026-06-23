package entity

import (
	"fmt"
	"time"
)

//Post represents an post on a tenant board
type Post struct {
	ID            int           `json:"id"`
	Number        int           `json:"number"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
	CreatedAt     time.Time     `json:"createdAt"`
	User          *User         `json:"user"`
	HasVoted      bool          `json:"hasVoted"`
	VotesCount    int           `json:"votesCount"`
	CommentsCount int           `json:"commentsCount"`
	StatusSlug    string        `json:"status"`
	StatusKind    string        `json:"statusKind"`
	Response      *PostResponse `json:"response,omitempty"`
	Tags          []string      `json:"tags"`
	IsApproved    bool          `json:"isApproved"`
}

// CanBeVoted returns true if this post can have its vote changed.
// Kind-based check covers all custom statuses an admin marked as closed.
func (i *Post) CanBeVoted() bool {
	return i.StatusKind != "closed-completed" && i.StatusKind != "closed-declined" && i.StatusKind != "duplicate"
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
	Number     int    `json:"number"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	StatusSlug string `json:"status"`
}

func (i *OriginalPost) Url(baseURL string) string {
	return fmt.Sprintf("%s/posts/%d/%s", baseURL, i.Number, i.Slug)
}
