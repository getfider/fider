package models

import (
	"time"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID              int        `json:"id" db:"id"`
	Number          int        `json:"number" db:"number"`
	Title           string     `json:"title" db:"title"`
	Slug            string     `json:"slug" db:"slug"`
	Description     string     `json:"description" db:"description"`
	CreatedOn       time.Time  `json:"createdOn" db:"created_on"`
	User            *User      `json:"user" db:"user"`
	TotalSupporters int        `json:"totalSupporters" db:"supporters"`
	Status          IdeaStatus `json:"status" db:"status"`
}

//IdeaStatus represents the status of an idea
type IdeaStatus int

var (
	//IdeaNew is the default status
	IdeaNew = IdeaStatus(0)
	//IdeaStarted is used when the idea has been accepted and work is in progress
	IdeaStarted = IdeaStatus(1)
	//IdeaCompleted is used when the idea has been accepted and already implemented
	IdeaCompleted = IdeaStatus(2)
	//IdeaDeclined is used when organizers decide to decline an idea
	IdeaDeclined = IdeaStatus(3)
)

//Comment represents an user comment on an idea
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	CreatedOn time.Time `json:"createdOn" db:"created_on"`
	User      *User     `json:"user" db:"user"`
}
