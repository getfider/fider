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
	lastCommentID    int
	ideas            []*models.Idea
	ideasSupportedBy map[int][]int
	ideaSubscribers  map[int][]int
	ideaComments     map[int][]*models.Comment
	tenant           *models.Tenant
	user             *models.User
}

// NewIdeaStorage creates a new IdeaStorage
func NewIdeaStorage() *IdeaStorage {
	return &IdeaStorage{
		ideas:            make([]*models.Idea, 0),
		ideasSupportedBy: make(map[int][]int, 0),
		ideaSubscribers:  make(map[int][]int, 0),
		ideaComments:     make(map[int][]*models.Comment, 0),
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
func (s *IdeaStorage) Update(idea *models.Idea, title, description string) (*models.Idea, error) {
	idea.Title = title
	idea.Description = description
	return idea, nil
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

// GetCommentsByIdea returns all comments from given idea
func (s *IdeaStorage) GetCommentsByIdea(idea *models.Idea) ([]*models.Comment, error) {
	return s.ideaComments[idea.ID], nil
}

// Add a new idea in the database
func (s *IdeaStorage) Add(title, description string) (*models.Idea, error) {
	s.lastID = s.lastID + 1
	idea := &models.Idea{
		ID:          s.lastID,
		Number:      s.lastID,
		Title:       title,
		Slug:        slug.Make(title),
		Description: description,
		User:        s.user,
	}
	s.ideas = append(s.ideas, idea)
	s.ideasSupportedBy[s.user.ID] = append(s.ideasSupportedBy[s.user.ID], idea.ID)
	return idea, nil
}

// AddComment places a new comment on an idea
func (s *IdeaStorage) AddComment(idea *models.Idea, content string) (int, error) {
	s.lastCommentID++
	s.ideaComments[idea.ID] = append(s.ideaComments[idea.ID], &models.Comment{
		ID:        s.lastCommentID,
		Content:   content,
		CreatedOn: time.Now(),
		User:      s.user,
	})

	return s.lastCommentID, nil
}

// GetCommentByID returns a comment by given ID
func (s *IdeaStorage) GetCommentByID(id int) (*models.Comment, error) {
	for _, comments := range s.ideaComments {
		for _, comment := range comments {
			if comment.ID == id && comment.User.Tenant == s.tenant {
				return comment, nil
			}
		}
	}
	return nil, app.ErrNotFound
}

// UpdateComment with given ID and content
func (s *IdeaStorage) UpdateComment(id int, content string) error {
	now := time.Now()
	comment, err := s.GetCommentByID(id)
	if err != nil {
		return err
	}
	comment.Content = content
	comment.EditedOn = &now
	comment.EditedBy = s.user
	return nil
}

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(idea *models.Idea, user *models.User) error {
	s.ideasSupportedBy[user.ID] = append(s.ideasSupportedBy[user.ID], idea.ID)
	idea.TotalSupporters = idea.TotalSupporters + 1
	return nil
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(idea *models.Idea, user *models.User) error {
	for i, id := range s.ideasSupportedBy[user.ID] {
		if id == idea.ID {
			s.ideasSupportedBy[user.ID] = append(s.ideasSupportedBy[user.ID][:i], s.ideasSupportedBy[user.ID][i+1:]...)
			break
		}
	}
	idea.TotalSupporters = idea.TotalSupporters - 1
	return nil
}

// SetResponse changes current idea response
func (s *IdeaStorage) SetResponse(idea *models.Idea, text string, status int) error {
	for i, storedIdea := range s.ideas {
		if storedIdea.Number == idea.Number {
			if status == models.IdeaDeleted {
				s.ideas = append(s.ideas[:i], s.ideas[i+1:]...)
			} else {
				storedIdea.Status = status
				storedIdea.Response = &models.IdeaResponse{
					Text:        text,
					User:        s.user,
					RespondedOn: time.Now(),
				}
			}
		}
	}
	return nil
}

// MarkAsDuplicate set idea as a duplicate of another idea
func (s *IdeaStorage) MarkAsDuplicate(idea *models.Idea, original *models.Idea) error {
	idea.Status = models.IdeaDuplicate
	idea.Response = &models.IdeaResponse{
		Original: &models.OriginalIdea{
			Number: original.Number,
			Title:  original.Title,
			Slug:   original.Slug,
			Status: original.Status,
		},
		Text:        "",
		User:        s.user,
		RespondedOn: time.Now(),
	}
	return nil
}

// IsReferenced returns true if another idea is referencing given idea
func (s *IdeaStorage) IsReferenced(idea *models.Idea) (bool, error) {
	for _, i := range s.ideas {
		if i.Status == models.IdeaDuplicate && i.Response.Original.Number == idea.Number {
			return true, nil
		}
	}
	return false, nil
}

// SupportedBy returns a list of Idea ID supported by given user
func (s *IdeaStorage) SupportedBy() ([]int, error) {
	return s.ideasSupportedBy[s.user.ID], nil
}

// AddSubscriber adds user to the idea list of subscribers
func (s *IdeaStorage) AddSubscriber(idea *models.Idea, user *models.User) error {
	s.ideaSubscribers[idea.ID] = append(s.ideaSubscribers[idea.ID], user.ID)
	return nil
}

// RemoveSubscriber removes user from idea list of subscribers
func (s *IdeaStorage) RemoveSubscriber(idea *models.Idea, user *models.User) error {
	for i, id := range s.ideaSubscribers[idea.ID] {
		if id == user.ID {
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
