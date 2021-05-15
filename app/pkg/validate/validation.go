package validate

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
)

// Validatable defines which models can be validated against context
type Validatable interface {
	Validate(ctx context.Context, user *entity.User) *Result
	IsAuthorized(ctx context.Context, user *entity.User) bool
}

// ErrorItem holds a reference to something that went wrong
type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// Result is returned after each validation
type Result struct {
	Ok         bool
	Authorized bool
	Err        error
	Errors     []ErrorItem
}

//AddFieldFailure add failure message to specific field
func (r *Result) AddFieldFailure(field string, messages ...string) {
	if r.Errors == nil {
		r.Errors = make([]ErrorItem, 0)
	}

	for _, message := range messages {
		r.Errors = append(r.Errors, ErrorItem{
			Field:   field,
			Message: message,
		})
		r.Ok = false
	}
}

// Success returns a successful validation
func Success() *Result {
	return &Result{Ok: true, Authorized: true}
}

// Failed returns a failed validation result
func Failed(messages ...string) *Result {
	r := &Result{Ok: false, Authorized: true}
	r.AddFieldFailure("", messages...)
	return r
}

// Error returns a failed validation result
func Error(err error) *Result {
	return &Result{Ok: false, Authorized: true, Err: err}
}

// Unauthorized returns an unauthorized validation result
func Unauthorized() *Result {
	return &Result{Ok: false, Authorized: false}
}
