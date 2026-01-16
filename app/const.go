package app

import "errors"

// ErrNotFound represents an object not found error
var ErrNotFound = errors.New("object not found")

// ErrCommercialLicenseRequired is used when a commercial feature is accessed without a license
var ErrCommercialLicenseRequired = errors.New("content moderation requires a commercial license")

// InvitePlaceholder represents the placeholder used by members to invite other users
var InvitePlaceholder = "%invite%"

// ErrUserIDRequired is used when OAuth integration returns an empty user ID
var ErrUserIDRequired = errors.New("UserID is required during OAuth integration")

type key string

func createKey(name string) key {
	return key("FIDER_CTX_" + name)
}

const (
	//FacebookProvider is const for 'facebook'
	FacebookProvider = "facebook"
	//GoogleProvider is const for 'google'
	GoogleProvider = "google"
	//GitHubProvider is const for 'github'
	GitHubProvider = "github"
)

var (
	RequestCtxKey     = createKey("REQUEST")
	TransactionCtxKey = createKey("TRANSACTION")
	TenantCtxKey      = createKey("TENANT")
	LocaleCtxKey      = createKey("LOCALE")
	UserCtxKey        = createKey("USER")
	LogPropsCtxKey    = createKey("LOG_PROPS")
)
