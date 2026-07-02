package cmd

import "github.com/getfider/fider/app/models/entity"

// CreateStatus inserts a custom status row for the current tenant.
type CreateStatus struct {
	Slug          string
	Label         string
	Kind          string
	Color         string
	Icon          string
	ShowOnHome    bool
	ShowOnRoadmap bool
	Filterable    bool
	SortOrder     int

	Result *entity.Status
}

// UpdateStatus mutates a status the current tenant owns.
// System statuses (is_system=true) can be re-labeled / re-colored / reordered
// but the Slug and Kind are immutable to preserve enum-compat behavior.
type UpdateStatus struct {
	ID            int
	Label         string
	Color         string
	Icon          string
	ShowOnHome    bool
	ShowOnRoadmap bool
	Filterable    bool
	SortOrder     int
	IsActive      bool
}

// DeleteStatus removes a status. Refused for system rows; refused if any post
// currently uses the status (caller should reassign first).
type DeleteStatus struct {
	ID int
}

// SeedTenantStatuses inserts the 6 built-in statuses for a newly created tenant.
// Idempotent: skips rows whose (tenant_id, slug) already exists.
type SeedTenantStatuses struct {
	TenantID int
}
