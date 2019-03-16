package app

import "errors"

// ErrNotFound represents an object not found error
var ErrNotFound = errors.New("Object not found")

// InvitePlaceholder represents the placeholder used by members to invite other users
var InvitePlaceholder = "%invite%"

func contextKey(name string) string {
	return "FIDER_CTX_" + name
}

var (
	TransactionCtxKey = contextKey("TRANSACTION")
	TenantCtxKey      = contextKey("TENANT")
	UserCtxKey        = contextKey("USER")
	ClaimsCtxKey      = contextKey("CLAIMS")
	ServicesCtxKey    = contextKey("SERVICES")
	TasksCtxKey       = contextKey("TASKS")
)
