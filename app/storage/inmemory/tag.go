package inmemory

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/gosimple/slug"
)

// TagStorage contains read and write operations for tags
type TagStorage struct {
	lastID   int
	tags     []*models.Tag
	user     *models.User
	tenant   *models.Tenant
	assigned map[int][]*models.Tag
}

// NewTagStorage creates a new TagStorage
func NewTagStorage() *TagStorage {
	return &TagStorage{
		tags:     make([]*models.Tag, 0),
		assigned: make(map[int][]*models.Tag),
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
func (s *TagStorage) Update(tag *models.Tag, name, color string, isPublic bool) (*models.Tag, error) {
	for _, storedTag := range s.tags {
		if storedTag.ID == tag.ID {
			storedTag.Name = name
			storedTag.Slug = slug.Make(name)
			storedTag.Color = color
			storedTag.IsPublic = isPublic
			return tag, nil
		}
	}
	return nil, app.ErrNotFound
}

// Delete a tag by its id
func (s *TagStorage) Delete(tag *models.Tag) error {
	for i, storedTag := range s.tags {
		if storedTag.ID == tag.ID {
			s.tags = append(s.tags[:i], s.tags[i+1:]...)
			break
		}
	}
	return nil
}

// AssignTag adds a tag to an post
func (s *TagStorage) AssignTag(tag *models.Tag, post *models.Post) error {
	assigned := s.assigned[post.ID]
	if assigned != nil {
		for _, assignedTag := range assigned {
			if assignedTag.ID == tag.ID {
				return nil
			}
		}
	}

	s.assigned[post.ID] = append(s.assigned[post.ID], tag)
	return nil
}

// UnassignTag removes a tag from an post
func (s *TagStorage) UnassignTag(tag *models.Tag, post *models.Post) error {
	assigned := s.assigned[post.ID]
	if assigned != nil {
		for i, assignedTag := range assigned {
			if assignedTag.ID == tag.ID {
				s.assigned[post.ID] = append(s.assigned[post.ID][:i], s.assigned[post.ID][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

// GetAssigned returns all tags assigned to given post
func (s *TagStorage) GetAssigned(post *models.Post) ([]*models.Tag, error) {
	assigned := s.assigned[post.ID]
	if assigned != nil {
		return assigned, nil
	}
	return make([]*models.Tag, 0), nil
}

// GetAll returns all tags
func (s *TagStorage) GetAll() ([]*models.Tag, error) {
	if s.user != nil && s.user.IsCollaborator() {
		return s.tags, nil
	}
	tags := make([]*models.Tag, 0)
	for _, tag := range s.tags {
		if tag.IsPublic {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
