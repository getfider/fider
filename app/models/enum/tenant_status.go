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

var tenantStatusName = map[string]TenantStatus{
	"active":   TenantActive,
	"pending":  TenantPending,
	"locked":   TenantLocked,
	"disabled": TenantDisabled,
}

// String returns the string version of the tenant status
func (status TenantStatus) String() string {
	return tenantStatusIDs[status]
}

// MarshalText returns the Text version of the tenant status
func (status TenantStatus) MarshalText() ([]byte, error) {
	return []byte(tenantStatusIDs[status]), nil
}

// UnmarshalText parse string into a tenant status
func (status *TenantStatus) UnmarshalText(text []byte) error {
	*status = tenantStatusName[string(text)]
	return nil
}
