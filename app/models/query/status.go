package query

import "github.com/getfider/fider/app/models/entity"

// ListActiveStatusesForTenant returns active statuses ordered by sort_order
// for the current tenant. Used by both admin UI and runtime post-rendering.
type ListActiveStatusesForTenant struct {
	Result []*entity.Status
}

// GetStatusByID fetches a single status from the current tenant by primary key.
type GetStatusByID struct {
	ID     int
	Result *entity.Status
}

// GetStatusBySlug fetches a single status from the current tenant by slug.
// Returns ErrNotFound when no row matches.
type GetStatusBySlug struct {
	Slug   string
	Result *entity.Status
}

// CountPostsByStatus returns the number of posts currently using a status.
// Used by the admin DeleteStatus flow to refuse deletion of in-use rows.
type CountPostsByStatus struct {
	StatusID int
	Result   int
}
