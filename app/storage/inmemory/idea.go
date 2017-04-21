package inmemory

import (
	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
)

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	lastID int
	ideas  []*models.Idea
}

// GetByID returns idea by given id
func (s *IdeaStorage) GetByID(tenantID, ideaID int) (*models.Idea, error) {
	for _, idea := range s.ideas {
		if idea.ID == ideaID {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetByNumber returns idea by tenant and number
func (s *IdeaStorage) GetByNumber(tenantID, number int) (*models.Idea, error) {
	for _, idea := range s.ideas {
		if idea.Number == number {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll(tenantID int) ([]*models.Idea, error) {
	return s.ideas, nil
}

// GetCommentsByIdeaID returns all coments from given idea
func (s *IdeaStorage) GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error) {
	return make([]*models.Comment, 0), nil
}

// Save a new idea in the database
func (s *IdeaStorage) Save(tenantID, userID int, title, description string) (*models.Idea, error) {
	s.lastID = s.lastID + 1
	idea := &models.Idea{
		ID:          s.lastID,
		Number:      s.lastID,
		Title:       title,
		Description: description,
	}
	s.ideas = append(s.ideas, idea)
	return idea, nil
}

// AddComment places a new comment on an idea
func (s *IdeaStorage) AddComment(userID, ideaID int, content string) (int, error) {
	return 0, nil
}

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(userID, ideaID int) error {
	for _, idea := range s.ideas {
		if idea.ID == ideaID {
			idea.TotalSupporters = idea.TotalSupporters + 1
		}
	}
	return nil
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(userID, ideaID int) error {
	for _, idea := range s.ideas {
		if idea.ID == ideaID {
			idea.TotalSupporters = idea.TotalSupporters - 1
		}
	}
	return nil
}

// GetSupportedIdeas returns all ideas supported by given user
func (s *IdeaStorage) GetSupportedIdeas(userID int) ([]int, error) {
	return make([]int, 0), nil
}
