package actions

import (
	"context"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/validate"
)

var statusSlugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var allowedStatusKinds = map[string]bool{
	"open":             true,
	"active":           true,
	"closed-completed": true,
	"closed-declined":  true,
	"duplicate":        true,
}

var allowedStatusColors = map[string]bool{
	"blue": true, "green": true, "yellow": true, "red": true, "gray": true,
}

// CreateStatus is the admin action to add a new custom status.
type CreateStatus struct {
	Slug          string `json:"slug" format:"lower"`
	Label         string `json:"label"`
	Kind          string `json:"kind"`
	Color         string `json:"color"`
	Icon          string `json:"icon"`
	ShowOnHome    bool   `json:"showOnHome"`
	ShowOnRoadmap bool   `json:"showOnRoadmap"`
	Filterable    bool   `json:"filterable"`
	SortOrder     int    `json:"sortOrder"`
}

func (action *CreateStatus) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

func (action *CreateStatus) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	action.Slug = strings.TrimSpace(strings.ToLower(action.Slug))
	action.Label = strings.TrimSpace(action.Label)
	action.Kind = strings.TrimSpace(action.Kind)
	if action.Color == "" {
		action.Color = "blue"
	}
	if action.Icon == "" {
		action.Icon = "lightbulb"
	}

	if action.Slug == "" {
		result.AddFieldFailure("slug", "Slug is required.")
	} else if len(action.Slug) > 50 {
		result.AddFieldFailure("slug", "Slug must be 50 characters or fewer.")
	} else if !statusSlugRegex.MatchString(action.Slug) {
		result.AddFieldFailure("slug", "Slug must use only lowercase letters, numbers, and hyphens.")
	}

	if action.Label == "" {
		result.AddFieldFailure("label", "Label is required.")
	} else if len(action.Label) > 50 {
		result.AddFieldFailure("label", "Label must be 50 characters or fewer.")
	}

	if !allowedStatusKinds[action.Kind] {
		result.AddFieldFailure("kind", "Kind must be one of: open, active, closed-completed, closed-declined, duplicate.")
	}

	if !allowedStatusColors[action.Color] {
		result.AddFieldFailure("color", "Color must be one of: blue, green, yellow, red, gray.")
	}

	return result
}

// UpdateStatus is the admin action to edit an existing status. ID comes from
// the URL; Slug and Kind are not editable here (immutable to preserve any
// outstanding posts and webhook integrations).
type UpdateStatus struct {
	ID            int    `route:"id"`
	Label         string `json:"label"`
	Color         string `json:"color"`
	Icon          string `json:"icon"`
	ShowOnHome    bool   `json:"showOnHome"`
	ShowOnRoadmap bool   `json:"showOnRoadmap"`
	Filterable    bool   `json:"filterable"`
	SortOrder     int    `json:"sortOrder"`
	IsActive      bool   `json:"isActive"`
}

func (action *UpdateStatus) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

func (action *UpdateStatus) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	action.Label = strings.TrimSpace(action.Label)
	if action.Color == "" {
		action.Color = "blue"
	}
	if action.Icon == "" {
		action.Icon = "lightbulb"
	}

	if action.Label == "" {
		result.AddFieldFailure("label", "Label is required.")
	} else if len(action.Label) > 50 {
		result.AddFieldFailure("label", "Label must be 50 characters or fewer.")
	}

	if !allowedStatusColors[action.Color] {
		result.AddFieldFailure("color", "Color must be one of: blue, green, yellow, red, gray.")
	}

	return result
}

// DeleteStatus removes a non-system status. ID from the URL only.
type DeleteStatus struct {
	ID int `route:"id"`
}

func (action *DeleteStatus) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.Role == enum.RoleAdministrator
}

func (action *DeleteStatus) Validate(ctx context.Context, user *entity.User) *validate.Result {
	return validate.Success()
}
