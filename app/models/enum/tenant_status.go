package enum

// TenantStatus is the status of a tenant
type TenantStatus int

var (
	//TenantActive is the default status for most tenants
	TenantActive TenantStatus = 1
	//TenantPending is used for signup via email that requires user confirmation
	TenantPending TenantStatus = 2
	//TenantLocked is used to set tenant on a read-only mode
	TenantLocked TenantStatus = 3
	//TenantDisabled is used to block all access
	TenantDisabled TenantStatus = 4
)

var tenantStatusIDs = map[TenantStatus]string{
	TenantActive:   "active",
	TenantPending:  "pending",
	TenantLocked:   "locked",
	TenantDisabled: "disabled",
}

// String returns the string version of the tenant status
func (status TenantStatus) String() string {
	return tenantStatusIDs[status]
}
