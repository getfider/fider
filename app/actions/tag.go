package actions

import (
	"regexp"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
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
func (input *CreateEditTag) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate is current model is valid
func (input *CreateEditTag) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()
	if input.Model.Slug != "" {
		tag, err := services.Tags.GetBySlug(input.Model.Slug)
		if err != nil {
			return validate.Error(err)
		}
		input.Tag = tag
	}

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(input.Model.Name) > 30 {
		result.AddFieldFailure("name", "Name must have less than 30 characters.")
	} else {
		duplicateTag, err := services.Tags.GetBySlug(slug.Make(input.Model.Name))
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil && (input.Tag == nil || input.Tag.ID != duplicateTag.ID) {
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
func (input *DeleteTag) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate is current model is valid
func (input *DeleteTag) Validate(user *models.User, services *app.Services) *validate.Result {
	tag, err := services.Tags.GetBySlug(input.Model.Slug)
	if err != nil {
		return validate.Error(err)
	}

	input.Tag = tag
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
func (input *AssignUnassignTag) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsCollaborator()
}

// Validate is current model is valid
func (input *AssignUnassignTag) Validate(user *models.User, services *app.Services) *validate.Result {
	post, err := services.Posts.GetByNumber(input.Model.Number)
	if err != nil {
		return validate.Error(err)
	}

	tag, err := services.Tags.GetBySlug(input.Model.Slug)
	if err != nil {
		return validate.Error(err)
	}

	input.Tag = tag
	input.Post = post
	return validate.Success()
}
