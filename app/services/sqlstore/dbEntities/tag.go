package dbEntities

import (
	"github.com/getfider/fider/app/models/entity"
)

type Tag struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Slug     string `db:"slug"`
	Color    string `db:"color"`
	IsPublic bool   `db:"is_public"`
}

func (t *Tag) ToModel() *entity.Tag {
	return &entity.Tag{
		ID:       t.ID,
		Name:     t.Name,
		Slug:     t.Slug,
		Color:    t.Color,
		IsPublic: t.IsPublic,
	}
}
