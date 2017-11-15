package inmemory

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/gosimple/slug"
)

// TagStorage contains read and write operations for tags
type TagStorage struct {
	lastID int
	tags   []*models.Tag
}

// NewTagStorage creates a new TagStorage
func NewTagStorage() *TagStorage {
	return &TagStorage{
		tags: make([]*models.Tag, 0),
	}
}

// Add creates a new tag with given input
func (s *TagStorage) Add(name, color string, isPublic bool) (*models.Tag, error) {
	s.lastID = s.lastID + 1
	tag := &models.Tag{
		ID:       s.lastID,
		Name:     name,
		Slug:     slug.Make(name),
		Color:    color,
		IsPublic: isPublic,
	}
	s.tags = append(s.tags, tag)
	return tag, nil
}

// GetBySlug returns tag by given slug
func (s *TagStorage) GetBySlug(slug string) (*models.Tag, error) {
	for _, tag := range s.tags {
		if tag.Slug == slug {
			return tag, nil
		}
	}
	return nil, app.ErrNotFound
}

// Update a tag with given input
func (s *TagStorage) Update(tagID int, name, color string, isPublic bool) (*models.Tag, error) {
	for _, tag := range s.tags {
		if tag.ID == tagID {
			tag.Name = name
			tag.Slug = slug.Make(name)
			tag.Color = color
			tag.IsPublic = isPublic
			return tag, nil
		}
	}
	return nil, app.ErrNotFound
}

// Remove a tag by its id
func (s *TagStorage) Remove(tagID int) error {
	for i, tag := range s.tags {
		if tag.ID == tagID {
			s.tags = append(s.tags[:i], s.tags[i+1:]...)
			break
		}
	}
	return nil
}
