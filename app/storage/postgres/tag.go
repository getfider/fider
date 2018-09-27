package postgres

import (
	"fmt"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
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
	user   *models.User
}

// NewTagStorage creates a new TagStorage
func NewTagStorage(trx *dbx.Trx) *TagStorage {
	return &TagStorage{
		trx: trx,
	}
}

// SetCurrentTenant to current context
func (s *TagStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *TagStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// Add creates a new tag with given input
func (s *TagStorage) Add(name, color string, isPublic bool) (*models.Tag, error) {
	tagSlug := slug.Make(name)

	var id int
	err := s.trx.Get(&id, `
		INSERT INTO tags (name, slug, color, is_public, created_at, tenant_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, name, tagSlug, color, isPublic, time.Now(), s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new tag")
	}

	return s.GetBySlug(tagSlug)
}

// GetBySlug returns tag by given slug
func (s *TagStorage) GetBySlug(slug string) (*models.Tag, error) {
	tag := dbTag{}

	err := s.trx.Get(&tag, "SELECT id, name, slug, color, is_public FROM tags WHERE tenant_id = $1 AND slug = $2", s.tenant.ID, slug)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tag with slug '%s'", slug)
	}

	return tag.toModel(), nil
}

// Update a tag with given input
func (s *TagStorage) Update(tag *models.Tag, name, color string, isPublic bool) (*models.Tag, error) {
	tagSlug := slug.Make(name)

	_, err := s.trx.Execute(`UPDATE tags SET name = $1, slug = $2, color = $3, is_public = $4
												WHERE id = $5 AND tenant_id = $6`, name, tagSlug, color, isPublic, tag.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update tag")
	}

	return s.GetBySlug(tagSlug)
}

// Delete a tag by its id
func (s *TagStorage) Delete(tag *models.Tag) error {
	_, err := s.trx.Execute(`DELETE FROM post_tags WHERE tag_id = $1 AND tenant_id = $2`, tag.ID, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to remove tag with id '%d' from all posts", tag.ID)
	}

	_, err = s.trx.Execute(`DELETE FROM tags WHERE id = $1 AND tenant_id = $2`, tag.ID, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete tag with id '%d'", tag.ID)
	}
	return nil
}

// AssignTag adds a tag to an post
func (s *TagStorage) AssignTag(tag *models.Tag, post *models.Post) error {
	alreadyAssigned, err := s.trx.Exists("SELECT 1 FROM post_tags WHERE post_id = $1 AND tag_id = $2 AND tenant_id = $3", post.ID, tag.ID, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to check if tag is already assigned")
	}

	if alreadyAssigned {
		return nil
	}

	_, err = s.trx.Execute(
		`INSERT INTO post_tags (tag_id, post_id, created_at, created_by_id, tenant_id) VALUES ($1, $2, $3, $4, $5)`,
		tag.ID, post.ID, time.Now(), s.user.ID, s.tenant.ID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to assign tag to post")
	}
	return nil
}

// UnassignTag removes a tag from an post
func (s *TagStorage) UnassignTag(tag *models.Tag, post *models.Post) error {
	_, err := s.trx.Execute(
		`DELETE FROM post_tags WHERE tag_id = $1 AND post_id = $2 AND tenant_id = $3`,
		tag.ID, post.ID, s.tenant.ID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to unassign tag from post")
	}
	return nil
}

// GetAssigned returns all tags assigned to given post
func (s *TagStorage) GetAssigned(post *models.Post) ([]*models.Tag, error) {
	tags, err := s.getTags(`
		SELECT t.id, t.name, t.slug, t.color, t.is_public 
		FROM tags t
		INNER JOIN post_tags pt
		ON pt.tag_id = t.id
		AND pt.tenant_id = t.tenant_id
		WHERE pt.post_id = $1 AND t.tenant_id = $2
		ORDER BY t.name
	`, post.ID, s.tenant.ID)

	if err != nil {
		return nil, errors.Wrap(err, "failed get assigned tags")
	}
	return tags, nil
}

// GetAll returns all tags
func (s *TagStorage) GetAll() ([]*models.Tag, error) {
	condition := `AND t.is_public = true`
	if s.user != nil && s.user.IsCollaborator() {
		condition = ``
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.name, t.slug, t.color, t.is_public 
		FROM tags t
		WHERE t.tenant_id = $1 %s
		ORDER BY t.name
	`, condition)
	tags, err := s.getTags(query, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get all tags")
	}
	return tags, nil
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
