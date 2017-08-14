package actions

import (
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
	return validate.Success()
}

type Actionable interface {
	NewModel() interface{}
	Validate(services *app.Services) *validate.Result
	IsAuthorized(user *models.User) bool
}
