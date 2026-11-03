package query

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type IsCNAMEAvailable struct {
	CNAME string

	// Output
	Result bool
}

type IsSubdomainAvailable struct {
	Subdomain string

	// Output
	Result bool
}

type GetVerificationByKey struct {
	Kind enum.EmailVerificationKind
	Key  string

	// Output
	Result *entity.EmailVerification
}

type GetVerificationByEmailAndCode struct {
	Email string
	Code  string
	Kind  enum.EmailVerificationKind

	// Output
	Result *entity.EmailVerification
}

type GetActiveVerificationByEmail struct {
	Email string
	Kind  enum.EmailVerificationKind

	// Output
	Result *entity.EmailVerification
}

type GetFirstTenant struct {

	// Output
	Result *entity.Tenant
}

type GetTenantByDomain struct {
	Domain string

	// Output
	Result *entity.Tenant
}

type GetPendingSignUpVerification struct {
	// Output
	Result *entity.EmailVerification
}

// GetTenantsPendingDeletion returns tenants whose scheduled_deletion_at has passed,
// ordered oldest-first. Consumed by the deletion cron job.
type GetTenantsPendingDeletion struct {
	// Output
	Result []*entity.Tenant
}

// GetTenantByCancelKey looks up the tenant holding the given deletion cancel key.
type GetTenantByCancelKey struct {
	Key string

	// Output
	Result *entity.Tenant
}

// GetTenantOwner returns the account owner: the active Administrator with the lowest
// users.id (the user that created the tenant at signup).
type GetTenantOwner struct {
	TenantID int

	// Output
	Result *entity.User
}
