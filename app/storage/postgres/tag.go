package postgres

import (
	"database/sql"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/gosimple/slug"
)

type dbTag struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Slug     string `db:"slug"`
	Color    string `db:"color"`
	IsPublic bool   `db:"is_public"`
}

func (t *dbTag) toModel() *models.Tag {
	return &models.Tag{
		ID:       t.ID,
		Name:     t.Name,
		Slug:     t.Slug,
		Color:    t.Color,
		IsPublic: t.IsPublic,
	}
}

// TagStorage contains read and write operations for tags
type TagStorage struct {
	trx    *dbx.Trx
	tenant *models.Tenant
}

// NewTagStorage creates a new TagStorage
func NewTagStorage(tenant *models.Tenant, trx *dbx.Trx) *TagStorage {
	return &TagStorage{
		trx:    trx,
		tenant: tenant,
	}
}

// Add creates a new tag with given input
func (s *TagStorage) Add(name, color string, isPublic bool) (*models.Tag, error) {
	tagSlug := slug.Make(name)

	var id int
	row := s.trx.QueryRow(`
		INSERT INTO tags (name, slug, color, is_public, created_on, tenant_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, name, tagSlug, color, isPublic, time.Now(), s.tenant.ID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return s.GetBySlug(tagSlug)
}

// GetBySlug returns tag by given slug
func (s *TagStorage) GetBySlug(slug string) (*models.Tag, error) {
	tag := dbTag{}

	err := s.trx.Get(&tag, "SELECT id, name, slug, color, is_public FROM tags WHERE tenant_id = $1 AND slug = $2", s.tenant.ID, slug)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return tag.toModel(), nil
}

// Update a tag with given input
func (s *TagStorage) Update(tagID int, name, color string, isPublic bool) (*models.Tag, error) {
	tagSlug := slug.Make(name)

	err := s.trx.Execute(`UPDATE tags SET name = $1, slug = $2, color = $3, is_public = $4
												WHERE id = $5 AND tenant_id = $6`, name, tagSlug, color, isPublic, tagID, s.tenant.ID)
	if err != nil {
		return nil, err
	}

	return s.GetBySlug(tagSlug)
}

// Remove a tag by its id
func (s *TagStorage) Remove(tagID int) error {
	return s.trx.Execute(`DELETE FROM tags WHERE id = $1 AND tenant_id = $2`, tagID, s.tenant.ID)
}
