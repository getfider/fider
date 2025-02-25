package dto

import "github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"

type UserListUpdateCompany struct {
	TenantID      int
	Name          string
	BillingStatus enum.BillingStatus
}
