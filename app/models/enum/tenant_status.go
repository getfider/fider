package enum

var (
	//TenantActive is the default status for most tenants
	TenantActive = 1
	//TenantPending is used for signup via email that requires user confirmation
	TenantPending = 2
	//TenantLocked is used when tenants are locked for various reasons
	TenantLocked = 3
)
