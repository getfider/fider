package actions

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateNewIdea is used to create a new idea
type CreateNewIdea struct {
	Model *models.NewIdea
}

// NewModel initializes the model
func (input *CreateNewIdea) NewModel() interface{} {
	input.Model = new(models.NewIdea)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateNewIdea) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *CreateNewIdea) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	input.Model.Title = strings.Trim(input.Model.Title, " ")

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Model.Title) < 10 || len(strings.Split(input.Model.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	return result
}

// AddNewComment represents a new comment to be added
type AddNewComment struct {
	Model *models.NewComment
}

// NewModel initializes the model
func (input *AddNewComment) NewModel() interface{} {
	input.Model = new(models.NewComment)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *AddNewComment) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *AddNewComment) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if strings.Trim(input.Model.Content, " ") == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	return result
}

// SetResponse represents the action to update an idea response
type SetResponse struct {
	Model *models.SetResponse
}

// NewModel initializes the model
func (input *SetResponse) NewModel() interface{} {
	input.Model = new(models.SetResponse)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *SetResponse) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (input *SetResponse) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Status < models.IdeaNew || input.Model.Status > models.IdeaDeclined {
		result.AddFieldFailure("status", "Status is invalid.")
	}

	if strings.Trim(input.Model.Text, " ") == "" {
		result.AddFieldFailure("text", "Text is required.")
	}

	return result
}
