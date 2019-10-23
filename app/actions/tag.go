package actions

import (
	"context"
	"regexp"

	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/gosimple/slug"
)

var colorRegex = regexp.MustCompile(`^([A-Fa-f0-9]{6})$`)

// CreateEditTag is used to create a new tag or edit existing
type CreateEditTag struct {
	Tag   *models.Tag
	Model *models.CreateEditTag
}

// Initialize the model
func (input *CreateEditTag) Initialize() interface{} {
	input.Model = new(models.CreateEditTag)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateEditTag) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *CreateEditTag) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Slug != "" {
		getSlug := &query.GetTagBySlug{Slug: input.Model.Slug}
		err := bus.Dispatch(ctx, getSlug)
		if err != nil {
			return validate.Error(err)
		}
		input.Tag = getSlug.Result
	}

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(input.Model.Name) > 30 {
		result.AddFieldFailure("name", "Name must have less than 30 characters.")
	} else {
		getDuplicateSlug := &query.GetTagBySlug{Slug: slug.Make(input.Model.Name)}
		err := bus.Dispatch(ctx, getDuplicateSlug)
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil && (input.Tag == nil || input.Tag.ID != getDuplicateSlug.Result.ID) {
			result.AddFieldFailure("name", "This tag name is already in use.")
		}
	}

	if input.Model.Color == "" {
		result.AddFieldFailure("color", "Color is required.")
	} else if len(input.Model.Color) != 6 {
		result.AddFieldFailure("color", "Color must be exactly 6 characters.")
	} else if !colorRegex.MatchString(input.Model.Color) {
		result.AddFieldFailure("color", "Color is invalid.")
	}

	return result
}

// DeleteTag is used to delete an existing tag
type DeleteTag struct {
	Tag   *models.Tag
	Model *models.DeleteTag
}

// Initialize the model
func (input *DeleteTag) Initialize() interface{} {
	input.Model = new(models.DeleteTag)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *DeleteTag) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *DeleteTag) Validate(ctx context.Context, user *models.User) *validate.Result {
	getSlug := &query.GetTagBySlug{Slug: input.Model.Slug}
	err := bus.Dispatch(ctx, getSlug)
	if err != nil {
		return validate.Error(err)
	}

	input.Tag = getSlug.Result
	return validate.Success()
}

// AssignUnassignTag is used to assign or remove a tag to/from an post
type AssignUnassignTag struct {
	Tag   *models.Tag
	Post  *models.Post
	Model *models.AssignUnassignTag
}

// Initialize the model
func (input *AssignUnassignTag) Initialize() interface{} {
	input.Model = new(models.AssignUnassignTag)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *AssignUnassignTag) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (input *AssignUnassignTag) Validate(ctx context.Context, user *models.User) *validate.Result {
	getPost := &query.GetPostByNumber{Number: input.Model.Number}
	getSlug := &query.GetTagBySlug{Slug: input.Model.Slug}
	if err := bus.Dispatch(ctx, getPost, getSlug); err != nil {
		return validate.Error(err)
	}

	input.Post = getPost.Result
	input.Tag = getSlug.Result
	return validate.Success()
}
