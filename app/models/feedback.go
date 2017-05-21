package models

import (
	"encoding/json"
	"time"

	"github.com/getfider/fider/app/pkg/dbx"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID              int           `json:"id" db:"id"`
	Number          int           `json:"number" db:"number"`
	Title           string        `json:"title" db:"title"`
	Slug            string        `json:"slug" db:"slug"`
	Description     string        `json:"description" db:"description"`
	CreatedOn       time.Time     `json:"createdOn" db:"created_on"`
	User            *User         `json:"user" db:"user"`
	TotalSupporters int           `json:"totalSupporters" db:"supporters"`
	Status          IdeaStatus    `json:"status" db:"status"`
	Response        *IdeaResponse `json:"response" db:"response"`
}

//IdeaResponse is a staff response to a given idea
type IdeaResponse struct {
	Text      dbx.NullString `db:"text"`
	CreatedOn dbx.NullTime   `db:"date"`
	UserID    dbx.NullInt    `db:"user_id"`
}

// MarshalJSON interface redefinition
func (r IdeaResponse) MarshalJSON() ([]byte, error) {
	if r.Text.Valid || r.CreatedOn.Valid || r.UserID.Valid {
		return json.Marshal(struct {
			Text      dbx.NullString `json:"text"`
			CreatedOn dbx.NullTime   `json:"createdOn"`
			UserID    dbx.NullInt    `json:"userId"`
		}{
			Text:      r.Text,
			CreatedOn: r.CreatedOn,
			UserID:    r.UserID,
		})
	}
	return json.Marshal(nil)
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
