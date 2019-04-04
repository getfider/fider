package inmemory

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
)

// UserStorage is used for user operations
type UserStorage struct {
	user            *models.User
	tenant          *models.Tenant
	users           []*models.User
	lastID          int
	settingsPerUser map[int]map[string]string
	userByAPIKey    map[string]*models.User
}

// GetByID returns a user based on given id
func (s *UserStorage) GetByID(userID int) (*models.User, error) {
	for _, user := range s.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// SetCurrentTenant tenant
func (s *UserStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *UserStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// GetByEmail returns a user based on given email
func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	for _, user := range s.users {
		if user.Email == email && user.Tenant.ID == s.tenant.ID {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetByProvider returns a user based on provider details
func (s *UserStorage) GetByProvider(provider string, uid string) (*models.User, error) {
	for _, user := range s.users {
		for _, item := range user.Providers {
			if item.Name == provider && item.UID == uid {
				return user, nil
			}
		}
	}
	return nil, app.ErrNotFound
}

// Register creates a new user based on given information
func (s *UserStorage) Register(user *models.User) error {
	user.Status = models.UserActive
	if user.ID == 0 {
		s.lastID = s.lastID + 1
		user.ID = s.lastID
	}
	_, err := s.GetByEmail(user.Email)
	if errors.Cause(err) == app.ErrNotFound || user.Email == "" {
		s.users = append(s.users, user)
		return nil
	}
	return errors.New("User already registered")
}

// RegisterProvider adds given provider to userID
func (s *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	user, err := s.GetByID(userID)
	if err == nil {
		user.Providers = append(user.Providers, provider)
	}
	return err
}

// Update user settings
func (s *UserStorage) Update(settings *models.UpdateUserSettings) error {
	user, err := s.GetByID(s.user.ID)
	if err == nil {
		user.Name = settings.Name
	}
	return err
}

// ChangeRole of given user
func (s *UserStorage) ChangeRole(userID int, role models.Role) error {
	user, err := s.GetByID(userID)
	if err == nil {
		user.Role = role
	}
	return err
}

// ChangeEmail of given user
func (s *UserStorage) ChangeEmail(userID int, email string) error {
	user, err := s.GetByID(userID)
	if err == nil {
		user.Email = email
	}
	return err
}

// GetAll return all users of current tenant
func (s *UserStorage) GetAll() ([]*models.User, error) {
	return s.users, nil
}

// UpdateSettings of given user
func (s *UserStorage) UpdateSettings(settings map[string]string) error {
	if s.user == nil {
		return nil
	}
	if s.settingsPerUser == nil {
		s.settingsPerUser = make(map[int]map[string]string, 0)
	}
	s.settingsPerUser[s.user.ID] = settings
	return nil
}

// GetUserSettings returns current user's settings
func (s *UserStorage) GetUserSettings() (map[string]string, error) {
	if s.settingsPerUser == nil || s.user == nil {
		return make(map[string]string, 0), nil
	}
	settings, ok := s.settingsPerUser[s.user.ID]
	if ok {
		return settings, nil
	}
	return make(map[string]string, 0), nil
}

// HasSubscribedTo returns true if current user is receiving notification from specific post
func (s *UserStorage) HasSubscribedTo(postID int) (bool, error) {
	return false, nil
}

// Delete removes current user personal data and mark it as deleted
func (s *UserStorage) Delete() error {
	s.user.Name = ""
	s.user.Email = ""
	s.user.Role = models.RoleVisitor
	s.user.Status = models.UserDeleted
	s.user.Providers = make([]*models.UserProvider, 0)
	return nil
}

// RegenerateAPIKey generates a new API Key and returns it
func (s *UserStorage) RegenerateAPIKey() (string, error) {
	apiKey := models.GenerateSecretKey()
	if s.userByAPIKey == nil {
		s.userByAPIKey = make(map[string]*models.User)
	}
	s.userByAPIKey[apiKey] = s.user
	return apiKey, nil
}

// GetByAPIKey returns a user based on its API key
func (s *UserStorage) GetByAPIKey(apiKey string) (*models.User, error) {
	if s.userByAPIKey == nil {
		s.userByAPIKey = make(map[string]*models.User)
	}

	for userAPIKey, user := range s.userByAPIKey {
		if userAPIKey == apiKey {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// Block a given user from using Fider
func (s *UserStorage) Block(userID int) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	if user.Status == models.UserActive {
		user.Status = models.UserBlocked
	}
	return nil
}

// Unblock a given user so that they can use Fider again
func (s *UserStorage) Unblock(userID int) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	if user.Status == models.UserBlocked {
		user.Status = models.UserActive
	}
	return nil
}
