package query

import (
	"github.com/getfider/fider/app/models"
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

	Result *models.EmailVerification
}

type GetFirstTenant struct {
	Result *models.Tenant
}

type GetTenantByDomain struct {
	Domain string

	Result *models.Tenant
}
