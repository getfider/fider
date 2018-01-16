package actions

import (
	"strings"

	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateNewIdea is used to create a new idea
type CreateNewIdea struct {
	Model *models.NewIdea
}

// Initialize the model
func (input *CreateNewIdea) Initialize() interface{} {
	input.Model = new(models.NewIdea)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateNewIdea) IsAuthorized(user *models.User) bool {
	return user != nil
}

// Validate is current model is valid
func (input *CreateNewIdea) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must be less than 100 characters.")
	}

	if len(input.Model.Title) < 10 || len(strings.Split(input.Model.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	idea, err := services.Ideas.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && err != app.ErrNotFound {
		return validate.Error(err)
	} else if idea != nil {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	return result
}

// UpdateIdea is used to edit an existing new idea
type UpdateIdea struct {
	Model *models.UpdateIdea
}

// Initialize the model
func (input *UpdateIdea) Initialize() interface{} {
	input.Model = new(models.UpdateIdea)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateIdea) IsAuthorized(user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *UpdateIdea) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	idea, err := services.Ideas.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(input.Model.Title) < 10 || len(strings.Split(input.Model.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	another, err := services.Ideas.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && err != app.ErrNotFound {
		return validate.Error(err)
	} else if another != nil && another.ID != idea.ID {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	return result
}

// AddNewComment represents a new comment to be added
type AddNewComment struct {
	Model *models.NewComment
}

// Initialize the model
func (input *AddNewComment) Initialize() interface{} {
	input.Model = new(models.NewComment)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *AddNewComment) IsAuthorized(user *models.User) bool {
	return user != nil
}

// Validate is current model is valid
func (input *AddNewComment) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Content == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	return result
}

// SetResponse represents the action to update an idea response
type SetResponse struct {
	Model    *models.SetResponse
	Original *models.Idea
}

// Initialize the model
func (input *SetResponse) Initialize() interface{} {
	input.Model = new(models.SetResponse)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *SetResponse) IsAuthorized(user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *SetResponse) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Status < models.IdeaOpen || input.Model.Status > models.IdeaDuplicate {
		result.AddFieldFailure("status", "Status is invalid.")
	}

	if input.Model.Status == models.IdeaDuplicate {
		if input.Model.OriginalNumber == input.Model.Number {
			result.AddFieldFailure("originalNumber", "Cannot be a duplicate of itself")
		}

		original, err := services.Ideas.GetByNumber(input.Model.OriginalNumber)
		if err != nil {
			if err == app.ErrNotFound {
				result.AddFieldFailure("originalNumber", "Original idea not found")
			} else {
				return validate.Error(err)
			}
		}
		if original != nil {
			input.Original = original
		}
	} else if input.Model.Status != models.IdeaOpen && input.Model.Text == "" {
		result.AddFieldFailure("text", "Description is required.")
	}

	return result
}
