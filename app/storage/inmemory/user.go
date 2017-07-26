package inmemory

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// UserStorage is used for user operations
type UserStorage struct {
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
	s.users = append(s.users, user)
	return nil
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
