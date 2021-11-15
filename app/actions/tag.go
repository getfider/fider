package actions

import (
	"context"
	"regexp"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/gosimple/slug"
)

var colorRegex = regexp.MustCompile(`^([A-Fa-f0-9]{6})$`)

// CreateEditTag is used to create a new tag or edit existing
type CreateEditTag struct {
	Slug     string `route:"slug"`
	Name     string `json:"name"`
	Color    string `json:"color" format:"upper"`
	IsPublic bool   `json:"isPublic"`

	Tag *entity.Tag
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateEditTag) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateEditTag) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Slug != "" {
		getSlug := &query.GetTagBySlug{Slug: action.Slug}
		err := bus.Dispatch(ctx, getSlug)
		if err != nil {
			return validate.Error(err)
		}
		action.Tag = getSlug.Result
	}

	if action.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(action.Name) > 30 {
		result.AddFieldFailure("name", "Name must have less than 30 characters.")
	} else {
		getDuplicateSlug := &query.GetTagBySlug{Slug: slug.Make(action.Name)}
		err := bus.Dispatch(ctx, getDuplicateSlug)
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil && (action.Tag == nil || action.Tag.ID != getDuplicateSlug.Result.ID) {
			result.AddFieldFailure("name", "This tag name is already in use.")
		}
	}

	if action.Color == "" {
		result.AddFieldFailure("color", "Color is required.")
	} else if len(action.Color) != 6 {
		result.AddFieldFailure("color", "Color must be exactly 6 characters.")
	} else if !colorRegex.MatchString(action.Color) {
		result.AddFieldFailure("color", "Color is invalid.")
	}

	return result
}

// DeleteTag is used to delete an existing tag
type DeleteTag struct {
	Slug string `route:"slug"`

	Tag *entity.Tag
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *DeleteTag) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *DeleteTag) Validate(ctx context.Context, user *entity.User) *validate.Result {
	getSlug := &query.GetTagBySlug{Slug: action.Slug}
	err := bus.Dispatch(ctx, getSlug)
	if err != nil {
		return validate.Error(err)
	}

	action.Tag = getSlug.Result
	return validate.Success()
}

// AssignUnassignTag is used to assign or remove a tag to/from an post
type AssignUnassignTag struct {
	Slug   string `route:"slug"`
	Number int    `route:"number"`

	Tag  *entity.Tag
	Post *entity.Post
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *AssignUnassignTag) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (action *AssignUnassignTag) Validate(ctx context.Context, user *entity.User) *validate.Result {
	getPost := &query.GetPostByNumber{Number: action.Number}
	getSlug := &query.GetTagBySlug{Slug: action.Slug}
	if err := bus.Dispatch(ctx, getPost, getSlug); err != nil {
		return validate.Error(err)
	}

	action.Post = getPost.Result
	action.Tag = getSlug.Result
	return validate.Success()
}
