package cmd

type UserListCreateCompany struct {
	Name          string
	TenantId      int
	SignedUpAt    string
	BillingStatus string
	Subdomain     string
	UserId        int
	UserEmail     string
	UserName      string
}

type UserListUpdateCompany struct {
	Id            int
	Name          string
	BillingStatus string
	Subdomain     string
}

type UpdateUserListUser struct {
	Id    int
	Email string
	Name  string
}
