package query

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type IsCNAMEAvailable struct {
	CNAME string

	Result bool
}

type IsSubdomainAvailable struct {
	Subdomain string

	Result bool
}

type GetVerificationByKey struct {
	Kind enum.EmailVerificationKind
	Key  string

	Result *entity.EmailVerification
}

type GetFirstTenant struct {
	Result *entity.Tenant
}

type GetTenantByDomain struct {
	Domain string

	Result *entity.Tenant
}
