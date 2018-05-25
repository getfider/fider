package actions

import (
	"strings"

	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
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
func (input *CreateNewIdea) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil
}

// Validate is current model is valid
func (input *CreateNewIdea) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 || len(strings.Split(input.Model.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	}

	idea, err := services.Ideas.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if idea != nil {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	return result
}

// UpdateIdea is used to edit an existing new idea
type UpdateIdea struct {
	Model *models.UpdateIdea
	Idea  *models.Idea
}

// Initialize the model
func (input *UpdateIdea) Initialize() interface{} {
	input.Model = new(models.UpdateIdea)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdateIdea) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *UpdateIdea) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	idea, err := services.Ideas.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 || len(strings.Split(input.Model.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	}

	another, err := services.Ideas.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if another != nil && another.ID != idea.ID {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	input.Idea = idea

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
func (input *AddNewComment) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil
}

// Validate is current model is valid
func (input *AddNewComment) Validate(user *models.User, services *app.Services) *validate.Result {
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
func (input *SetResponse) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *SetResponse) Validate(user *models.User, services *app.Services) *validate.Result {
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
			if errors.Cause(err) == app.ErrNotFound {
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

// DeleteIdea represents the action of an administrator deleting an existing Idea
type DeleteIdea struct {
	Model *models.DeleteIdea
	Idea  *models.Idea
}

// Initialize the model
func (input *DeleteIdea) Initialize() interface{} {
	input.Model = new(models.DeleteIdea)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *DeleteIdea) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *DeleteIdea) Validate(user *models.User, services *app.Services) *validate.Result {
	idea, err := services.Ideas.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	isReferenced, err := services.Ideas.IsReferenced(idea)
	if err != nil {
		return validate.Error(err)
	}

	if isReferenced {
		return validate.Failed([]string{
			"This idea cannot be deleted because it's being referenced by a duplicated idea.",
		})
	}

	input.Idea = idea

	return validate.Success()
}

// EditComment represents the action to update an existing comment
type EditComment struct {
	Model *models.EditComment
}

// Initialize the model
func (input *EditComment) Initialize() interface{} {
	input.Model = new(models.EditComment)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *EditComment) IsAuthorized(user *models.User, services *app.Services) bool {
	comment, err := services.Ideas.GetCommentByID(input.Model.ID)
	if err != nil {
		return false
	}

	return user.ID == comment.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (input *EditComment) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Content == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	return result
}
