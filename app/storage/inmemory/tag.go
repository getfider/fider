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

// Delete a tag by its id
func (s *TagStorage) Delete(tagID int) error {
	for i, tag := range s.tags {
		if tag.ID == tagID {
			s.tags = append(s.tags[:i], s.tags[i+1:]...)
			break
		}
	}
	return nil
}

// AssignTag adds a tag to an idea
func (s *TagStorage) AssignTag(tagID, ideaID, userID int) error {
	var tagToAssign *models.Tag
	for _, tag := range s.tags {
		if tag.ID == tagID {
			tagToAssign = tag
			break
		}
	}

	if tagToAssign == nil {
		return app.ErrNotFound
	}

	assigned := s.assigned[ideaID]
	if assigned != nil {
		for _, tag := range assigned {
			if tag.ID == tagID {
				return nil
			}
		}
	}

	s.assigned[ideaID] = append(s.assigned[ideaID], tagToAssign)
	return nil
}

// UnassignTag removes a tag from an idea
func (s *TagStorage) UnassignTag(tagID, ideaID int) error {
	assigned := s.assigned[ideaID]
	if assigned != nil {
		for i, tag := range assigned {
			if tag.ID == tagID {
				s.assigned[ideaID] = append(s.assigned[ideaID][:i], s.assigned[ideaID][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

// GetAssigned returns all tags assigned to given idea
func (s *TagStorage) GetAssigned(ideaID int) ([]*models.Tag, error) {
	assigned := s.assigned[ideaID]
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
