package query

import (
	"time"

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

type GetFirstTenant struct {

	// Output
	Result *entity.Tenant
}

type GetTenantByDomain struct {
	Domain string

	// Output
	Result *entity.Tenant
}

type GetTrialingTenantContacts struct {
	TrialExpiresOn time.Time

	// Output
	Contacts []*entity.User
}
