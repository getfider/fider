package validate

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// Validatable defines which models can be validated against context
type Validatable interface {
	Validate(services *app.Services) *Result
	IsAuthorized(user *models.User) bool
}

// Result is returned after each validation
type Result struct {
	Ok         bool
	Authorized bool
	Error      error
	Messages   []string
}

// Success returns a successful validation
func Success() *Result {
	return &Result{Ok: true, Authorized: true}
}

// Failed returns a failed validation result
func Failed(messages []string) *Result {
	return &Result{Ok: false, Authorized: true, Messages: messages}
}

// Error returns a failed validation result
func Error(err error) *Result {
	return &Result{Ok: false, Authorized: true, Error: err}
}

// Unauthorized returns an unauthorized validation result
func Unauthorized() *Result {
	return &Result{Ok: false, Authorized: false}
}
