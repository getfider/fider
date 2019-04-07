package query

type IsCNAMEAvailable struct {
	CNAME string

	Result bool
}

type IsSubdomainAvailable struct {
	Subdomain string

	Result bool
}
