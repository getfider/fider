package inmemory

import (
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	lastID int
	ideas  []*models.Idea
}

// GetByID returns idea by given id
func (s *IdeaStorage) GetByID(ideaID int) (*models.Idea, error) {
	for _, idea := range s.ideas {
		if idea.ID == ideaID {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

// Update given idea
func (s *IdeaStorage) Update(number int, title, description string) (*models.Idea, error) {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return nil, err
	}
	idea.Title = title
	idea.Description = description
	return idea, err
}

// GetByNumber returns idea by tenant and number
func (s *IdeaStorage) GetByNumber(number int) (*models.Idea, error) {
	for _, idea := range s.ideas {
		if idea.Number == number {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll() ([]*models.Idea, error) {
	return s.ideas, nil
}

// GetCommentsByIdea returns all coments from given idea
func (s *IdeaStorage) GetCommentsByIdea(number int) ([]*models.Comment, error) {
	return make([]*models.Comment, 0), nil
}

// Add a new idea in the database
func (s *IdeaStorage) Add(title, description string, userID int) (*models.Idea, error) {
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
func (s *IdeaStorage) AddComment(number int, content string, userID int) (int, error) {
	return 0, nil
}

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}
	idea.TotalSupporters = idea.TotalSupporters + 1
	return nil
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}
	idea.TotalSupporters = idea.TotalSupporters - 1
	return nil
}

// SetResponse changes current idea response
func (s *IdeaStorage) SetResponse(number int, text string, userID, status int) error {
	for _, idea := range s.ideas {
		if idea.Number == number {
			idea.Response = &models.IdeaResponse{
				Text:        text,
				User:        &models.User{ID: userID},
				RespondedOn: time.Now(),
			}
		}
	}
	return nil
}

// SupportedBy returns a list of Idea ID supported by given user
func (s *IdeaStorage) SupportedBy(userID int) ([]int, error) {
	return make([]int, 0), nil
}
