package actions

import (
	"regexp"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/gosimple/slug"
)

var colorRegex = regexp.MustCompile(`^([A-Fa-f0-9]{6})$`)

// CreateNewTag is used to create a new tag
type CreateNewTag struct {
	Model *models.CreateNewTag
}

// Initialize the model
func (input *CreateNewTag) Initialize() interface{} {
	input.Model = new(models.CreateNewTag)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateNewTag) IsAuthorized(user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate is current model is valid
func (input *CreateNewTag) Validate(services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(input.Model.Name) > 30 {
		result.AddFieldFailure("name", "Name must be less than 30 characters.")
	} else {
		_, err := services.Tags.GetBySlug(slug.Make(input.Model.Name))
		if err != nil && err != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil {
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
