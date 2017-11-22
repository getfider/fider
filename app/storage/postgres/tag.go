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

// Delete a tag by its id
func (s *TagStorage) Delete(tagID int) error {
	err := s.trx.Execute(`DELETE FROM idea_tags WHERE tag_id = $1`, tagID)
	if err != nil {
		return err
	}
	return s.trx.Execute(`DELETE FROM tags WHERE id = $1 AND tenant_id = $2;`, tagID, s.tenant.ID)
}

// AssignTag adds a tag to an idea
func (s *TagStorage) AssignTag(tagID, ideaID, userID int) error {
	alreadyAssigned, err := s.trx.Exists("SELECT 1 FROM idea_tags WHERE idea_id = $1", ideaID)
	if err != nil {
		return err
	}

	if alreadyAssigned {
		return nil
	}

	if err := s.trx.Execute(
		`INSERT INTO idea_tags (tag_id, idea_id, created_on, created_by_id) VALUES ($1, $2, $3, $4)`,
		tagID, ideaID, time.Now(), userID,
	); err != nil {
		return err
	}

	return nil
}

// UnassignTag removes a tag from an idea
func (s *TagStorage) UnassignTag(tagID, ideaID int) error {
	return s.trx.Execute(
		`DELETE FROM idea_tags WHERE tag_id = $1 AND idea_id = $2`,
		tagID, ideaID,
	)
}

// GetAssigned returns all tags assigned to given idea
func (s *TagStorage) GetAssigned(ideaID int) ([]*models.Tag, error) {
	return s.getTags(`
		SELECT t.id, t.name, t.slug, t.color, t.is_public 
		FROM tags t
		INNER JOIN idea_tags it
		ON it.tag_id = t.id
		WHERE it.idea_id = $1 AND t.tenant_id = $2
	`, ideaID, s.tenant.ID)
}

// GetAll returns all tags
func (s *TagStorage) GetAll() ([]*models.Tag, error) {
	return s.getTags(`
		SELECT t.id, t.name, t.slug, t.color, t.is_public 
		FROM tags t
		WHERE t.tenant_id = $1
	`, s.tenant.ID)
}

// GetVisibleFor returns all tags visible for given user
func (s *TagStorage) GetVisibleFor(user *models.User) ([]*models.Tag, error) {
	if user != nil && user.IsCollaborator() {
		return s.GetAll()
	}
	return s.getTags(`
		SELECT t.id, t.name, t.slug, t.color, t.is_public 
		FROM tags t
		WHERE t.tenant_id = $1
		AND t.is_public = true
	`, s.tenant.ID)
}

// GetAll returns all tags
func (s *TagStorage) getTags(query string, args ...interface{}) ([]*models.Tag, error) {
	tags := []*dbTag{}
	err := s.trx.Select(&tags, query, args...)
	if err != nil {
		return nil, err
	}

	var result = make([]*models.Tag, len(tags))
	for i, tag := range tags {
		result[i] = tag.toModel()
	}
	return result, nil
}
