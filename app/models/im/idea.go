package im

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// Idea represents an unsaved idea
type Idea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *Idea) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *Idea) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	input.Title = strings.Trim(input.Title, " ")

	if input.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Title) < 10 || len(strings.Split(input.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	return result
}

// NewComment represents a new comment
type NewComment struct {
	Content string `json:"content"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (c *NewComment) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (c *NewComment) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if strings.Trim(c.Content, " ") == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	return result
}

// SetResponse represents the action to update an idea response
type SetResponse struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *SetResponse) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *SetResponse) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Status < models.IdeaNew || input.Status > models.IdeaDeclined {
		result.AddFieldFailure("status", "Status is invalid.")
	}

	if strings.Trim(input.Text, " ") == "" {
		result.AddFieldFailure("text", "Text is required.")
	}

	return result
}
