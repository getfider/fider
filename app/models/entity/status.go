package entity

import "time"

// Status is a tenant-configurable label that classifies posts. Each tenant owns
// a list of statuses, seeded with the built-in set (Open, Planned, Started,
// Completed, Declined, Duplicate) and extendable by admins via the UI.
//
// Kind anchors a custom label to a fixed semantic so the rest of Fider still
// knows what "closed" means without parsing names. Allowed kinds:
//
//	open              — default initial state
//	active            — accepted onto the roadmap; in progress
//	closed-completed  — done; positive resolution
//	closed-declined   — closed without action
//	duplicate         — closed as a dupe of another post
type Status struct {
	ID         int       `json:"id"`
	Slug       string    `json:"slug"`
	Label      string    `json:"label"`
	Kind       string    `json:"kind"`
	Color      string    `json:"color"`
	Icon       string    `json:"icon"`
	ShowOnHome    bool `json:"showOnHome"`
	ShowOnRoadmap bool `json:"showOnRoadmap"`
	Filterable    bool `json:"filterable"`
	SortOrder  int       `json:"sortOrder"`
	IsSystem   bool      `json:"isSystem"`
	IsActive   bool      `json:"isActive"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// IsClosed reports whether posts in this status are considered closed.
// Used to drive UI filtering, vote-eligibility, and notification gating.
func (s Status) IsClosed() bool {
	return s.Kind == "closed-completed" || s.Kind == "closed-declined" || s.Kind == "duplicate"
}
