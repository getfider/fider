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
	DatabaseCtxKey    = contextKey("DATABASE")
	TransactionCtxKey = contextKey("TRANSACTION")
	TenantCtxKey      = contextKey("TENANT")
	UserCtxKey        = contextKey("USER")
	ClaimsCtxKey      = contextKey("CLAIMS")   // EXPERIMENTAL-BUS: remove this key
	ServicesCtxKey    = contextKey("SERVICES") // EXPERIMENTAL-BUS: remove this key
	TasksCtxKey       = contextKey("TASKS")    // EXPERIMENTAL-BUS: remove this key
	LogPropsCtxKey    = contextKey("LOG_PROPS")
)
