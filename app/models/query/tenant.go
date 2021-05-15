package query

import (
	"github.com/getfider/fider/app/models/entities"
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

	Result *entities.EmailVerification
}

type GetFirstTenant struct {
	Result *entities.Tenant
}

type GetTenantByDomain struct {
	Domain string

	Result *entities.Tenant
}
