package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// Actionable is any action that the user can perform using the web app
type Actionable interface {
	Initialize() interface{}
	IsAuthorized(user *models.User, services *app.Services) bool
	Validate(user *models.User, services *app.Services) *validate.Result
}
