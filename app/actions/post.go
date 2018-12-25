package actions

import (
	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateNewPost is used to create a new post
type CreateNewPost struct {
	Model *models.NewPost
}

// Initialize the model
func (input *CreateNewPost) Initialize() interface{} {
	input.Model = new(models.NewPost)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateNewPost) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil
}

// Validate is current model is valid
func (input *CreateNewPost) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	}

	post, err := services.Posts.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if post != nil {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	return result
}

// UpdatePost is used to edit an existing new post
type UpdatePost struct {
	Model *models.UpdatePost
	Post  *models.Post
}

// Initialize the model
func (input *UpdatePost) Initialize() interface{} {
	input.Model = new(models.UpdatePost)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdatePost) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *UpdatePost) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	post, err := services.Posts.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	}

	another, err := services.Posts.GetBySlug(slug.Make(input.Model.Title))
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if another != nil && another.ID != post.ID {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	input.Post = post

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

// SetResponse represents the action to update an post response
type SetResponse struct {
	Model    *models.SetResponse
	Original *models.Post
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

	if input.Model.Status < models.PostOpen || input.Model.Status > models.PostDuplicate {
		result.AddFieldFailure("status", "Status is invalid.")
	}

	if input.Model.Status == models.PostDuplicate {
		if input.Model.OriginalNumber == input.Model.Number {
			result.AddFieldFailure("originalNumber", "Cannot be a duplicate of itself")
		}

		original, err := services.Posts.GetByNumber(input.Model.OriginalNumber)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				result.AddFieldFailure("originalNumber", "Original post not found")
			} else {
				return validate.Error(err)
			}
		}
		if original != nil {
			input.Original = original
		}
	}

	return result
}

// DeletePost represents the action of an administrator deleting an existing Post
type DeletePost struct {
	Model *models.DeletePost
	Post  *models.Post
}

// Initialize the model
func (input *DeletePost) Initialize() interface{} {
	input.Model = new(models.DeletePost)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *DeletePost) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *DeletePost) Validate(user *models.User, services *app.Services) *validate.Result {
	post, err := services.Posts.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	isReferenced, err := services.Posts.IsReferenced(post)
	if err != nil {
		return validate.Error(err)
	}

	if isReferenced {
		return validate.Failed("This post cannot be deleted because it's being referenced by a duplicated post.")
	}

	input.Post = post

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
	comment, err := services.Posts.GetCommentByID(input.Model.ID)
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

// DeleteComment represents the action of deleting an existing comment
type DeleteComment struct {
	Model *models.DeleteComment
}

// Initialize the model
func (input *DeleteComment) Initialize() interface{} {
	input.Model = new(models.DeleteComment)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *DeleteComment) IsAuthorized(user *models.User, services *app.Services) bool {
	comment, err := services.Posts.GetCommentByID(input.Model.CommentID)
	if err != nil {
		return false
	}

	return user.ID == comment.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (input *DeleteComment) Validate(user *models.User, services *app.Services) *validate.Result {
	return validate.Success()
}
