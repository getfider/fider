package app

import "errors"

// ErrNotFound represents an object not found error
var ErrNotFound = errors.New("Object not found")

// InvitePlaceholder represents the placeholder used by members to invite other users
var InvitePlaceholder = "%invite%"

type key string

func createKey(name string) key {
	return key("FIDER_CTX_" + name)
}

var (
	RequestCtxKey     = createKey("REQUEST")
	DatabaseCtxKey    = createKey("DATABASE")
	TransactionCtxKey = createKey("TRANSACTION")
	TenantCtxKey      = createKey("TENANT")
	UserCtxKey        = createKey("USER")
	ClaimsCtxKey      = createKey("CLAIMS")   // EXPERIMENTAL-BUS: remove this key
	ServicesCtxKey    = createKey("SERVICES") // EXPERIMENTAL-BUS: remove this key
	LogPropsCtxKey    = createKey("LOG_PROPS")
)
