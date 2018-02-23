package inmemory

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	lastID           int
	ideas            []*models.Idea
	ideasSupportedBy map[int][]int
	ideaSubscribers  map[int][]int
	tenant           *models.Tenant
	user             *models.User
}

// NewIdeaStorage creates a new IdeaStorage
func NewIdeaStorage() *IdeaStorage {
	return &IdeaStorage{
		ideas:            make([]*models.Idea, 0),
		ideasSupportedBy: make(map[int][]int, 0),
		ideaSubscribers:  make(map[int][]int, 0),
	}
}

// SetCurrentTenant to current context
func (s *IdeaStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *IdeaStorage) SetCurrentUser(user *models.User) {
	s.user = user
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

// GetBySlug returns idea by tenant and slug
func (s *IdeaStorage) GetBySlug(slug string) (*models.Idea, error) {
	for _, idea := range s.ideas {
		if idea.Slug == slug {
			return idea, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll() ([]*models.Idea, error) {
	return s.ideas, nil
}

// CountPerStatus returns total number of ideas per status
func (s *IdeaStorage) CountPerStatus() (map[int]int, error) {
	return make(map[int]int, 0), nil
}

// Search existing ideas based on input
func (s *IdeaStorage) Search(query, filter string, tags []string) ([]*models.Idea, error) {
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
		Slug:        slug.Make(title),
		Description: description,
		User: &models.User{
			ID: userID,
		},
	}
	s.ideas = append(s.ideas, idea)
	s.ideasSupportedBy[userID] = append(s.ideasSupportedBy[userID], idea.ID)
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
	s.ideasSupportedBy[userID] = append(s.ideasSupportedBy[userID], idea.ID)
	idea.TotalSupporters = idea.TotalSupporters + 1
	return nil
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	for i, id := range s.ideasSupportedBy[userID] {
		if id == idea.ID {
			s.ideasSupportedBy[userID] = append(s.ideasSupportedBy[userID][:i], s.ideasSupportedBy[userID][i+1:]...)
			break
		}
	}
	idea.TotalSupporters = idea.TotalSupporters - 1
	return nil
}

// SetResponse changes current idea response
func (s *IdeaStorage) SetResponse(number int, text string, userID, status int) error {
	for _, idea := range s.ideas {
		if idea.Number == number {
			idea.Status = status
			idea.Response = &models.IdeaResponse{
				Text:        text,
				User:        &models.User{ID: userID},
				RespondedOn: time.Now(),
			}
		}
	}
	return nil
}

// MarkAsDuplicate set idea as a duplicate of another idea
func (s *IdeaStorage) MarkAsDuplicate(number, originalNumber, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}
	original, err := s.GetByNumber(originalNumber)
	if err != nil {
		return err
	}

	idea.Status = models.IdeaDuplicate
	idea.Response = &models.IdeaResponse{
		Original: &models.OriginalIdea{
			Number: original.Number,
			Title:  original.Title,
			Slug:   original.Slug,
			Status: original.Status,
		},
		Text:        "",
		User:        &models.User{ID: userID},
		RespondedOn: time.Now(),
	}
	return nil
}

// SupportedBy returns a list of Idea ID supported by given user
func (s *IdeaStorage) SupportedBy(userID int) ([]int, error) {
	return s.ideasSupportedBy[userID], nil
}

// AddSubscriber adds user to the idea list of subscribers
func (s *IdeaStorage) AddSubscriber(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}
	s.ideaSubscribers[idea.ID] = append(s.ideaSubscribers[idea.ID], userID)
	return nil
}

// RemoveSubscriber removes user from idea list of subscribers
func (s *IdeaStorage) RemoveSubscriber(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	for i, id := range s.ideaSubscribers[idea.ID] {
		if id == userID {
			s.ideaSubscribers[idea.ID] = append(s.ideasSupportedBy[idea.ID][:i], s.ideasSupportedBy[idea.ID][i+1:]...)
			break
		}
	}
	return nil
}

// GetActiveSubscribers based on input and settings
func (s *IdeaStorage) GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error) {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return nil, err
	}
	subscribers, ok := s.ideaSubscribers[idea.ID]
	if ok {
		users := make([]*models.User, len(subscribers))
		for i, id := range subscribers {
			users[i] = &models.User{
				ID:    id,
				Name:  fmt.Sprintf("User %d", id),
				Email: fmt.Sprintf("user%d@test.com", id),
			}
		}
	}
	return make([]*models.User, 0), nil
}
