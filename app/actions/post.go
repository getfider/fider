package actions

import (
	"context"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
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
func (input *CreateNewPost) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil
}

// Validate if current model is valid
func (input *CreateNewPost) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	} else if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	} else {
		err := bus.Dispatch(ctx, &query.GetPostBySlug{Slug: slug.Make(input.Model.Title)})
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil {
			result.AddFieldFailure("title", "This has already been posted before.")
		}
	}

	messages, err := validate.MultiImageUpload(nil, input.Model.Attachments, validate.MultiImageUploadOpts{
		MaxUploads:   3,
		MaxKilobytes: 5120,
		ExactRatio:   false,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("attachments", messages...)

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
func (input *UpdatePost) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (input *UpdatePost) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Title == "" {
		result.AddFieldFailure("title", "Title is required.")
	} else if len(input.Model.Title) < 10 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	if len(input.Model.Title) > 100 {
		result.AddFieldFailure("title", "Title must have less than 100 characters.")
	}

	postByNumber := &query.GetPostByNumber{Number: input.Model.Number}
	if err := bus.Dispatch(ctx, postByNumber); err != nil {
		return validate.Error(err)
	}

	input.Post = postByNumber.Result

	postBySlug := &query.GetPostBySlug{Slug: slug.Make(input.Model.Title)}
	err := bus.Dispatch(ctx, postBySlug)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if err == nil && postBySlug.Result.ID != input.Post.ID {
		result.AddFieldFailure("title", "This has already been posted before.")
	}

	if len(input.Model.Attachments) > 0 {
		getAttachments := &query.GetAttachments{Post: input.Post}
		err = bus.Dispatch(ctx, getAttachments)
		if err != nil {
			return validate.Error(err)
		}

		messages, err := validate.MultiImageUpload(getAttachments.Result, input.Model.Attachments, validate.MultiImageUploadOpts{
			MaxUploads:   3,
			MaxKilobytes: 5120,
			ExactRatio:   false,
		})
		if err != nil {
			return validate.Error(err)
		}
		result.AddFieldFailure("attachments", messages...)
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
func (input *AddNewComment) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil
}

// Validate if current model is valid
func (input *AddNewComment) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Content == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	messages, err := validate.MultiImageUpload(nil, input.Model.Attachments, validate.MultiImageUploadOpts{
		MaxUploads:   2,
		MaxKilobytes: 5120,
		ExactRatio:   false,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("attachments", messages...)

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
func (input *SetResponse) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (input *SetResponse) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Status < enum.PostOpen || input.Model.Status > enum.PostDuplicate {
		result.AddFieldFailure("status", "Status is invalid.")
	}

	if input.Model.Status == enum.PostDuplicate {
		if input.Model.OriginalNumber == input.Model.Number {
			result.AddFieldFailure("originalNumber", "Cannot be a duplicate of itself")
		}

		getOriginaPost := &query.GetPostByNumber{Number: input.Model.OriginalNumber}
		err := bus.Dispatch(ctx, getOriginaPost)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				result.AddFieldFailure("originalNumber", "Original post not found")
			} else {
				return validate.Error(err)
			}
		}

		if getOriginaPost.Result != nil {
			input.Original = getOriginaPost.Result
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
func (input *DeletePost) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *DeletePost) Validate(ctx context.Context, user *models.User) *validate.Result {
	getPost := &query.GetPostByNumber{Number: input.Model.Number}
	if err := bus.Dispatch(ctx, getPost); err != nil {
		return validate.Error(err)
	}

	input.Post = getPost.Result

	isReferencedQuery := &query.PostIsReferenced{PostID: input.Post.ID}
	if err := bus.Dispatch(ctx, isReferencedQuery); err != nil {
		return validate.Error(err)
	}

	if isReferencedQuery.Result {
		return validate.Failed("This post cannot be deleted because it's being referenced by a duplicated post.")
	}

	return validate.Success()
}

// EditComment represents the action to update an existing comment
type EditComment struct {
	Model   *models.EditComment
	Post    *models.Post
	Comment *models.Comment
}

// Initialize the model
func (input *EditComment) Initialize() interface{} {
	input.Model = new(models.EditComment)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *EditComment) IsAuthorized(ctx context.Context, user *models.User) bool {
	postByNumber := &query.GetPostByNumber{Number: input.Model.PostNumber}
	commentByID := &query.GetCommentByID{CommentID: input.Model.ID}
	if err := bus.Dispatch(ctx, postByNumber, commentByID); err != nil {
		return false
	}

	input.Post = postByNumber.Result
	input.Comment = commentByID.Result
	return user.ID == input.Comment.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (input *EditComment) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Content == "" {
		result.AddFieldFailure("content", "Comment is required.")
	}

	if len(input.Model.Attachments) > 0 {
		getAttachments := &query.GetAttachments{Post: input.Post, Comment: input.Comment}
		err := bus.Dispatch(ctx, getAttachments)
		if err != nil {
			return validate.Error(err)
		}

		messages, err := validate.MultiImageUpload(getAttachments.Result, input.Model.Attachments, validate.MultiImageUploadOpts{
			MaxUploads:   2,
			MaxKilobytes: 5120,
			ExactRatio:   false,
		})
		if err != nil {
			return validate.Error(err)
		}
		result.AddFieldFailure("attachments", messages...)
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
func (input *DeleteComment) IsAuthorized(ctx context.Context, user *models.User) bool {
	commentByID := &query.GetCommentByID{CommentID: input.Model.CommentID}
	if err := bus.Dispatch(ctx, commentByID); err != nil {
		return false
	}

	return user.ID == commentByID.Result.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (input *DeleteComment) Validate(ctx context.Context, user *models.User) *validate.Result {
	return validate.Success()
}
