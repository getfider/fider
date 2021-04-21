package enum

var (
	//TenantActive is the default status for most tenants
	TenantActive = 1
	//TenantPending is used for signup via email that requires user confirmation
	TenantPending = 2
	//TenantLocked is used to set tenant on a read-only mode
	TenantLocked = 3
	//TenantDisabled is used to block all access
	TenantDisabled = 4
)
