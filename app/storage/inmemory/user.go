package inmemory

import (
	"errors"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// UserStorage is used for user operations
type UserStorage struct {
	user   *models.User
	tenant *models.Tenant
	users  []*models.User
	lastID int
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
func (s *UserStorage) GetByEmail(tenantID int, email string) (*models.User, error) {
	for _, user := range s.users {
		if user.Email == email && user.Tenant.ID == tenantID {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetByProvider returns a user based on provider details
func (s *UserStorage) GetByProvider(tenantID int, provider string, uid string) (*models.User, error) {
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
	if user.ID == 0 {
		s.lastID = s.lastID + 1
		user.ID = s.lastID
	}
	_, err := s.GetByEmail(user.Tenant.ID, user.Email)
	if err == app.ErrNotFound || user.Email == "" {
		s.users = append(s.users, user)
		return nil
	}
	return errors.New("User already registered")
}

// RegisterProvider adds given provider to userID
func (s *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	for _, user := range s.users {
		if user.ID == userID {
			user.Providers = append(user.Providers, provider)
			return nil
		}
	}
	return app.ErrNotFound
}

// Update user settings
func (s *UserStorage) Update(userID int, settings *models.UpdateUserSettings) error {
	for _, user := range s.users {
		if user.ID == userID {
			user.Name = settings.Name
			return nil
		}
	}
	return app.ErrNotFound
}

// ChangeRole of given user
func (s *UserStorage) ChangeRole(userID int, role models.Role) error {
	for _, user := range s.users {
		if user.ID == userID {
			user.Role = role
			return nil
		}
	}
	return app.ErrNotFound
}

// GetAll return all users of current tenant
func (s *UserStorage) GetAll() ([]*models.User, error) {
	return s.users, nil
}
