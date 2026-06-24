package dbEntities

import (
	"time"

	"github.com/getfider/fider/app/models/entity"
)

type Status struct {
	ID         int       `db:"id"`
	TenantID   int       `db:"tenant_id"`
	Slug       string    `db:"slug"`
	Label      string    `db:"label"`
	Kind       string    `db:"kind"`
	Color      string    `db:"color"`
	Icon       string    `db:"icon"`
	ShowOnHome    bool `db:"show_on_home"`
	ShowOnRoadmap bool `db:"show_on_roadmap"`
	Filterable    bool `db:"filterable"`
	SortOrder  int       `db:"sort_order"`
	IsSystem   bool      `db:"is_system"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (s *Status) ToModel() *entity.Status {
	if s == nil {
		return nil
	}
	return &entity.Status{
		ID:         s.ID,
		Slug:       s.Slug,
		Label:      s.Label,
		Kind:       s.Kind,
		Color:      s.Color,
		Icon:       s.Icon,
		ShowOnHome:    s.ShowOnHome,
		ShowOnRoadmap: s.ShowOnRoadmap,
		Filterable:    s.Filterable,
		SortOrder:  s.SortOrder,
		IsSystem:   s.IsSystem,
		IsActive:   s.IsActive,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
