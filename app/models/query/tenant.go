package query

import "github.com/getfider/fider/app/models"

type IsCNAMEAvailable struct {
	CNAME string

	Result bool
}

type IsSubdomainAvailable struct {
	Subdomain string

	Result bool
}

type GetVerificationByKey struct {
	Kind models.EmailVerificationKind
	Key  string

	Result *models.EmailVerification
}
