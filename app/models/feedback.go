package models

import (
	"time"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID          int       `json:"id" db:"id"`
	Number      int       `json:"number" db:"number"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedOn   time.Time `json:"createdOn" db:"created_on"`
	User        *User     `json:"user" db:"user"`
}

//Comment represents an user comment on an idea
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	CreatedOn time.Time `json:"createdOn" db:"created_on"`
	User      *User     `json:"user" db:"user"`
}
