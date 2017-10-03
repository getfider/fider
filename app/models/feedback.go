package models

import (
	"time"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID              int           `json:"id"`
	Number          int           `json:"number"`
	Title           string        `json:"title"`
	Slug            string        `json:"slug"`
	Description     string        `json:"description"`
	CreatedOn       time.Time     `json:"createdOn"`
	User            *User         `json:"user"`
	TotalSupporters int           `json:"totalSupporters"`
	TotalComments   int           `json:"totalComments"`
	Status          int           `json:"status"`
	Response        *IdeaResponse `json:"response"`
}

//CanBeChangedBy returns true if given user can change this idea
func (i *Idea) CanBeChangedBy(user *User) bool {
	return user.IsStaff() || i.User.ID == user.ID
}

// NewIdea represents a new idea
type NewIdea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NewComment represents a new comment
type NewComment struct {
	Number  int    `route:"number"`
	Content string `json:"content"`
}

// SetResponse represents the action to update an idea response
type SetResponse struct {
	Number int    `route:"number"`
	Status int    `json:"status"`
	Text   string `json:"text"`
}

//IdeaResponse is a staff response to a given idea
type IdeaResponse struct {
	Text        string    `json:"text"`
	RespondedOn time.Time `json:"respondedOn"`
	User        *User     `json:"user"`
}

//Comment represents an user comment on an idea
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedOn time.Time `json:"createdOn"`
	User      *User     `json:"user"`
}

var (
	//IdeaNew is the default status
	IdeaNew = 0
	//IdeaStarted is used when the idea has been accepted and work is in progress
	IdeaStarted = 1
	//IdeaCompleted is used when the idea has been accepted and already implemented
	IdeaCompleted = 2
	//IdeaDeclined is used when organizers decide to decline an idea
	IdeaDeclined = 3
)
