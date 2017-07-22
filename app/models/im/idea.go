package im

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// Idea represents an unsaved idea
type Idea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (i *Idea) IsAuthorized(user *models.User) bool {
	return true
}

// Validate is current model is valid
func (i *Idea) Validate(services *app.Services) *validate.Result {
	result := validate.Success()
	if strings.Trim(i.Title, " ") == "" {
		result.AddFieldFailure("title", "Title is required.")
	}

	if len(i.Title) < 10 || len(strings.Split(i.Title, " ")) < 3 {
		result.AddFieldFailure("title", "Title needs to be more descriptive.")
	}

	return result
}
