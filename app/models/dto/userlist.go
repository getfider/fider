package dto

import "github.com/getfider/fider/app/models/enum"

type UserListUpdateCompany struct {
	TenantID      int
	Name          string
	BillingStatus enum.BillingStatus
}
