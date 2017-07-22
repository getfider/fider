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
	Failures   map[string][]string
}

//AddFailure as a general message
func (r *Result) AddFailure(message string) {
	if r.Messages == nil {
		r.Messages = []string{}
	}
	r.Messages = append(r.Messages, message)
	r.Ok = false
}

//AddFieldFailure add failure message to specific field
func (r *Result) AddFieldFailure(field string, messages ...string) {
	if r.Failures == nil {
		r.Failures = make(map[string][]string)
	}

	if r.Failures[field] == nil {
		r.Failures[field] = []string{}
	}

	for _, message := range messages {
		r.Failures[field] = append(r.Failures[field], message)
	}
	r.Ok = false
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
